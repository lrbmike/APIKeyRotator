package converters

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ResponseConverter interface for format conversion
type ResponseConverter interface {
	// Convert transforms a complete (non-streaming) response body
	Convert(body []byte) ([]byte, error)
	// ConvertStreamChunk transforms a single SSE chunk's JSON payload
	ConvertStreamChunk(chunk []byte) ([]byte, error)
	// GetContentType returns the content type for the converted response
	GetContentType() string
}

// PassthroughConverter returns the response as-is without conversion
type PassthroughConverter struct{}

func (c *PassthroughConverter) Convert(body []byte) ([]byte, error) {
	return body, nil
}

func (c *PassthroughConverter) ConvertStreamChunk(chunk []byte) ([]byte, error) {
	return chunk, nil
}

func (c *PassthroughConverter) GetContentType() string {
	return "application/json"
}

// NewResponseConverter creates a converter based on input and output formats
// inputFormat: the format of the upstream API response (openai_compatible, anthropic_native, gemini_native)
// outputFormat: the desired output format (none, openai, anthropic, gemini)
func NewResponseConverter(inputFormat, outputFormat string) ResponseConverter {
	// Normalize input format to match output format naming
	normalizedInput := NormalizeFormat(inputFormat)

	// If no conversion needed or same format, use passthrough
	if outputFormat == "none" || outputFormat == "" || normalizedInput == outputFormat {
		return &PassthroughConverter{}
	}

	// Select appropriate converter based on input → output combination
	switch {
	case normalizedInput == "openai" && outputFormat == "anthropic":
		return &OpenAIToAnthropicConverter{}
	case normalizedInput == "openai" && outputFormat == "gemini":
		return &OpenAIToGeminiConverter{}
	case normalizedInput == "anthropic" && outputFormat == "openai":
		return &AnthropicToOpenAIConverter{}
	case normalizedInput == "gemini" && outputFormat == "openai":
		return &GeminiToOpenAIConverter{}
	case normalizedInput == "anthropic" && outputFormat == "gemini":
		// anthropic → gemini: chain through openai
		return &ChainedConverter{
			first:  &AnthropicToOpenAIConverter{},
			second: &OpenAIToGeminiConverter{},
		}
	case normalizedInput == "gemini" && outputFormat == "anthropic":
		// gemini → anthropic: chain through openai
		return &ChainedConverter{
			first:  &GeminiToOpenAIConverter{},
			second: &OpenAIToAnthropicConverter{},
		}
	default:
		// Unsupported conversion, use passthrough
		return &PassthroughConverter{}
	}
}

// NormalizeFormat converts api_format values to client_format naming convention
func NormalizeFormat(format string) string {
	switch format {
	case "openai_compatible":
		return "openai"
	case "anthropic_native":
		return "anthropic"
	case "gemini_native":
		return "gemini"
	default:
		return format
	}
}

// ChainedConverter chains two converters together
type ChainedConverter struct {
	first  ResponseConverter
	second ResponseConverter
}

func (c *ChainedConverter) Convert(body []byte) ([]byte, error) {
	intermediate, err := c.first.Convert(body)
	if err != nil {
		return nil, err
	}
	return c.second.Convert(intermediate)
}

func (c *ChainedConverter) ConvertStreamChunk(chunk []byte) ([]byte, error) {
	intermediate, err := c.first.ConvertStreamChunk(chunk)
	if err != nil {
		return nil, err
	}
	return c.second.ConvertStreamChunk(intermediate)
}

func (c *ChainedConverter) GetContentType() string {
	return c.second.GetContentType()
}

// Helper function to parse SSE data line
func ParseSSEChunk(data []byte) ([]byte, bool) {
	str := strings.TrimSpace(string(data))
	if strings.HasPrefix(str, "data: ") {
		payload := strings.TrimPrefix(str, "data: ")
		if payload == "[DONE]" {
			return nil, false
		}
		return []byte(payload), true
	}
	return nil, false
}

// Helper function to format SSE data line
func FormatSSEChunk(data []byte) []byte {
	return []byte(fmt.Sprintf("data: %s\n\n", string(data)))
}

// Common JSON helper
func parseJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

func toJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
