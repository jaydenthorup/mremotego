package secrets

import (
	"fmt"
	"io"
	"sync"
)

// Registry holds the configured secret providers and routes references to the
// provider that understands them.
type Registry struct {
	providers []Provider
}

// NewRegistry creates a registry from the given providers.
func NewRegistry(providers ...Provider) *Registry {
	return &Registry{providers: providers}
}

// Providers returns every registered provider.
func (r *Registry) Providers() []Provider {
	return r.providers
}

// ByScheme returns the provider registered for a scheme such as "op" or "bw".
func (r *Registry) ByScheme(scheme string) (Provider, bool) {
	for _, p := range r.providers {
		if p.Scheme() == scheme {
			return p, true
		}
	}
	return nil, false
}

// ProviderFor returns the provider that owns the given reference.
func (r *Registry) ProviderFor(value string) (Provider, bool) {
	for _, p := range r.providers {
		if p.IsReference(value) {
			return p, true
		}
	}
	return nil, false
}

// IsReference reports whether value is a reference for any registered provider.
func (r *Registry) IsReference(value string) bool {
	_, ok := r.ProviderFor(value)
	return ok
}

// Resolve resolves a reference. Values that are not references are returned
// unchanged, which lets callers pass plain passwords through.
func (r *Registry) Resolve(value string) (string, error) {
	provider, ok := r.ProviderFor(value)
	if !ok {
		return value, nil
	}

	secret, err := provider.ResolveSecret(value)
	if err != nil {
		return "", fmt.Errorf("%s: %w", provider.Name(), err)
	}
	return secret, nil
}

// Close releases resources held by providers, such as helper processes.
func (r *Registry) Close() {
	for _, p := range r.providers {
		if closer, ok := p.(io.Closer); ok {
			_ = closer.Close()
		}
	}
}

var (
	defaultOnce     sync.Once
	defaultMu       sync.Mutex
	defaultRegistry *Registry
)

// Default returns the process wide registry. Providers are shared because some
// of them own a helper process (see BitwardenProvider) that must exist only
// once per application run.
func Default() *Registry {
	defaultOnce.Do(func() {
		registry := NewRegistry(
			NewOnePasswordProvider(),
			NewBitwardenProvider(),
		)
		defaultMu.Lock()
		defaultRegistry = registry
		defaultMu.Unlock()
	})

	defaultMu.Lock()
	defer defaultMu.Unlock()
	return defaultRegistry
}

// Shutdown closes the default registry if it was ever created. It is safe to
// call multiple times and must be called before the process exits so that
// helper processes are terminated.
func Shutdown() {
	defaultMu.Lock()
	registry := defaultRegistry
	defaultMu.Unlock()

	if registry != nil {
		registry.Close()
	}
}
