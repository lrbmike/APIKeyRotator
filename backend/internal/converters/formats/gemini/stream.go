package gemini

import (
	"encoding/json"
	"fmt"

	"api-key-rotator/backend/internal/converters/formats"
)

// StreamHandler implements formats.StreamHandler for Gemini format
type StreamHandler struct{}

func (h *StreamHandler) Name() string {
	return "gemini"
}

// ParseStreamChunk implements StreamHandler
func (h *StreamHandler) ParseStreamChunk(chunk []byte) (*formats.UniversalStreamChunk, error) {
	var resp Response
	if err := json.Unmarshal(chunk, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini stream chunk: %w", err)
	}

	universal := &formats.UniversalStreamChunk{}

	if len(resp.Candidates) > 0 {
		candidate := resp.Candidates[0]
		if candidate.Content != nil {
			for _, part := range candidate.Content.Parts {
				universal.Delta += part.Text
			}
			role := candidate.Content.Role
			if role == "model" {
				role = "assistant"
			}
			universal.Role = role
		}

		if candidate.FinishReason != "" {
			reason := mapGeminiFinishReason(candidate.FinishReason)
			universal.StopReason = &reason
			universal.IsLast = true
		}
	}

	return universal, nil
}

// BuildStreamChunk implements StreamHandler
func (h *StreamHandler) BuildStreamChunk(chunk *formats.UniversalStreamChunk) ([]byte, error) {
	role := chunk.Role
	if role == "assistant" || role == "" {
		role = "model"
	}

	geminiChunk := Response{
		Candidates: []Candidate{
			{
				Index: 0,
				Content: &ResponseContent{
					Role:  role,
					Parts: []Part{{Text: chunk.Delta}},
				},
			},
		},
	}

	if chunk.StopReason != nil {
		geminiChunk.Candidates[0].FinishReason = mapToGeminiFinishReason(*chunk.StopReason)
	}

	return json.Marshal(geminiChunk)
}

// BuildStartEvent implements StreamHandler
func (h *StreamHandler) BuildStartEvent(model string, id string) [][]byte {
	return nil
}

// BuildEndEvent implements StreamHandler
func (h *StreamHandler) BuildEndEvent() [][]byte {
	return nil
}
