package formats

import (
	"fmt"
	"sync"
)

// Global format registry
var (
	registry = make(map[string]*FormatInfo)
	mu       sync.RWMutex
)

// RegisterFormat registers a format handler with the given name
func RegisterFormat(name string, handler FormatHandler, streamHandler StreamHandler) {
	mu.Lock()
	defer mu.Unlock()
	registry[name] = &FormatInfo{
		Handler:       handler,
		StreamHandler: streamHandler,
	}
}

// GetFormat returns the format info for the given name
func GetFormat(name string) (*FormatInfo, error) {
	mu.RLock()
	defer mu.RUnlock()
	info, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("unknown format: %s", name)
	}
	return info, nil
}

// GetHandler returns the FormatHandler for the given format name
func GetHandler(name string) (FormatHandler, error) {
	info, err := GetFormat(name)
	if err != nil {
		return nil, err
	}
	return info.Handler, nil
}

// GetStreamHandler returns the StreamHandler for the given format name
func GetStreamHandler(name string) (StreamHandler, error) {
	info, err := GetFormat(name)
	if err != nil {
		return nil, err
	}
	return info.StreamHandler, nil
}

// ListFormats returns a list of all registered format names
func ListFormats() []string {
	mu.RLock()
	defer mu.RUnlock()
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
