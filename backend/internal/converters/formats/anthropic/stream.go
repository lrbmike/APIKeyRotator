package anthropic

import (
	"encoding/json"
	"fmt"

	"api-key-rotator/backend/internal/converters/formats"
)

// StreamHandler implements formats.StreamHandler for Anthropic format
type StreamHandler struct{}

func (h *StreamHandler) Name() string {
	return "anthropic"
}

// Stream event types
type StreamEvent struct {
	Type         string       `json:"type"`
	Message      *Response    `json:"message,omitempty"`
	Index        *int         `json:"index,omitempty"`
	ContentBlock *Content     `json:"content_block,omitempty"`
	Delta        *StreamDelta `json:"delta,omitempty"`
	Usage        *StreamUsage `json:"usage,omitempty"`
}

type StreamDelta struct {
	Type       string  `json:"type,omitempty"`
	Text       string  `json:"text,omitempty"`
	StopReason *string `json:"stop_reason,omitempty"`
}

type StreamUsage struct {
	OutputTokens int `json:"output_tokens"`
}

// ParseStreamChunk implements StreamHandler
func (h *StreamHandler) ParseStreamChunk(chunk []byte) (*formats.UniversalStreamChunk, error) {
	var event StreamEvent
	if err := json.Unmarshal(chunk, &event); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic stream event: %w", err)
	}

	universal := &formats.UniversalStreamChunk{}

	switch event.Type {
	case "message_start":
		if event.Message != nil {
			universal.ID = event.Message.ID
			universal.Model = event.Message.Model
			universal.Role = event.Message.Role
			universal.IsFirst = true
		}
	case "content_block_delta":
		if event.Delta != nil {
			universal.Delta = event.Delta.Text
		}
	case "message_delta":
		if event.Delta != nil && event.Delta.StopReason != nil {
			universal.StopReason = event.Delta.StopReason
			universal.IsLast = true
		}
	case "message_stop":
		universal.IsLast = true
	}

	return universal, nil
}

// BuildStreamChunk implements StreamHandler
func (h *StreamHandler) BuildStreamChunk(chunk *formats.UniversalStreamChunk) ([]byte, error) {
	if chunk.Delta != "" {
		event := StreamEvent{
			Type:  "content_block_delta",
			Index: intPtr(0),
			Delta: &StreamDelta{
				Type: "text_delta",
				Text: chunk.Delta,
			},
		}
		return json.Marshal(event)
	}

	if chunk.StopReason != nil {
		stopReason := *chunk.StopReason
		switch stopReason {
		case "stop":
			stopReason = "end_turn"
		case "length":
			stopReason = "max_tokens"
		}
		event := StreamEvent{
			Type:  "message_delta",
			Delta: &StreamDelta{StopReason: &stopReason},
		}
		return json.Marshal(event)
	}

	if chunk.IsFirst {
		event := StreamEvent{
			Type: "message_start",
			Message: &Response{
				ID:    chunk.ID,
				Type:  "message",
				Role:  "assistant",
				Model: chunk.Model,
			},
		}
		return json.Marshal(event)
	}

	return nil, nil
}

// BuildStartEvent implements StreamHandler
func (h *StreamHandler) BuildStartEvent(model string, id string) [][]byte {
	events := make([][]byte, 0)

	messageStart := StreamEvent{
		Type: "message_start",
		Message: &Response{
			ID:      id,
			Type:    "message",
			Role:    "assistant",
			Model:   model,
			Content: []Content{},
		},
	}
	if data, err := json.Marshal(messageStart); err == nil {
		events = append(events, data)
	}

	contentStart := map[string]interface{}{
		"type":          "content_block_start",
		"index":         0,
		"content_block": map[string]string{"type": "text", "text": ""},
	}
	if data, err := json.Marshal(contentStart); err == nil {
		events = append(events, data)
	}

	return events
}

// BuildEndEvent implements StreamHandler
func (h *StreamHandler) BuildEndEvent() [][]byte {
	events := make([][]byte, 0)

	contentStop := map[string]interface{}{"type": "content_block_stop", "index": 0}
	if data, err := json.Marshal(contentStop); err == nil {
		events = append(events, data)
	}

	messageStop := map[string]string{"type": "message_stop"}
	if data, err := json.Marshal(messageStop); err == nil {
		events = append(events, data)
	}

	return events
}

func intPtr(i int) *int {
	return &i
}
