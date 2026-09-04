package secrets

import (
	"errors"
	"strings"
	"testing"
)

// fakeProvider is a Provider that resolves references from a static map.
type fakeProvider struct {
	name       string
	scheme     string
	enabled    bool
	authed     bool
	values     map[string]string
	resolveErr error
}

func (f *fakeProvider) Name() string   { return f.name }
func (f *fakeProvider) Scheme() string { return f.scheme }
func (f *fakeProvider) IsEnabled() bool {
	return f.enabled
}
func (f *fakeProvider) IsAuthenticated() bool                 { return f.authed }
func (f *fakeProvider) GetAuthenticationInstructions() string { return "sign in to " + f.name }
func (f *fakeProvider) IsReference(value string) bool {
	return strings.HasPrefix(value, f.scheme+"://")
}
func (f *fakeProvider) ResolveSecret(reference string) (string, error) {
	if f.resolveErr != nil {
		return "", f.resolveErr
	}
	value, ok := f.values[reference]
	if !ok {
		return "", errors.New("not found")
	}
	return value, nil
}

func newTestRegistry() (*Registry, *fakeProvider, *fakeProvider) {
	op := &fakeProvider{
		name: "1Password", scheme: SchemeOnePassword, enabled: true, authed: true,
		values: map[string]string{"op://Private/Server/password": "op-secret"},
	}
	bw := &fakeProvider{
		name: "Bitwarden", scheme: SchemeBitwarden, enabled: true, authed: true,
		values: map[string]string{"bw://item-id": "bw-secret"},
	}
	return NewRegistry(op, bw), op, bw
}

func TestRegistryProviderFor(t *testing.T) {
	registry, op, bw := newTestRegistry()

	tests := []struct {
		value string
		want  Provider
	}{
		{"op://Private/Server/password", op},
		{"bw://item-id", bw},
		{"plain-password", nil},
		{"", nil},
		{"enc:AAAA", nil},
	}

	for _, tt := range tests {
		got, ok := registry.ProviderFor(tt.value)
		if tt.want == nil {
			if ok {
				t.Errorf("ProviderFor(%q) = %s, want no provider", tt.value, got.Name())
			}
			continue
		}
		if !ok || got != tt.want {
			t.Errorf("ProviderFor(%q) = %v, want %s", tt.value, got, tt.want.Name())
		}
	}
}

func TestRegistryResolve(t *testing.T) {
	registry, _, _ := newTestRegistry()

	got, err := registry.Resolve("bw://item-id")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if got != "bw-secret" {
		t.Errorf("Resolve = %q, want %q", got, "bw-secret")
	}

	// Plain values pass through untouched.
	got, err = registry.Resolve("literal")
	if err != nil {
		t.Fatalf("Resolve returned error: %v", err)
	}
	if got != "literal" {
		t.Errorf("Resolve = %q, want %q", got, "literal")
	}
}

func TestRegistryResolveErrorNamesProvider(t *testing.T) {
	op := &fakeProvider{name: "1Password", scheme: SchemeOnePassword, resolveErr: errors.New("boom")}
	registry := NewRegistry(op)

	_, err := registry.Resolve("op://a/b/c")
	if err == nil {
		t.Fatal("expected an error")
	}
	if !strings.Contains(err.Error(), "1Password") {
		t.Errorf("error %q does not name the provider", err)
	}
}

func TestRegistryByScheme(t *testing.T) {
	registry, _, bw := newTestRegistry()

	got, ok := registry.ByScheme(SchemeBitwarden)
	if !ok || got != bw {
		t.Errorf("ByScheme(%q) = %v, want Bitwarden provider", SchemeBitwarden, got)
	}

	if _, ok := registry.ByScheme("nope"); ok {
		t.Error("ByScheme returned a provider for an unknown scheme")
	}
}

func TestIsKnownReference(t *testing.T) {
	tests := map[string]bool{
		"op://Private/Server/password": true,
		"bw://8f3c":                    true,
		"enc:AAAA":                     false,
		"hunter2":                      false,
		"":                             false,
		"opsomething":                  false,
	}

	for value, want := range tests {
		if got := IsKnownReference(value); got != want {
			t.Errorf("IsKnownReference(%q) = %v, want %v", value, got, want)
		}
	}
}

func TestKnownSchemesIsACopy(t *testing.T) {
	schemes := KnownSchemes()
	if len(schemes) == 0 {
		t.Fatal("expected at least one known scheme")
	}
	schemes[0] = "mutated"

	if KnownSchemes()[0] == "mutated" {
		t.Error("KnownSchemes returned the underlying slice")
	}
}
