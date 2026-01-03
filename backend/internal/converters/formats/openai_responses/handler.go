package openai_responses

import (
	"encoding/json"
	"fmt"
	"time"

	"api-key-rotator/backend/internal/converters/formats"

	"github.com/google/uuid"
)

func init() {
	formats.RegisterFormat("openai_responses", &Handler{}, &StreamHandler{})
}

// Handler implements FormatHandler for OpenAI Responses API format (v1/responses)
type Handler struct{}

func (h *Handler) Name() string {
	return "openai_responses"
}

// Request types for OpenAI Responses API
type Request struct {
	Model              string      `json:"model"`
	Input              interface{} `json:"input"`                  // Can be string or array of input items
	Instructions       string      `json:"instructions,omitempty"` // System instructions
	MaxOutputTokens    *int        `json:"max_output_tokens,omitempty"`
	Temperature        *float64    `json:"temperature,omitempty"`
	TopP               *float64    `json:"top_p,omitempty"`
	Stream             bool        `json:"stream,omitempty"`
	PreviousResponseID string      `json:"previous_response_id,omitempty"` // For multi-turn conversations
	Store              *bool       `json:"store,omitempty"`
	Metadata           interface{} `json:"metadata,omitempty"`
}

// InputItem represents an item in the input array
type InputItem struct {
	Type    string      `json:"type"`              // "message"
	Role    string      `json:"role,omitempty"`    // "user", "assistant", "system"
	Content interface{} `json:"content,omitempty"` // Can be string or array of content parts
}

// Response types for OpenAI Responses API
type Response struct {
	ID        string       `json:"id"`
	Object    string       `json:"object"` // "response"
	CreatedAt int64        `json:"created_at"`
	Model     string       `json:"model"`
	Output    []OutputItem `json:"output"`
	Usage     *Usage       `json:"usage,omitempty"`
	Status    string       `json:"status"` // "completed", "failed", etc.
	Error     interface{}  `json:"error,omitempty"`
}

// OutputItem represents an item in the output array
type OutputItem struct {
	Type    string        `json:"type"` // "message"
	ID      string        `json:"id"`
	Role    string        `json:"role,omitempty"`   // "assistant"
	Status  string        `json:"status,omitempty"` // "completed"
	Content []ContentPart `json:"content,omitempty"`
}

// ContentPart represents a part of the content
type ContentPart struct {
	Type string `json:"type"` // "output_text"
	Text string `json:"text,omitempty"`
}

// Usage token usage statistics
type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// ParseRequest implements FormatHandler - parses v1/responses request to universal format
func (h *Handler) ParseRequest(body []byte) (*formats.UniversalRequest, error) {
	var req Request
	if err := json.Unmarshal(body, &req); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI Responses request: %w", err)
	}

	universal := &formats.UniversalRequest{
		Model:       req.Model,
		Messages:    make([]formats.UniversalMessage, 0),
		System:      req.Instructions,
		Stream:      req.Stream,
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}

	if req.MaxOutputTokens != nil {
		universal.MaxTokens = *req.MaxOutputTokens
	}

	// Parse input - can be string or array of input items
	switch input := req.Input.(type) {
	case string:
		// Simple string input - treat as user message
		universal.Messages = append(universal.Messages, formats.UniversalMessage{
			Role:    "user",
			Content: input,
		})
	case []interface{}:
		// Array of input items
		for _, item := range input {
			if itemMap, ok := item.(map[string]interface{}); ok {
				msg := parseInputItem(itemMap)
				if msg != nil {
					if msg.Role == "system" {
						// Move system message to System field
						universal.System = msg.Content
					} else {
						universal.Messages = append(universal.Messages, *msg)
					}
				}
			}
		}
	}

	return universal, nil
}

