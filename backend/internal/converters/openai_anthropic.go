package converters

import (
	"fmt"
	"time"
)

// OpenAIToAnthropicConverter converts OpenAI responses to Anthropic format
type OpenAIToAnthropicConverter struct{}

func (c *OpenAIToAnthropicConverter) Convert(body []byte) ([]byte, error) {
	var openaiResp OpenAIResponse
	if err := parseJSON(body, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	anthropicResp := c.convertResponse(&openaiResp)
	return toJSON(anthropicResp)
}

func (c *OpenAIToAnthropicConverter) ConvertStreamChunk(chunk []byte) ([]byte, error) {
	var openaiChunk OpenAIStreamChunk
	if err := parseJSON(chunk, &openaiChunk); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI stream chunk: %w", err)
	}

	// Convert to Anthropic stream event
	anthropicEvent := c.convertStreamChunk(&openaiChunk)
	return toJSON(anthropicEvent)
}

func (c *OpenAIToAnthropicConverter) GetContentType() string {
	return "application/json"
}

func (c *OpenAIToAnthropicConverter) convertResponse(openai *OpenAIResponse) *AnthropicResponse {
	// Default stop_reason
	defaultStopReason := "end_turn"

	resp := &AnthropicResponse{
		ID:         openai.ID,
		Type:       "message",
		Role:       "assistant",
		Model:      openai.Model,
		Content:    []AnthropicContent{}, // Initialize empty array
		StopReason: &defaultStopReason,   // Default stop reason
	}

	// Convert content from choices
	if len(openai.Choices) > 0 {
		choice := openai.Choices[0]
		if choice.Message != nil && choice.Message.Content != "" {
			resp.Content = []AnthropicContent{
				{Type: "text", Text: choice.Message.Content},
			}
		}
		if choice.FinishReason != nil {
			stopReason := mapFinishReasonToStopReason(*choice.FinishReason)
			resp.StopReason = &stopReason
		}
	}

	// Convert usage - always include even if zero
	if openai.Usage != nil {
		resp.Usage = &AnthropicUsage{
			InputTokens:  openai.Usage.PromptTokens,
			OutputTokens: openai.Usage.CompletionTokens,
		}
	} else {
		resp.Usage = &AnthropicUsage{
			InputTokens:  0,
			OutputTokens: 0,
		}
	}

	return resp
}

func (c *OpenAIToAnthropicConverter) convertStreamChunk(openai *OpenAIStreamChunk) *AnthropicStreamEvent {
	if len(openai.Choices) == 0 {
		return &AnthropicStreamEvent{Type: "ping"}
	}

	choice := openai.Choices[0]

	// Check for finish
	if choice.FinishReason != nil && *choice.FinishReason != "" {
		stopReason := mapFinishReasonToStopReason(*choice.FinishReason)
		return &AnthropicStreamEvent{
			Type: "message_delta",
			Delta: &AnthropicStreamDelta{
				StopReason: &stopReason,
			},
		}
	}

	// Content delta
	if choice.Delta != nil && choice.Delta.Content != "" {
		return &AnthropicStreamEvent{
			Type:  "content_block_delta",
			Index: intPtr(0),
			Delta: &AnthropicStreamDelta{
				Type: "text_delta",
				Text: choice.Delta.Content,
			},
		}
	}

	return &AnthropicStreamEvent{Type: "ping"}
}

// mapFinishReasonToStopReason maps OpenAI finish_reason to Anthropic stop_reason
func mapFinishReasonToStopReason(finishReason string) string {
	switch finishReason {
	case "stop":
		return "end_turn"
	case "length":
		return "max_tokens"
	case "content_filter":
		return "stop_sequence"
	default:
		return "end_turn"
	}
}

// AnthropicToOpenAIConverter converts Anthropic responses to OpenAI format
type AnthropicToOpenAIConverter struct{}

func (c *AnthropicToOpenAIConverter) Convert(body []byte) ([]byte, error) {
	var anthropicResp AnthropicResponse
	if err := parseJSON(body, &anthropicResp); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic response: %w", err)
	}

	openaiResp := c.convertResponse(&anthropicResp)
	return toJSON(openaiResp)
}

func (c *AnthropicToOpenAIConverter) ConvertStreamChunk(chunk []byte) ([]byte, error) {
	var anthropicEvent AnthropicStreamEvent
	if err := parseJSON(chunk, &anthropicEvent); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic stream event: %w", err)
	}

	openaiChunk := c.convertStreamChunk(&anthropicEvent)
	if openaiChunk == nil {
		return nil, nil // Skip this event
	}
	return toJSON(openaiChunk)
}

func (c *AnthropicToOpenAIConverter) GetContentType() string {
	return "application/json"
}

func (c *AnthropicToOpenAIConverter) convertResponse(anthropic *AnthropicResponse) *OpenAIResponse {
	resp := &OpenAIResponse{
		ID:      anthropic.ID,
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   anthropic.Model,
		Choices: []OpenAIChoice{
			{
				Index: 0,
				Message: &OpenAIMessage{
					Role:    "assistant",
					Content: "",
				},
			},
		},
	}

	// Extract content
	if len(anthropic.Content) > 0 {
		var contentText string
		for _, content := range anthropic.Content {
			if content.Type == "text" {
				contentText += content.Text
			}
		}
		resp.Choices[0].Message.Content = contentText
	}

	// Map stop reason to finish reason
	if anthropic.StopReason != nil {
		finishReason := mapStopReasonToFinishReason(*anthropic.StopReason)
		resp.Choices[0].FinishReason = &finishReason
	}

	// Convert usage
	if anthropic.Usage != nil {
		resp.Usage = &OpenAIUsage{
			PromptTokens:     anthropic.Usage.InputTokens,
			CompletionTokens: anthropic.Usage.OutputTokens,
			TotalTokens:      anthropic.Usage.InputTokens + anthropic.Usage.OutputTokens,
		}
	}

	return resp
}

func (c *AnthropicToOpenAIConverter) convertStreamChunk(anthropic *AnthropicStreamEvent) *OpenAIStreamChunk {
	switch anthropic.Type {
	case "content_block_delta":
		if anthropic.Delta != nil && anthropic.Delta.Text != "" {
			return &OpenAIStreamChunk{
				ID:      "chatcmpl-stream",
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Choices: []OpenAIChoice{
					{
						Index: 0,
						Delta: &OpenAIMessage{
							Content: anthropic.Delta.Text,
						},
					},
				},
			}
		}
	case "message_delta":
		if anthropic.Delta != nil && anthropic.Delta.StopReason != nil {
			finishReason := mapStopReasonToFinishReason(*anthropic.Delta.StopReason)
			return &OpenAIStreamChunk{
				ID:      "chatcmpl-stream",
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Choices: []OpenAIChoice{
					{
						Index:        0,
						Delta:        &OpenAIMessage{},
						FinishReason: &finishReason,
					},
				},
			}
		}
	}
	return nil
}

// mapStopReasonToFinishReason maps Anthropic stop_reason to OpenAI finish_reason
func mapStopReasonToFinishReason(stopReason string) string {
	switch stopReason {
	case "end_turn":
		return "stop"
	case "max_tokens":
		return "length"
	case "stop_sequence":
		return "stop"
	default:
		return "stop"
	}
}

// Helper function
func intPtr(i int) *int {
	return &i
}
