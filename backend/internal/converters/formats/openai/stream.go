package openai

import (
	"encoding/json"
	"fmt"
	"time"

	"api-key-rotator/backend/internal/converters/formats"
)

// StreamHandler implements formats.StreamHandler for OpenAI format
type StreamHandler struct{}

func (h *StreamHandler) Name() string {
	return "openai"
}

// StreamChunk represents an OpenAI streaming chunk
type StreamChunk struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

type StreamChoice struct {
	Index        int          `json:"index"`
	Delta        *StreamDelta `json:"delta,omitempty"`
	FinishReason *string      `json:"finish_reason,omitempty"`
}

type StreamDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// ParseStreamChunk implements StreamHandler
func (h *StreamHandler) ParseStreamChunk(chunk []byte) (*formats.UniversalStreamChunk, error) {
	var streamChunk StreamChunk
	if err := json.Unmarshal(chunk, &streamChunk); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI stream chunk: %w", err)
	}

	universal := &formats.UniversalStreamChunk{
		ID:    streamChunk.ID,
		Model: streamChunk.Model,
	}

	if len(streamChunk.Choices) > 0 {
		choice := streamChunk.Choices[0]
		if choice.Delta != nil {
			universal.Delta = choice.Delta.Content
			universal.Role = choice.Delta.Role
		}
		if choice.FinishReason != nil && *choice.FinishReason != "" {
			universal.StopReason = choice.FinishReason
			universal.IsLast = true
		}
	}

	return universal, nil
}

// BuildStreamChunk implements StreamHandler
func (h *StreamHandler) BuildStreamChunk(chunk *formats.UniversalStreamChunk) ([]byte, error) {
	streamChunk := StreamChunk{
		ID:      chunk.ID,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   chunk.Model,
		Choices: []StreamChoice{
			{
				Index: 0,
				Delta: &StreamDelta{
					Content: chunk.Delta,
				},
			},
		},
	}

	if chunk.Role != "" {
		streamChunk.Choices[0].Delta.Role = chunk.Role
	}

	if chunk.StopReason != nil {
		streamChunk.Choices[0].FinishReason = chunk.StopReason
		streamChunk.Choices[0].Delta = &StreamDelta{}
	}

	return json.Marshal(streamChunk)
}

// BuildStartEvent implements StreamHandler
func (h *StreamHandler) BuildStartEvent(model string, id string) [][]byte {
	return nil
}

// BuildEndEvent implements StreamHandler
func (h *StreamHandler) BuildEndEvent() [][]byte {
	return [][]byte{[]byte("[DONE]")}
}
