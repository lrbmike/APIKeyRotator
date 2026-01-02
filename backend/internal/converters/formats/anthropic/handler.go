package anthropic

import (
	"encoding/json"
	"fmt"

	"api-key-rotator/backend/internal/converters/formats"
)

func init() {
	formats.RegisterFormat("anthropic", &Handler{}, &StreamHandler{})
}

// Handler implements FormatHandler for Anthropic format
type Handler struct{}

func (h *Handler) Name() string {
	return "anthropic"
}

// Request types
type Request struct {
	Model         string           `json:"model"`
	Messages      []RequestMessage `json:"messages"`
	System        interface{}      `json:"system,omitempty"`
	MaxTokens     int              `json:"max_tokens"`
	Temperature   *float64         `json:"temperature,omitempty"`
	TopP          *float64         `json:"top_p,omitempty"`
	Stream        bool             `json:"stream,omitempty"`
	StopSequences []string         `json:"stop_sequences,omitempty"`
}

type RequestMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// Response types
type Response struct {
	ID           string    `json:"id"`
	Type         string    `json:"type"`
	Role         string    `json:"role"`
	Content      []Content `json:"content"`
	Model        string    `json:"model"`
	StopReason   *string   `json:"stop_reason,omitempty"`
	StopSequence *string   `json:"stop_sequence,omitempty"`
	Usage        *Usage    `json:"usage,omitempty"`
}

type Content struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// ParseRequest implements FormatHandler
func (h *Handler) ParseRequest(body []byte) (*formats.UniversalRequest, error) {
	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic request: %w", err)
	}

	universal := &formats.UniversalRequest{
		Model:       req.Model,
		Messages:    make([]formats.UniversalMessage, 0),
		MaxTokens:   req.MaxTokens,
		Stream:      req.Stream,
		Stop:        req.StopSequences,
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}

	universal.System = extractSystemContent(req.System)

	for _, msg := range req.Messages {
		universal.Messages = append(universal.Messages, formats.UniversalMessage{
			Role:    msg.Role,
			Content: extractTextContent(msg.Content),
		})
	}

	return universal, nil
}

// BuildRequest implements FormatHandler
func (h *Handler) BuildRequest(req *formats.UniversalRequest) ([]byte, error) {
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 4096
	}

	anthropicReq := Request{
		Model:         req.Model,
		Messages:      make([]RequestMessage, 0),
		MaxTokens:     maxTokens,
		Stream:        req.Stream,
		StopSequences: req.Stop,
		Temperature:   req.Temperature,
		TopP:          req.TopP,
	}

	if req.System != "" {
		anthropicReq.System = req.System
	}

	for _, msg := range req.Messages {
		anthropicReq.Messages = append(anthropicReq.Messages, RequestMessage{
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	return json.Marshal(anthropicReq)
}

// ParseResponse implements FormatHandler
func (h *Handler) ParseResponse(body []byte) (*formats.UniversalResponse, error) {
	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic response: %w", err)
	}

	universal := &formats.UniversalResponse{
		ID:    resp.ID,
		Model: resp.Model,
		Role:  resp.Role,
	}

	for _, content := range resp.Content {
		if content.Type == "text" {
			universal.Content += content.Text
		}
	}

	if resp.StopReason != nil {
		universal.StopReason = *resp.StopReason
	}

	if resp.Usage != nil {
		universal.Usage = &formats.UniversalUsage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
			TotalTokens:  resp.Usage.InputTokens + resp.Usage.OutputTokens,
		}
	}

	return universal, nil
}

// BuildResponse implements FormatHandler
func (h *Handler) BuildResponse(resp *formats.UniversalResponse) ([]byte, error) {
	stopReason := resp.StopReason
	if stopReason == "" {
		stopReason = "end_turn"
	} else {
		switch stopReason {
		case "stop":
			stopReason = "end_turn"
		case "length":
			stopReason = "max_tokens"
		}
	}

	anthropicResp := Response{
		ID:         resp.ID,
		Type:       "message",
		Role:       "assistant",
		Model:      resp.Model,
		StopReason: &stopReason,
		Content:    []Content{{Type: "text", Text: resp.Content}},
	}

	if resp.Usage != nil {
		anthropicResp.Usage = &Usage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
		}
	} else {
		anthropicResp.Usage = &Usage{InputTokens: 0, OutputTokens: 0}
	}

	return json.Marshal(anthropicResp)
}

// GetAPIPath implements FormatHandler
func (h *Handler) GetAPIPath(action string) string {
	switch action {
	case "chat/completions", "v1/chat/completions":
		return "v1/messages"
	default:
		return action
	}
}

// GetClientAction implements FormatHandler
func (h *Handler) GetClientAction(apiPath string) string {
	switch apiPath {
	case "v1/messages", "messages":
		return "chat/completions"
	default:
		return apiPath
	}
}

// Helper functions
func extractSystemContent(system interface{}) string {
	switch s := system.(type) {
	case string:
		return s
	case []interface{}:
		var text string
		for _, block := range s {
			if blockMap, ok := block.(map[string]interface{}); ok {
				if blockMap["type"] == "text" {
					if t, ok := blockMap["text"].(string); ok {
						text += t
					}
				}
			}
		}
		return text
	default:
		return ""
	}
}

func extractTextContent(content interface{}) string {
	switch c := content.(type) {
	case string:
		return c
	case []interface{}:
		var text string
		for _, block := range c {
			if blockMap, ok := block.(map[string]interface{}); ok {
				if blockMap["type"] == "text" {
					if t, ok := blockMap["text"].(string); ok {
						text += t
					}
				}
			}
		}
		return text
	default:
		return ""
	}
}
