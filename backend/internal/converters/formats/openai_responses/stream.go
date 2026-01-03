package openai_responses

import (
	"encoding/json"
	"fmt"
	"time"

	"api-key-rotator/backend/internal/converters/formats"

	"github.com/google/uuid"
)

// StreamHandler implements formats.StreamHandler for OpenAI Responses API format
type StreamHandler struct{}

func (h *StreamHandler) Name() string {
	return "openai_responses"
}

// Stream event types for Responses API
type StreamEvent struct {
	Type         string    `json:"type"`
	Response     *Response `json:"response,omitempty"`
	Delta        string    `json:"delta,omitempty"`
	ItemID       string    `json:"item_id,omitempty"`
	OutputIndex  int       `json:"output_index,omitempty"`
	ContentIndex int       `json:"content_index,omitempty"`
}

// ResponseCreatedEvent is emitted at the start of streaming
type ResponseCreatedEvent struct {
	Type     string   `json:"type"` // "response.created"
	Response Response `json:"response"`
}

// ResponseDoneEvent is emitted at the end of streaming
type ResponseDoneEvent struct {
	Type     string   `json:"type"` // "response.done"
	Response Response `json:"response"`
}

// OutputTextDeltaEvent is emitted for each text chunk
type OutputTextDeltaEvent struct {
	Type         string `json:"type"` // "response.output_text.delta"
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Delta        string `json:"delta"`
}

// OutputTextDoneEvent is emitted when text output is complete
type OutputTextDoneEvent struct {
	Type         string `json:"type"` // "response.output_text.done"
	ItemID       string `json:"item_id"`
	OutputIndex  int    `json:"output_index"`
	ContentIndex int    `json:"content_index"`
	Text         string `json:"text"`
}

// ParseStreamChunk implements StreamHandler - parses Responses API stream chunk to universal format
func (h *StreamHandler) ParseStreamChunk(chunk []byte) (*formats.UniversalStreamChunk, error) {
	// First try to parse as a generic event to determine type
	var event struct {
		Type     string    `json:"type"`
		Delta    string    `json:"delta,omitempty"`
		Text     string    `json:"text,omitempty"`
		ItemID   string    `json:"item_id,omitempty"`
		Response *Response `json:"response,omitempty"`
	}

	if err := json.Unmarshal(chunk, &event); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI Responses stream chunk: %w", err)
	}

	universal := &formats.UniversalStreamChunk{}

	switch event.Type {
	case "response.created":
		// Initial event - mark as first
		universal.IsFirst = true
		if event.Response != nil {
			universal.ID = event.Response.ID
			universal.Model = event.Response.Model
		}
		universal.Role = "assistant"

	case "response.output_text.delta":
		// Text delta event
		universal.Delta = event.Delta

	case "response.output_text.done":
		// Text complete event - could include full text
		// Not marking as last since response.done comes after

	case "response.done":
		// Final event
		universal.IsLast = true
		if event.Response != nil {
			stopReason := "stop"
			if event.Response.Status != "completed" {
				stopReason = event.Response.Status
			}
			universal.StopReason = &stopReason
		}

	case "response.output_item.added", "response.output_item.done",
		"response.content_part.added", "response.content_part.done":
		// Structural events - skip

	default:
		// Unknown event type - skip
	}

	return universal, nil
}

// BuildStreamChunk implements StreamHandler - builds Responses API stream chunk from universal format
func (h *StreamHandler) BuildStreamChunk(chunk *formats.UniversalStreamChunk) ([]byte, error) {
	// Generate consistent IDs for the stream
	itemID := "item_" + uuid.New().String()[:8]

	if chunk.IsFirst {
		// Build response.created event
		event := ResponseCreatedEvent{
			Type: "response.created",
			Response: Response{
				ID:        chunk.ID,
				Object:    "response",
				CreatedAt: time.Now().Unix(),
				Model:     chunk.Model,
				Status:    "in_progress",
				Output:    []OutputItem{},
			},
		}
		return json.Marshal(event)
	}

	if chunk.StopReason != nil {
		// Build response.done event
		status := "completed"
		if *chunk.StopReason != "stop" && *chunk.StopReason != "" {
			status = *chunk.StopReason
		}

		event := ResponseDoneEvent{
			Type: "response.done",
			Response: Response{
				ID:        chunk.ID,
				Object:    "response",
				CreatedAt: time.Now().Unix(),
				Model:     chunk.Model,
				Status:    status,
				Output:    []OutputItem{},
			},
		}
		return json.Marshal(event)
	}

	if chunk.Delta != "" {
		// Build response.output_text.delta event
		event := OutputTextDeltaEvent{
			Type:         "response.output_text.delta",
			ItemID:       itemID,
			OutputIndex:  0,
			ContentIndex: 0,
			Delta:        chunk.Delta,
		}
		return json.Marshal(event)
	}

	return nil, nil
}

// BuildStartEvent implements StreamHandler - returns initial events for streaming
func (h *StreamHandler) BuildStartEvent(model string, id string) [][]byte {
	if id == "" {
		id = "resp_" + uuid.New().String()[:12]
	}

	itemID := "item_" + uuid.New().String()[:8]

	events := make([][]byte, 0)

	// 1. response.created
	createdEvent := ResponseCreatedEvent{
		Type: "response.created",
		Response: Response{
			ID:        id,
			Object:    "response",
			CreatedAt: time.Now().Unix(),
			Model:     model,
			Status:    "in_progress",
			Output:    []OutputItem{},
		},
	}
	if data, err := json.Marshal(createdEvent); err == nil {
		events = append(events, data)
	}

	// 2. response.output_item.added
	itemAddedEvent := struct {
		Type        string     `json:"type"`
		OutputIndex int        `json:"output_index"`
		Item        OutputItem `json:"item"`
	}{
		Type:        "response.output_item.added",
		OutputIndex: 0,
		Item: OutputItem{
			Type:    "message",
			ID:      itemID,
			Role:    "assistant",
			Status:  "in_progress",
			Content: []ContentPart{},
		},
	}
	if data, err := json.Marshal(itemAddedEvent); err == nil {
		events = append(events, data)
	}

	// 3. response.content_part.added
	contentAddedEvent := struct {
		Type         string      `json:"type"`
		ItemID       string      `json:"item_id"`
		OutputIndex  int         `json:"output_index"`
		ContentIndex int         `json:"content_index"`
		Part         ContentPart `json:"part"`
	}{
		Type:         "response.content_part.added",
		ItemID:       itemID,
		OutputIndex:  0,
		ContentIndex: 0,
		Part: ContentPart{
			Type: "output_text",
			Text: "",
		},
	}
	if data, err := json.Marshal(contentAddedEvent); err == nil {
		events = append(events, data)
	}

	return events
}

// BuildEndEvent implements StreamHandler - returns final events for streaming
func (h *StreamHandler) BuildEndEvent() [][]byte {
	events := make([][]byte, 0)

	// response.done - simplified
	doneEvent := struct {
		Type string `json:"type"`
	}{
		Type: "response.done",
	}
	if data, err := json.Marshal(doneEvent); err == nil {
		events = append(events, data)
	}

	return events
}
