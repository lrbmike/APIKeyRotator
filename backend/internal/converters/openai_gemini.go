package converters

import (
	"fmt"
	"time"
)

// OpenAIToGeminiConverter converts OpenAI responses to Gemini format
type OpenAIToGeminiConverter struct{}

func (c *OpenAIToGeminiConverter) Convert(body []byte) ([]byte, error) {
	var openaiResp OpenAIResponse
	if err := parseJSON(body, &openaiResp); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w", err)
	}

	geminiResp := c.convertResponse(&openaiResp)
	return toJSON(geminiResp)
}

func (c *OpenAIToGeminiConverter) ConvertStreamChunk(chunk []byte) ([]byte, error) {
	var openaiChunk OpenAIStreamChunk
	if err := parseJSON(chunk, &openaiChunk); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI stream chunk: %w", err)
	}

	geminiChunk := c.convertStreamChunk(&openaiChunk)
	return toJSON(geminiChunk)
}

func (c *OpenAIToGeminiConverter) GetContentType() string {
	return "application/json"
}

func (c *OpenAIToGeminiConverter) convertResponse(openai *OpenAIResponse) *GeminiResponse {
	resp := &GeminiResponse{
		Candidates:   make([]GeminiCandidate, 0),
		ModelVersion: openai.Model,
	}

	if len(openai.Choices) > 0 {
		choice := openai.Choices[0]
		candidate := GeminiCandidate{
			Index:        choice.Index,
			FinishReason: mapOpenAIFinishReasonToGemini(choice.FinishReason),
		}

		if choice.Message != nil {
			candidate.Content = &GeminiContent{
				Role: "model",
				Parts: []GeminiPart{
					{Text: choice.Message.Content},
				},
			}
		}

		resp.Candidates = append(resp.Candidates, candidate)
	}

	if openai.Usage != nil {
		resp.UsageMetadata = &GeminiUsageMetadata{
			PromptTokenCount:     openai.Usage.PromptTokens,
			CandidatesTokenCount: openai.Usage.CompletionTokens,
			TotalTokenCount:      openai.Usage.TotalTokens,
		}
	}

	return resp
}

func (c *OpenAIToGeminiConverter) convertStreamChunk(openai *OpenAIStreamChunk) *GeminiStreamChunk {
	chunk := &GeminiStreamChunk{
		Candidates: make([]GeminiCandidate, 0),
	}

	if len(openai.Choices) > 0 {
		choice := openai.Choices[0]
		candidate := GeminiCandidate{
			Index:        choice.Index,
			FinishReason: mapOpenAIFinishReasonToGemini(choice.FinishReason),
		}

		if choice.Delta != nil && choice.Delta.Content != "" {
			candidate.Content = &GeminiContent{
				Role: "model",
				Parts: []GeminiPart{
					{Text: choice.Delta.Content},
				},
			}
		}

		chunk.Candidates = append(chunk.Candidates, candidate)
	}

	return chunk
}

func mapOpenAIFinishReasonToGemini(finishReason *string) string {
	if finishReason == nil {
		return ""
	}
	switch *finishReason {
	case "stop":
		return "STOP"
	case "length":
		return "MAX_TOKENS"
	case "content_filter":
		return "SAFETY"
	default:
		return "STOP"
	}
}

// GeminiToOpenAIConverter converts Gemini responses to OpenAI format
type GeminiToOpenAIConverter struct{}

func (c *GeminiToOpenAIConverter) Convert(body []byte) ([]byte, error) {
	var geminiResp GeminiResponse
	if err := parseJSON(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini response: %w", err)
	}

	openaiResp := c.convertResponse(&geminiResp)
	return toJSON(openaiResp)
}

func (c *GeminiToOpenAIConverter) ConvertStreamChunk(chunk []byte) ([]byte, error) {
	var geminiChunk GeminiStreamChunk
	if err := parseJSON(chunk, &geminiChunk); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini stream chunk: %w", err)
	}

	openaiChunk := c.convertStreamChunk(&geminiChunk)
	return toJSON(openaiChunk)
}

func (c *GeminiToOpenAIConverter) GetContentType() string {
	return "application/json"
}

func (c *GeminiToOpenAIConverter) convertResponse(gemini *GeminiResponse) *OpenAIResponse {
	resp := &OpenAIResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   gemini.ModelVersion,
		Choices: make([]OpenAIChoice, 0),
	}

	if len(gemini.Candidates) > 0 {
		candidate := gemini.Candidates[0]
		choice := OpenAIChoice{
			Index: candidate.Index,
			Message: &OpenAIMessage{
				Role: "assistant",
			},
		}

		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			var content string
			for _, part := range candidate.Content.Parts {
				content += part.Text
			}
			choice.Message.Content = content
		}

		finishReason := mapGeminiFinishReasonToOpenAI(candidate.FinishReason)
		choice.FinishReason = &finishReason

		resp.Choices = append(resp.Choices, choice)
	}

	if gemini.UsageMetadata != nil {
		resp.Usage = &OpenAIUsage{
			PromptTokens:     gemini.UsageMetadata.PromptTokenCount,
			CompletionTokens: gemini.UsageMetadata.CandidatesTokenCount,
			TotalTokens:      gemini.UsageMetadata.TotalTokenCount,
		}
	}

	return resp
}

func (c *GeminiToOpenAIConverter) convertStreamChunk(gemini *GeminiStreamChunk) *OpenAIStreamChunk {
	chunk := &OpenAIStreamChunk{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().UnixNano()),
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Choices: make([]OpenAIChoice, 0),
	}

	if len(gemini.Candidates) > 0 {
		candidate := gemini.Candidates[0]
		choice := OpenAIChoice{
			Index: candidate.Index,
			Delta: &OpenAIMessage{},
		}

		if candidate.Content != nil && len(candidate.Content.Parts) > 0 {
			var content string
			for _, part := range candidate.Content.Parts {
				content += part.Text
			}
			choice.Delta.Content = content
		}

		if candidate.FinishReason != "" {
			finishReason := mapGeminiFinishReasonToOpenAI(candidate.FinishReason)
			choice.FinishReason = &finishReason
		}

		chunk.Choices = append(chunk.Choices, choice)
	}

	return chunk
}

func mapGeminiFinishReasonToOpenAI(finishReason string) string {
	switch finishReason {
	case "STOP":
		return "stop"
	case "MAX_TOKENS":
		return "length"
	case "SAFETY":
		return "content_filter"
	default:
		return "stop"
	}
}
