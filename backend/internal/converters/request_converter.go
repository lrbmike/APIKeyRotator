package converters

import (
	"fmt"
)

// RequestConverter interface for request format conversion
type RequestConverter interface {
	// Convert transforms a request body from one format to another
	Convert(body []byte) ([]byte, error)
	// GetTargetPath returns the converted API path
	GetTargetPath(originalPath string) string
}

// PassthroughRequestConverter returns the request as-is without conversion
type PassthroughRequestConverter struct{}

func (c *PassthroughRequestConverter) Convert(body []byte) ([]byte, error) {
	return body, nil
}

func (c *PassthroughRequestConverter) GetTargetPath(originalPath string) string {
	return originalPath
}

// NewRequestConverter creates a converter based on client and API formats
// clientFormat: the format of the incoming request from the client (openai, anthropic, gemini)
// apiFormat: the format expected by the upstream API (openai_compatible, anthropic_native, gemini_native)
func NewRequestConverter(clientFormat, apiFormat string) RequestConverter {
	// Normalize API format to match client format naming
	normalizedAPI := NormalizeFormat(apiFormat)

	// If same format or no conversion needed, use passthrough
	if clientFormat == "none" || clientFormat == "" || clientFormat == normalizedAPI {
		return &PassthroughRequestConverter{}
	}

	// Select appropriate converter based on client → API combination
	switch {
	case clientFormat == "anthropic" && normalizedAPI == "openai":
		return &AnthropicToOpenAIRequestConverter{}
	case clientFormat == "openai" && normalizedAPI == "anthropic":
		return &OpenAIToAnthropicRequestConverter{}
	case clientFormat == "gemini" && normalizedAPI == "openai":
		return &GeminiToOpenAIRequestConverter{}
	case clientFormat == "openai" && normalizedAPI == "gemini":
		return &OpenAIToGeminiRequestConverter{}
	default:
		// Unsupported conversion, use passthrough
		return &PassthroughRequestConverter{}
	}
}

// AnthropicToOpenAIRequestConverter converts Anthropic request to OpenAI format
type AnthropicToOpenAIRequestConverter struct{}

func (c *AnthropicToOpenAIRequestConverter) Convert(body []byte) ([]byte, error) {
	var anthropicReq AnthropicRequest
	if err := parseJSON(body, &anthropicReq); err != nil {
		return nil, fmt.Errorf("failed to parse Anthropic request: %w", err)
	}

	openaiReq := c.convertRequest(&anthropicReq)
	return toJSON(openaiReq)
}

func (c *AnthropicToOpenAIRequestConverter) GetTargetPath(originalPath string) string {
	// Map Anthropic paths to OpenAI paths
	switch originalPath {
	case "messages", "v1/messages":
		return "chat/completions"
	default:
		return originalPath
	}
}

func (c *AnthropicToOpenAIRequestConverter) convertRequest(anthropic *AnthropicRequest) *OpenAIRequest {
	req := &OpenAIRequest{
		Model:    anthropic.Model,
		Messages: make([]OpenAIRequestMessage, 0),
		Stream:   anthropic.Stream,
	}

	// Convert max_tokens
	if anthropic.MaxTokens > 0 {
		req.MaxTokens = &anthropic.MaxTokens
	}

	// Convert temperature
	if anthropic.Temperature != nil {
		req.Temperature = anthropic.Temperature
	}

	// Convert top_p
	if anthropic.TopP != nil {
		req.TopP = anthropic.TopP
	}

	// Convert system message - can be string or array of content blocks
	if anthropic.System != nil {
		systemContent := extractSystemContent(anthropic.System)
		if systemContent != "" {
			req.Messages = append(req.Messages, OpenAIRequestMessage{
				Role:    "system",
				Content: systemContent,
			})
		}
	}

	// Convert messages
	for _, msg := range anthropic.Messages {
		openaiMsg := OpenAIRequestMessage{
			Role: msg.Role,
		}

		// Handle content - can be string or array of content blocks
		if msg.Content != nil {
			openaiMsg.Content = extractTextContent(msg.Content)
		}

		req.Messages = append(req.Messages, openaiMsg)
	}

	return req
}

// extractSystemContent extracts text from system field (can be string or array)
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

// extractTextContent extracts text from content field (can be string or array)
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

// OpenAIToAnthropicRequestConverter converts OpenAI request to Anthropic format
type OpenAIToAnthropicRequestConverter struct{}

func (c *OpenAIToAnthropicRequestConverter) Convert(body []byte) ([]byte, error) {
	var openaiReq OpenAIRequest
	if err := parseJSON(body, &openaiReq); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI request: %w", err)
	}

	anthropicReq := c.convertRequest(&openaiReq)
	return toJSON(anthropicReq)
}

