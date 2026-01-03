package converters

import (
	"fmt"

	"api-key-rotator/backend/internal/converters/formats"
	// Import format packages to trigger their init() registration
	_ "api-key-rotator/backend/internal/converters/formats/anthropic"
	_ "api-key-rotator/backend/internal/converters/formats/gemini"
	_ "api-key-rotator/backend/internal/converters/formats/openai"
	_ "api-key-rotator/backend/internal/converters/formats/openai_responses"
)

// Converter handles format conversion between different LLM API formats
type Converter struct {
	from       formats.FormatHandler
	to         formats.FormatHandler
	fromStream formats.StreamHandler
	toStream   formats.StreamHandler
}

// NewConverter creates a new converter between two formats
// fromFormat and toFormat should be format names like "openai", "anthropic", "gemini"
func NewConverter(fromFormat, toFormat string) (*Converter, error) {
	// Normalize format names
	fromFormat = NormalizeFormat(fromFormat)
	toFormat = NormalizeFormat(toFormat)

	// Get handlers from registry
	fromInfo, err := formats.GetFormat(fromFormat)
	if err != nil {
		return nil, fmt.Errorf("source format error: %w", err)
	}

	toInfo, err := formats.GetFormat(toFormat)
	if err != nil {
		return nil, fmt.Errorf("target format error: %w", err)
	}

	return &Converter{
		from:       fromInfo.Handler,
		to:         toInfo.Handler,
		fromStream: fromInfo.StreamHandler,
		toStream:   toInfo.StreamHandler,
	}, nil
}

// ConvertRequest converts a request from source format to target format
func (c *Converter) ConvertRequest(body []byte) ([]byte, error) {
	// Parse source format to universal
	universal, err := c.from.ParseRequest(body)
	if err != nil {
		return nil, fmt.Errorf("parse request error: %w", err)
	}

	// Build target format from universal
	result, err := c.to.BuildRequest(universal)
	if err != nil {
		return nil, fmt.Errorf("build request error: %w", err)
	}

	return result, nil
}

// ConvertResponse converts a response from source format to target format
func (c *Converter) ConvertResponse(body []byte) ([]byte, error) {
	// Parse source format to universal
	universal, err := c.from.ParseResponse(body)
	if err != nil {
		return nil, fmt.Errorf("parse response error: %w", err)
	}

	// Build target format from universal
	result, err := c.to.BuildResponse(universal)
	if err != nil {
		return nil, fmt.Errorf("build response error: %w", err)
	}

	return result, nil
}

// ConvertStreamChunk converts a streaming chunk from source format to target format
func (c *Converter) ConvertStreamChunk(chunk []byte) ([]byte, error) {
	// Parse source format to universal
	universal, err := c.fromStream.ParseStreamChunk(chunk)
	if err != nil {
		return nil, fmt.Errorf("parse stream chunk error: %w", err)
	}

	// Skip empty chunks
	if universal.Delta == "" && universal.StopReason == nil && !universal.IsFirst && !universal.IsLast {
		return nil, nil
	}

	// Build target format from universal
	result, err := c.toStream.BuildStreamChunk(universal)
	if err != nil {
		return nil, fmt.Errorf("build stream chunk error: %w", err)
	}

	return result, nil
}

// GetTargetPath converts a client action path to the target API path
func (c *Converter) GetTargetPath(action string) string {
	// First convert from client format's perspective
	// Then convert to target format's API path
	return c.to.GetAPIPath(action)
}

// GetStreamStartEvents returns the start events needed for the target format
func (c *Converter) GetStreamStartEvents(model, id string) [][]byte {
	return c.toStream.BuildStartEvent(model, id)
}

// GetStreamEndEvents returns the end events needed for the target format
func (c *Converter) GetStreamEndEvents() [][]byte {
	return c.toStream.BuildEndEvent()
}

// NormalizeFormat converts api_format values to standard format names
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

// NeedsConversion checks if conversion is needed between two formats
func NeedsConversion(clientFormat, apiFormat string) bool {
	if clientFormat == "none" || clientFormat == "" {
		return false
	}
	normalizedAPI := NormalizeFormat(apiFormat)
	return clientFormat != normalizedAPI
}
