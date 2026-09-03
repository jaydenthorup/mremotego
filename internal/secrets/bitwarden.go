package secrets

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

// BitwardenProvider resolves "bw://" references through the Bitwarden CLI.
//
// The CLI has no library interface, so the provider runs "bw serve" as a child
// process and talks to its REST API over loopback. The server is started lazily
// on first use, bound to 127.0.0.1 on a random free port, and terminated when
// Close is called.
//
// The child inherits BW_SESSION from the environment. Unlocking therefore
// happens once, in the shell the application was started from, and no master
// password is ever handled by MremoteGO itself.
type BitwardenProvider struct {
	bwPath  string
	enabled bool

	mu     sync.Mutex
	server *bwServer
	client *bwClient

	// newClient creates the transport. Tests replace it to talk to an httptest
	// server instead of spawning a real CLI.
	newClient func() (*bwClient, *bwServer, error)
}

var (
	_ Provider    = (*BitwardenProvider)(nil)
	_ ItemCreator = (*BitwardenProvider)(nil)
)

// BitwardenItem is a login item as shown in the item picker.
type BitwardenItem struct {
	ID       string
	Name     string
	Username string
	URI      string
}

// NewBitwardenProvider creates a provider backed by the Bitwarden CLI. It only
// looks the binary up; no process is started until a secret is needed.
func NewBitwardenProvider() *BitwardenProvider {
	path, err := exec.LookPath("bw")
	provider := &BitwardenProvider{
		bwPath:  path,
		enabled: err == nil,
	}
	provider.newClient = provider.startServer
	return provider
}

// newBitwardenProviderWithClient builds a provider around an existing client.
// It exists for tests, which must not spawn the real CLI.
func newBitwardenProviderWithClient(client *bwClient) *BitwardenProvider {
	provider := &BitwardenProvider{
		bwPath:  "bw",
		enabled: true,
	}
	provider.newClient = func() (*bwClient, *bwServer, error) {
		return client, nil, nil
	}
	return provider
}

// Name returns the human readable provider name.
func (p *BitwardenProvider) Name() string { return "Bitwarden" }

// Scheme returns the reference scheme handled by this provider.
func (p *BitwardenProvider) Scheme() string { return SchemeBitwarden }

// IsEnabled reports whether the Bitwarden CLI is installed.
func (p *BitwardenProvider) IsEnabled() bool { return p.enabled }

// IsReference checks whether a value is a Bitwarden reference.
func (p *BitwardenProvider) IsReference(value string) bool {
	return strings.HasPrefix(value, bitwardenPrefix)
}

// Status returns the vault state: unlocked, locked or unauthenticated.
func (p *BitwardenProvider) Status() (string, error) {
	client, err := p.ensureClient()
	if err != nil {
		return "", err
	}

	status, err := client.status(context.Background())
	if err != nil {
		return "", err
	}
	return status.Status, nil
}

// IsAuthenticated reports whether the vault is unlocked and secrets can be
// read. Note that this may start the helper process, so it should not be called
// on a UI thread.
func (p *BitwardenProvider) IsAuthenticated() bool {
	status, err := p.Status()
	return err == nil && status == bwStatusUnlocked
}

// GetAuthenticationInstructions explains how to make the vault available.
func (p *BitwardenProvider) GetAuthenticationInstructions() string {
	return `Bitwarden CLI is installed but the vault is not available.

MremoteGO reads secrets through the Bitwarden CLI, which needs an unlocked
vault. Unlock it once in a terminal and start MremoteGO from that same
terminal, so it inherits the session key.

PowerShell:
  bw config server https://your-server   (self-hosted or Vaultwarden only)
  bw login                               (once)
  $env:BW_SESSION = bw unlock --raw
  .\mremotego.exe

bash / zsh:
  bw config server https://your-server   (self-hosted or Vaultwarden only)
  bw login                               (once)
  export BW_SESSION="$(bw unlock --raw)"
  ./mremotego

Run "bw status" to confirm the vault reports "unlocked".

MremoteGO never asks for or stores your master password. It starts
"bw serve" on 127.0.0.1 with a random port and stops it on exit.`
}