func (c *OpenAIToAnthropicRequestConverter) GetTargetPath(originalPath string) string {
	switch originalPath {
	case "chat/completions", "v1/chat/completions":
		return "v1/messages"
	default:
		return originalPath
	}
}

func (c *OpenAIToAnthropicRequestConverter) convertRequest(openai *OpenAIRequest) *AnthropicRequest {
	req := &AnthropicRequest{
		Model:       openai.Model,
		Messages:    make([]AnthropicRequestMessage, 0),
		Stream:      openai.Stream,
		Temperature: openai.Temperature,
		TopP:        openai.TopP,
	}

	if openai.MaxTokens != nil {
		req.MaxTokens = *openai.MaxTokens
	} else {
		req.MaxTokens = 4096 // Default for Anthropic
	}

	// Convert messages
	for _, msg := range openai.Messages {
		if msg.Role == "system" {
			// Extract system message
			if content, ok := msg.Content.(string); ok {
				req.System = content
			}
			continue
		}

		anthropicMsg := AnthropicRequestMessage{
			Role: msg.Role,
		}

		// Convert content
		if content, ok := msg.Content.(string); ok {
			anthropicMsg.Content = content
		}

		req.Messages = append(req.Messages, anthropicMsg)
	}

	return req
}

// GeminiToOpenAIRequestConverter converts Gemini request to OpenAI format
type GeminiToOpenAIRequestConverter struct{}

func (c *GeminiToOpenAIRequestConverter) Convert(body []byte) ([]byte, error) {
	var geminiReq GeminiRequest
	if err := parseJSON(body, &geminiReq); err != nil {
		return nil, fmt.Errorf("failed to parse Gemini request: %w", err)
	}

	openaiReq := c.convertRequest(&geminiReq)
	return toJSON(openaiReq)
}

func (c *GeminiToOpenAIRequestConverter) GetTargetPath(originalPath string) string {
	// Gemini uses model-specific paths, convert to OpenAI chat/completions
	return "chat/completions"
}

func (c *GeminiToOpenAIRequestConverter) convertRequest(gemini *GeminiRequest) *OpenAIRequest {
	req := &OpenAIRequest{
		Messages: make([]OpenAIRequestMessage, 0),
	}

	// Extract model from path if not in body (Gemini includes model in URL)
	// This will be handled separately

	// Convert contents to messages
	for _, content := range gemini.Contents {
		role := content.Role
		if role == "model" {
			role = "assistant"
		}

		var textContent string
		for _, part := range content.Parts {
			textContent += part.Text
		}

		req.Messages = append(req.Messages, OpenAIRequestMessage{
			Role:    role,
			Content: textContent,
		})
	}

	// Convert generation config
	if gemini.GenerationConfig != nil {
		if gemini.GenerationConfig.MaxOutputTokens > 0 {
			req.MaxTokens = &gemini.GenerationConfig.MaxOutputTokens
		}
		if gemini.GenerationConfig.Temperature != nil {
			req.Temperature = gemini.GenerationConfig.Temperature
		}
		if gemini.GenerationConfig.TopP != nil {
			req.TopP = gemini.GenerationConfig.TopP
		}
	}

	return req
}

// OpenAIToGeminiRequestConverter converts OpenAI request to Gemini format
type OpenAIToGeminiRequestConverter struct{}

func (c *OpenAIToGeminiRequestConverter) Convert(body []byte) ([]byte, error) {
	var openaiReq OpenAIRequest
	if err := parseJSON(body, &openaiReq); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI request: %w", err)
	}

	geminiReq := c.convertRequest(&openaiReq)
	return toJSON(geminiReq)
}

func (c *OpenAIToGeminiRequestConverter) GetTargetPath(originalPath string) string {
	// Gemini uses a different path structure
	return originalPath
}

func (c *OpenAIToGeminiRequestConverter) convertRequest(openai *OpenAIRequest) *GeminiRequest {
	req := &GeminiRequest{
		Contents: make([]GeminiRequestContent, 0),
	}

	// Convert messages to contents
	for _, msg := range openai.Messages {
		if msg.Role == "system" {
			// Gemini handles system differently via systemInstruction
			continue
		}

		role := msg.Role
		if role == "assistant" {
			role = "model"
		}

		var text string
		if content, ok := msg.Content.(string); ok {
			text = content
		}

		req.Contents = append(req.Contents, GeminiRequestContent{
			Role: role,
			Parts: []GeminiPart{
				{Text: text},
			},
		})
	}

	// Convert generation config
	req.GenerationConfig = &GeminiGenerationConfig{}
	if openai.MaxTokens != nil {
		req.GenerationConfig.MaxOutputTokens = *openai.MaxTokens
	}
	if openai.Temperature != nil {
		req.GenerationConfig.Temperature = openai.Temperature
	}
	if openai.TopP != nil {
		req.GenerationConfig.TopP = openai.TopP
	}

	return req
}
