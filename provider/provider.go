// Package provider implements the registry and base types for Gondolia's provider system.
//
// This package provides a type-safe, compile-time provider pattern for external integrations.
// Providers are registered via init() and discovered through the global registry.
package provider

import (
	"fmt"
	"sync"
)

// Metadata describes a provider implementation.
type Metadata struct {
	Name        string        // e.g. "saferpay"
	DisplayName string        // e.g. "Saferpay (SIX Payment Services)"
	Category    string        // e.g. "payment"
	Version     string        // e.g. "1.0.0"
	Description string        // Human-readable description
	ConfigSpec  []ConfigField // Configuration schema
}

// ConfigField describes a configuration field required by the provider.
type ConfigField struct {
	Key         string // Field name
	Type        string // "string", "int", "bool", "secret"
	Required    bool   // Is this field required?
	Default     any    // Default value (if any)
	Description string // Human-readable description
}

// ProviderFactory creates a provider instance from configuration.
type ProviderFactory[T any] func(config map[string]any) (T, error)

// --- Global Registry ---

var (
	registry = make(map[string]map[string]any)      // category -> name -> factory
	metadata = make(map[string]map[string]Metadata) // category -> name -> metadata
	mu       sync.RWMutex
)

// Register registers a provider factory with the global registry.
// This function is typically called from init() functions in provider packages.
//
// Example:
//
//	func init() {
//	    provider.Register[payment.PaymentProvider]("payment", "saferpay",
//	        provider.Metadata{Name: "saferpay", ...},
//	        NewProvider,
//	    )
//	}
func Register[T any](category, name string, meta Metadata, factory ProviderFactory[T]) {
	mu.Lock()
	defer mu.Unlock()

	if registry[category] == nil {
		registry[category] = make(map[string]any)
		metadata[category] = make(map[string]Metadata)
	}

	// Prevent duplicate registrations
	if _, exists := registry[category][name]; exists {
		panic(fmt.Sprintf("provider %s.%s is already registered", category, name))
	}

	registry[category][name] = factory
	metadata[category][name] = meta
}

// Get retrieves a provider factory from the registry.
//
// Example:
//
//	factory, err := provider.Get[payment.PaymentProvider]("payment", "saferpay")
//	if err != nil {
//	    return err
//	}
//	provider, err := factory(config)
func Get[T any](category, name string) (ProviderFactory[T], error) {
	mu.RLock()
	defer mu.RUnlock()

	cat, ok := registry[category]
	if !ok {
		return nil, fmt.Errorf("unknown provider category: %s", category)
	}

	factory, ok := cat[name]
	if !ok {
		return nil, fmt.Errorf("unknown provider: %s.%s", category, name)
	}

	f, ok := factory.(ProviderFactory[T])
	if !ok {
		return nil, fmt.Errorf("provider %s.%s has wrong type", category, name)
	}

	return f, nil
}

// List returns all registered providers in a category.
func List(category string) []Metadata {
	mu.RLock()
	defer mu.RUnlock()

	var result []Metadata
	for _, m := range metadata[category] {
		result = append(result, m)
	}
	return result
}

// ListAll returns all registered providers grouped by category.
func ListAll() map[string][]Metadata {
	mu.RLock()
	defer mu.RUnlock()

	result := make(map[string][]Metadata)
	for cat, providers := range metadata {
		for _, m := range providers {
			result[cat] = append(result[cat], m)
		}
	}
	return result
}
