package openai

import (
	"encoding/json"
	"fmt"
	"time"

	"api-key-rotator/backend/internal/converters/formats"
)

func init() {
	formats.RegisterFormat("openai", &Handler{}, &StreamHandler{})
}

// Handler implements FormatHandler for OpenAI format
type Handler struct{}

func (h *Handler) Name() string {
	return "openai"
}

// Request types
type Request struct {
	Model       string           `json:"model"`
	Messages    []RequestMessage `json:"messages"`
	MaxTokens   *int             `json:"max_tokens,omitempty"`
	Temperature *float64         `json:"temperature,omitempty"`
	TopP        *float64         `json:"top_p,omitempty"`
	Stream      bool             `json:"stream,omitempty"`
	Stop        []string         `json:"stop,omitempty"`
}

type RequestMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Response types
type Response struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   *Usage   `json:"usage,omitempty"`
}

type Choice struct {
	Index        int      `json:"index"`
	Message      *Message `json:"message,omitempty"`
	FinishReason *string  `json:"finish_reason,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ParseRequest implements FormatHandler
func (h *Handler) ParseRequest(body []byte) (*formats.UniversalRequest, error) {
	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI request: %w", err)
	}

	universal := &formats.UniversalRequest{
		Model:       req.Model,
		Messages:    make([]formats.UniversalMessage, 0),
		Stream:      req.Stream,
		Stop:        req.Stop,
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}

	if req.MaxTokens != nil {
		universal.MaxTokens = *req.MaxTokens
	}

	for _, msg := range req.Messages {
		if msg.Role == "system" {
			universal.System = msg.Content
		} else {
			universal.Messages = append(universal.Messages, formats.UniversalMessage{
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	return universal, nil
}

// BuildRequest implements FormatHandler
func (h *Handler) BuildRequest(req *formats.UniversalRequest) ([]byte, error) {
	openaiReq := Request{
		Model:       req.Model,
		Messages:    make([]RequestMessage, 0),
		Stream:      req.Stream,
		Stop:        req.Stop,
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}

	if req.MaxTokens > 0 {
		openaiReq.MaxTokens = &req.MaxTokens
	}

	if req.System != "" {
		openaiReq.Messages = append(openaiReq.Messages, RequestMessage{
			Role:    "system",
			Content: req.System,
		})
	}

	for _, msg := range req.Messages {
		openaiReq.Messages = append(openaiReq.Messages, RequestMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return json.Marshal(openaiReq)
}

// ParseResponse implements FormatHandler
func (h *Handler) ParseResponse(body []byte) (*formats.UniversalResponse, error) {
	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	universal := &formats.UniversalResponse{
		ID:    resp.ID,
		Model: resp.Model,
		Role:  "assistant",
	}

	if len(resp.Choices) > 0 {
		choice := resp.Choices[0]
		if choice.Message != nil {
			universal.Content = choice.Message.Content
			universal.Role = choice.Message.Role
		}
		if choice.FinishReason != nil {
			universal.StopReason = *choice.FinishReason
		}
	}

	if resp.Usage != nil {
		universal.Usage = &formats.UniversalUsage{
			InputTokens:  resp.Usage.PromptTokens,
			OutputTokens: resp.Usage.CompletionTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		}
	}

	return universal, nil
}

// BuildResponse implements FormatHandler
func (h *Handler) BuildResponse(resp *formats.UniversalResponse) ([]byte, error) {
	finishReason := resp.StopReason
	if finishReason == "" {
		finishReason = "stop"
	}

	openaiResp := Response{
		ID:      resp.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   resp.Model,
		Choices: []Choice{
			{
				Index: 0,
				Message: &Message{
					Role:    resp.Role,
					Content: resp.Content,
				},
				FinishReason: &finishReason,
			},
		},
	}

	if resp.Usage != nil {
		openaiResp.Usage = &Usage{
			PromptTokens:     resp.Usage.InputTokens,
			CompletionTokens: resp.Usage.OutputTokens,
			TotalTokens:      resp.Usage.TotalTokens,
		}
	}

	return json.Marshal(openaiResp)
}

// GetAPIPath implements FormatHandler - converts client action to OpenAI API path
func (h *Handler) GetAPIPath(action string) string {
	switch action {
	case "v1/messages", "messages":
		// Anthropic endpoint -> OpenAI endpoint
		return "chat/completions"
	case "v1/chat/completions":
		return "chat/completions"
	default:
		return action
	}
}

// GetClientAction implements FormatHandler - converts API path to client action
func (h *Handler) GetClientAction(apiPath string) string {
	switch apiPath {
	case "chat/completions", "v1/chat/completions":
		return "v1/chat/completions"
	default:
		return apiPath
	}
}