// parseInputItem extracts a message from an input item map
func parseInputItem(item map[string]interface{}) *formats.UniversalMessage {
	itemType, _ := item["type"].(string)
	if itemType != "message" {
		return nil
	}

	role, _ := item["role"].(string)
	if role == "" {
		role = "user"
	}

	var content string
	switch c := item["content"].(type) {
	case string:
		content = c
	case []interface{}:
		// Array of content parts - extract text
		for _, part := range c {
			if partMap, ok := part.(map[string]interface{}); ok {
				if partMap["type"] == "input_text" || partMap["type"] == "text" {
					if text, ok := partMap["text"].(string); ok {
						content += text
					}
				}
			}
		}
	}

	return &formats.UniversalMessage{
		Role:    role,
		Content: content,
	}
}

// BuildRequest implements FormatHandler - builds v1/responses request from universal format
func (h *Handler) BuildRequest(req *formats.UniversalRequest) ([]byte, error) {
	responsesReq := Request{
		Model:        req.Model,
		Instructions: req.System,
		Stream:       req.Stream,
		Temperature:  req.Temperature,
		TopP:         req.TopP,
	}

	if req.MaxTokens > 0 {
		responsesReq.MaxOutputTokens = &req.MaxTokens
	}

	// Build input array from messages
	input := make([]InputItem, 0, len(req.Messages))
	for _, msg := range req.Messages {
		input = append(input, InputItem{
			Type:    "message",
			Role:    msg.Role,
			Content: msg.Content,
		})
	}

	if len(input) == 1 {
		// Single message - can be simplified to string
		responsesReq.Input = input[0].Content
	} else {
		responsesReq.Input = input
	}

	return json.Marshal(responsesReq)
}

// ParseResponse implements FormatHandler - parses v1/responses response to universal format
func (h *Handler) ParseResponse(body []byte) (*formats.UniversalResponse, error) {
	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI Responses response: %w", err)
	}

	universal := &formats.UniversalResponse{
		ID:    resp.ID,
		Model: resp.Model,
		Role:  "assistant",
	}

	// Extract content from output items
	for _, item := range resp.Output {
		if item.Type == "message" && item.Role == "assistant" {
			for _, part := range item.Content {
				if part.Type == "output_text" {
					universal.Content += part.Text
				}
			}
		}
	}

	// Map status to stop reason
	if resp.Status == "completed" {
		universal.StopReason = "stop"
	} else if resp.Status != "" {
		universal.StopReason = resp.Status
	}

	if resp.Usage != nil {
		universal.Usage = &formats.UniversalUsage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		}
	}

	return universal, nil
}

// BuildResponse implements FormatHandler - builds v1/responses response from universal format
func (h *Handler) BuildResponse(resp *formats.UniversalResponse) ([]byte, error) {
	messageID := "msg_" + uuid.New().String()[:12]

	status := "completed"
	if resp.StopReason != "" && resp.StopReason != "stop" {
		status = resp.StopReason
	}

	responsesResp := Response{
		ID:        resp.ID,
		Object:    "response",
		CreatedAt: time.Now().Unix(),
		Model:     resp.Model,
		Status:    status,
		Output: []OutputItem{
			{
				Type:   "message",
				ID:     messageID,
				Role:   "assistant",
				Status: "completed",
				Content: []ContentPart{
					{
						Type: "output_text",
						Text: resp.Content,
					},
				},
			},
		},
	}

	if resp.Usage != nil {
		responsesResp.Usage = &Usage{
			InputTokens:  resp.Usage.InputTokens,
			OutputTokens: resp.Usage.OutputTokens,
			TotalTokens:  resp.Usage.TotalTokens,
		}
	}

	return json.Marshal(responsesResp)
}

// GetAPIPath implements FormatHandler - converts client action to target API path
func (h *Handler) GetAPIPath(action string) string {
	// When client sends to v1/responses, we need to convert to chat/completions
	// for the actual backend API call
	if action == "v1/responses" || action == "responses" {
		return "v1/chat/completions"
	}
	return action
}

// GetClientAction implements FormatHandler - converts API path to client action
func (h *Handler) GetClientAction(apiPath string) string {
	if apiPath == "v1/chat/completions" || apiPath == "chat/completions" {
		return "v1/responses"
	}
	return apiPath
}
