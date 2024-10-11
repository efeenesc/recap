package models

import (
	"fmt"
	"sync"
)

// TextVisionAPIFactory is a type alias for clarity
type TextVisionAPIFactory func(string) TextVisionAPI

var (
	registeredAPIs = make(map[string]TextVisionAPIFactory)
	APIList        []string
	apiMutex       sync.RWMutex
)

// RegisterAPI safely registers a new API factory
func RegisterAPI(name string, factory TextVisionAPIFactory) {
	apiMutex.Lock()
	defer apiMutex.Unlock()

	if factory == nil {
		panic("Attempted to register nil factory for " + name)
	}

	if _, exists := registeredAPIs[name]; exists {
		panic("Attempted to register duplicate API: " + name)
	}

	registeredAPIs[name] = factory
	APIList = append(APIList, name)
}

// GetAPI safely retrieves a registered API factory
func GetAPI(name string) (TextVisionAPIFactory, error) {
	apiMutex.RLock()
	defer apiMutex.RUnlock()

	factory, ok := registeredAPIs[name]
	if !ok {
		return nil, fmt.Errorf("API not found: %s", name)
	}

	return factory, nil
}

// ListRegisteredAPIs returns a list of all registered API names
func ListRegisteredAPIs() []string {
	apiMutex.RLock()
	defer apiMutex.RUnlock()

	return APIList
}
