package formats

// FormatHandler defines the interface for handling request/response format conversion
// Each format (OpenAI, Anthropic, Gemini) implements this interface
type FormatHandler interface {
	// Name returns the format identifier (e.g., "openai", "anthropic", "gemini")
	Name() string

	// Request handling
	// ParseRequest parses a format-specific request body into UniversalRequest
	ParseRequest(body []byte) (*UniversalRequest, error)
	// BuildRequest builds a format-specific request body from UniversalRequest
	BuildRequest(req *UniversalRequest) ([]byte, error)

	// Response handling
	// ParseResponse parses a format-specific response body into UniversalResponse
	ParseResponse(body []byte) (*UniversalResponse, error)
	// BuildResponse builds a format-specific response body from UniversalResponse
	BuildResponse(resp *UniversalResponse) ([]byte, error)

	// Path mapping
	// GetAPIPath converts a client action to the format's API path
	GetAPIPath(action string) string
	// GetClientAction converts an API path to a client action
	GetClientAction(apiPath string) string
}

// StreamHandler defines the interface for handling streaming format conversion
// Separated from FormatHandler since streaming has unique requirements
type StreamHandler interface {
	// Name returns the format identifier
	Name() string

	// ParseStreamChunk parses a format-specific stream chunk into UniversalStreamChunk
	ParseStreamChunk(chunk []byte) (*UniversalStreamChunk, error)
	// BuildStreamChunk builds a format-specific stream chunk from UniversalStreamChunk
	BuildStreamChunk(chunk *UniversalStreamChunk) ([]byte, error)

	// BuildStartEvent builds the initial event(s) for streaming (some formats need this)
	BuildStartEvent(model string, id string) [][]byte
	// BuildEndEvent builds the final event(s) for streaming
	BuildEndEvent() [][]byte
}

// FormatInfo contains both handlers for a format
type FormatInfo struct {
	Handler       FormatHandler
	StreamHandler StreamHandler
}
