package converters

// Anthropic Messages API Response structures
type AnthropicResponse struct {
	ID           string             `json:"id"`
	Type         string             `json:"type"` // "message"
	Role         string             `json:"role"` // "assistant"
	Content      []AnthropicContent `json:"content"`
	Model        string             `json:"model"`
	StopReason   *string            `json:"stop_reason,omitempty"`
	StopSequence *string            `json:"stop_sequence,omitempty"`
	Usage        *AnthropicUsage    `json:"usage,omitempty"`
}

type AnthropicContent struct {
	Type string `json:"type"` // "text"
	Text string `json:"text"`
}

type AnthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

// Anthropic Streaming Event structures
type AnthropicStreamEvent struct {
	Type         string                `json:"type"`
	Message      *AnthropicResponse    `json:"message,omitempty"`
	Index        *int                  `json:"index,omitempty"`
	ContentBlock *AnthropicContent     `json:"content_block,omitempty"`
	Delta        *AnthropicStreamDelta `json:"delta,omitempty"`
	Usage        *AnthropicStreamUsage `json:"usage,omitempty"`
}

type AnthropicStreamDelta struct {
	Type       string  `json:"type,omitempty"`
	Text       string  `json:"text,omitempty"`
	StopReason *string `json:"stop_reason,omitempty"`
}

type AnthropicStreamUsage struct {
	OutputTokens int `json:"output_tokens"`
}

// Anthropic Request structures
type AnthropicRequest struct {
	Model       string                    `json:"model"`
	Messages    []AnthropicRequestMessage `json:"messages"`
	System      interface{}               `json:"system,omitempty"` // Can be string or array of content blocks
	MaxTokens   int                       `json:"max_tokens"`
	Temperature *float64                  `json:"temperature,omitempty"`
	TopP        *float64                  `json:"top_p,omitempty"`
	Stream      bool                      `json:"stream,omitempty"`
}

type AnthropicRequestMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"` // Can be string or array of content blocks
}