// ResolveSecret reads the field a reference points at.
func (p *BitwardenProvider) ResolveSecret(reference string) (string, error) {
	id, field, err := parseBitwardenReference(reference)
	if err != nil {
		return "", err
	}

	client, err := p.ensureClient()
	if err != nil {
		return "", err
	}

	value, err := client.getField(context.Background(), field, id)
	if err != nil {
		return "", err
	}

	if value == "" {
		return "", fmt.Errorf("item %s has no %s", id, field)
	}

	return value, nil
}

// CreateItem stores a new login item and returns its "bw://" reference.
func (p *BitwardenProvider) CreateItem(req CreateItemRequest) (string, error) {
	if req.Title == "" {
		return "", fmt.Errorf("a title is required")
	}

	client, err := p.ensureClient()
	if err != nil {
		return "", err
	}

	login := &bwLogin{
		Username: req.Username,
		Password: req.Password,
	}
	if req.URI != "" {
		login.URIs = []bwURI{{URI: req.URI}}
	}

	created, err := client.createItem(context.Background(), bwItem{
		Type:  bwTypeLogin,
		Name:  req.Title,
		Login: login,
	})
	if err != nil {
		return "", err
	}

	return bitwardenReference(created.ID), nil
}

// ListLoginItems returns the login items of the vault. Cards, identities and
// secure notes are skipped because they carry no connection credentials.
func (p *BitwardenProvider) ListLoginItems(search string) ([]BitwardenItem, error) {
	client, err := p.ensureClient()
	if err != nil {
		return nil, err
	}

	items, err := client.listItems(context.Background(), search)
	if err != nil {
		return nil, err
	}

	var logins []BitwardenItem
	for _, item := range items {
		if item.Type != bwTypeLogin {
			continue
		}

		entry := BitwardenItem{ID: item.ID, Name: item.Name}
		if item.Login != nil {
			entry.Username = item.Login.Username
			if len(item.Login.URIs) > 0 {
				entry.URI = item.Login.URIs[0].URI
			}
		}
		logins = append(logins, entry)
	}

	return logins, nil
}

// Sync pulls the latest vault contents from the server. The CLI serves items
// from a local cache, so items created elsewhere only appear after a sync.
func (p *BitwardenProvider) Sync() error {
	client, err := p.ensureClient()
	if err != nil {
		return err
	}
	return client.sync(context.Background())
}

// Warmup starts the helper process in the background so that a later call from
// the UI thread does not have to wait for the CLI to boot.
func (p *BitwardenProvider) Warmup() {
	if !p.enabled {
		return
	}
	go func() {
		_, _ = p.ensureClient()
	}()
}

// Close stops the helper process.
func (p *BitwardenProvider) Close() error {
	p.mu.Lock()
	server, client := p.server, p.client
	p.server, p.client = nil, nil
	p.mu.Unlock()

	if server != nil {
		server.Stop()
	}
	_ = client

	return nil
}

// ensureClient returns a client for a running server, starting one if needed.
func (p *BitwardenProvider) ensureClient() (*bwClient, error) {
	if !p.enabled {
		return nil, fmt.Errorf("Bitwarden CLI (bw) is not installed or not in PATH")
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	// Restart if a previously started server has died.
	if p.client != nil && (p.server == nil || !p.server.Exited()) {
		return p.client, nil
	}
	if p.server != nil {
		p.server.Stop()
		p.server, p.client = nil, nil
	}

	client, server, err := p.newClient()
	if err != nil {
		return nil, err
	}

	p.client, p.server = client, server
	return client, nil
}

// startServer spawns "bw serve" and waits until it answers requests.
func (p *BitwardenProvider) startServer() (*bwClient, *bwServer, error) {
	server, err := startBwServe(p.bwPath)
	if err != nil {
		return nil, nil, err
	}

	client := newBwClient(server.baseURL)
	if err := waitForReady(context.Background(), client, server.exited); err != nil {
		server.Stop()
		if tail := server.StderrTail(); tail != "" {
			return nil, nil, fmt.Errorf("%w: %s", err, tail)
		}
		return nil, nil, err
	}

	return client, server, nil
}
