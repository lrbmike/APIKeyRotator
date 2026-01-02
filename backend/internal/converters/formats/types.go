package formats

// Universal message types for format-agnostic conversion
// All format handlers convert to/from these types

// UniversalRequest represents a chat completion request in a format-agnostic way
type UniversalRequest struct {
	Model       string             `json:"model"`
	Messages    []UniversalMessage `json:"messages"`
	System      string             `json:"system,omitempty"`
	MaxTokens   int                `json:"max_tokens,omitempty"`
	Temperature *float64           `json:"temperature,omitempty"`
	TopP        *float64           `json:"top_p,omitempty"`
	Stream      bool               `json:"stream,omitempty"`
	Stop        []string           `json:"stop,omitempty"`
}

// UniversalMessage represents a single message in a conversation
type UniversalMessage struct {
	Role    string `json:"role"`    // "system", "user", "assistant"
	Content string `json:"content"` // Text content
}

// UniversalResponse represents a chat completion response in a format-agnostic way
type UniversalResponse struct {
	ID         string          `json:"id"`
	Model      string          `json:"model"`
	Content    string          `json:"content"`
	Role       string          `json:"role"`
	StopReason string          `json:"stop_reason,omitempty"`
	Usage      *UniversalUsage `json:"usage,omitempty"`
}

// UniversalUsage represents token usage statistics
type UniversalUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
	TotalTokens  int `json:"total_tokens"`
}

// UniversalStreamChunk represents a single streaming chunk in a format-agnostic way
type UniversalStreamChunk struct {
	ID         string  `json:"id,omitempty"`
	Model      string  `json:"model,omitempty"`
	Delta      string  `json:"delta"`          // The text delta
	Role       string  `json:"role,omitempty"` // Role (usually only in first chunk)
	StopReason *string `json:"stop_reason,omitempty"`
	IsFirst    bool    `json:"-"` // Internal flag for first chunk
	IsLast     bool    `json:"-"` // Internal flag for last chunk
}
