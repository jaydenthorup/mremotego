package secrets

import "strings"

// Scheme identifiers for the supported secret providers. A connection password
// that starts with "<scheme>://" is a reference and is resolved at launch time
// instead of being stored in the configuration file.
const (
	SchemeOnePassword = "op"
	SchemeBitwarden   = "bw"
)

// knownSchemes lists every scheme understood by this package. It is kept as a
// package level variable (rather than derived from a Registry) so that callers
// which must not instantiate providers - such as the crypto package deciding
// whether a value needs encrypting - can still recognise a reference.
var knownSchemes = []string{SchemeOnePassword, SchemeBitwarden}

// Provider resolves secret references belonging to a single scheme.
type Provider interface {
	// Name returns a human readable provider name, e.g. "1Password".
	Name() string

	// Scheme returns the URL scheme handled by this provider, e.g. "op".
	Scheme() string

	// IsEnabled reports whether the backing tool is available on this machine.
	IsEnabled() bool

	// IsAuthenticated reports whether secrets can currently be retrieved.
	IsAuthenticated() bool

	// GetAuthenticationInstructions returns user facing text explaining how to
	// authenticate when IsAuthenticated returns false.
	GetAuthenticationInstructions() string

	// IsReference reports whether value is a reference for this provider.
	IsReference(value string) bool

	// ResolveSecret retrieves the secret a reference points at.
	ResolveSecret(reference string) (string, error)
}

// CreateItemRequest is the provider agnostic input for creating a login item.
// Fields that a provider does not support are ignored by that provider.
type CreateItemRequest struct {
	Vault    string // 1Password vault name; ignored by Bitwarden
	Title    string
	Username string
	Password string
	URI      string // e.g. "ssh://host"; ignored by 1Password
}

// ItemCreator is implemented by providers able to store new login items.
type ItemCreator interface {
	// CreateItem stores a login item and returns a reference to it.
	CreateItem(req CreateItemRequest) (reference string, err error)
}

// KnownSchemes returns the schemes understood by this package.
func KnownSchemes() []string {
	out := make([]string, len(knownSchemes))
	copy(out, knownSchemes)
	return out
}

// IsKnownReference reports whether value uses any known scheme prefix. It does
// not require a provider to be installed or authenticated.
func IsKnownReference(value string) bool {
	for _, scheme := range knownSchemes {
		if strings.HasPrefix(value, scheme+"://") {
			return true
		}
	}
	return false
}
