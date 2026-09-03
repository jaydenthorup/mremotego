package secrets

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Errors returned by the Bitwarden client that callers may want to react to.
var (
	// ErrVaultLocked is returned when the vault is locked, so no secret can be
	// read until the user unlocks it.
	ErrVaultLocked = errors.New("bitwarden vault is locked")
	// ErrNotAuthenticated is returned when no user is logged in.
	ErrNotAuthenticated = errors.New("not logged in to bitwarden")
	// ErrNotFound is returned when the requested item does not exist.
	ErrNotFound = errors.New("item not found")
)

// Vault states reported by the Bitwarden CLI.
const (
	bwStatusUnlocked        = "unlocked"
	bwStatusLocked          = "locked"
	bwStatusUnauthenticated = "unauthenticated"
)

// Bitwarden item types; only logins carry credentials.
const bwTypeLogin = 1

const (
	bwRequestTimeout = 15 * time.Second
	bwSyncTimeout    = 60 * time.Second
)

// bwClient talks to the local REST API exposed by "bw serve".
type bwClient struct {
	baseURL string
	http    *http.Client
}

func newBwClient(baseURL string) *bwClient {
	return &bwClient{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		http:    &http.Client{Timeout: bwSyncTimeout},
	}
}

// bwEnvelope is the response wrapper used by every "bw serve" endpoint.
type bwEnvelope struct {
	Success bool            `json:"success"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

type bwStatus struct {
	ServerURL string `json:"serverUrl"`
	UserEmail string `json:"userEmail"`
	Status    string `json:"status"`
}

type bwURI struct {
	URI   string `json:"uri"`
	Match *int   `json:"match,omitempty"`
}

type bwLogin struct {
	Username string  `json:"username,omitempty"`
	Password string  `json:"password,omitempty"`
	Totp     string  `json:"totp,omitempty"`
	URIs     []bwURI `json:"uris,omitempty"`
}

type bwItem struct {
	ID       string   `json:"id,omitempty"`
	Type     int      `json:"type"`
	Name     string   `json:"name"`
	Notes    string   `json:"notes,omitempty"`
	Login    *bwLogin `json:"login,omitempty"`
	FolderID string   `json:"folderId,omitempty"`
}

// bwAPIError carries the message returned by the CLI so that it can be shown
// to the user unchanged.
type bwAPIError struct {
	Message string
}

func (e *bwAPIError) Error() string { return e.Message }

// do performs a request and unwraps the response envelope. The Origin header is
// deliberately never set: "bw serve" rejects any request that carries one.
func (c *bwClient) do(ctx context.Context, method, path string, body any) (json.RawMessage, error) {
	var reader *bytes.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to encode request: %w", err)
		}
		reader = bytes.NewReader(encoded)
	} else {
		reader = bytes.NewReader(nil)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var envelope bwEnvelope
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return nil, fmt.Errorf("unexpected response from bw serve (HTTP %d): %w", resp.StatusCode, err)
	}

	if !envelope.Success {
		return nil, classifyBwError(envelope.Message, resp.StatusCode)
	}

	return envelope.Data, nil
}

// classifyBwError maps CLI messages onto the sentinel errors above.
func classifyBwError(message string, statusCode int) error {
	lower := strings.ToLower(message)

	switch {
	case strings.Contains(lower, "vault is locked"):
		return fmt.Errorf("%w", ErrVaultLocked)
	case strings.Contains(lower, "not logged in"), strings.Contains(lower, "you are not logged in"):
		return fmt.Errorf("%w", ErrNotAuthenticated)
	case strings.Contains(lower, "not found"), statusCode == http.StatusNotFound:
		if message == "" {
			return ErrNotFound
		}
		return fmt.Errorf("%w: %s", ErrNotFound, message)
	case message == "":
		return fmt.Errorf("bw serve returned HTTP %d", statusCode)
	}

	return &bwAPIError{Message: message}
}

// status reports whether the vault is unlocked, locked or unauthenticated.
func (c *bwClient) status(ctx context.Context) (*bwStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, bwRequestTimeout)
	defer cancel()

	data, err := c.do(ctx, http.MethodGet, "/status", nil)
	if err != nil {
		return nil, err
	}

	// The payload is wrapped in a template object:
	// {"object":"template","template":{"status":"unlocked",...}}
	var wrapper struct {
		Template bwStatus `json:"template"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse status: %w", err)
	}

	return &wrapper.Template, nil
}

// getField reads a single field of an item, e.g. its password.
func (c *bwClient) getField(ctx context.Context, field, id string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, bwRequestTimeout)
	defer cancel()

	path := "/object/" + url.PathEscape(field) + "/" + url.PathEscape(id)
	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return "", err
	}

	// Scalar fields are returned as {"object":"string","data":"<value>"}.
	var wrapper struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return "", fmt.Errorf("failed to parse %s: %w", field, err)
	}

	return wrapper.Data, nil
}

// listItems returns the vault items, optionally narrowed by a search term.
func (c *bwClient) listItems(ctx context.Context, search string) ([]bwItem, error) {
	ctx, cancel := context.WithTimeout(ctx, bwRequestTimeout)
	defer cancel()

	path := "/list/object/items"
	if search != "" {
		path += "?search=" + url.QueryEscape(search)
	}

	data, err := c.do(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	var wrapper struct {
		Data []bwItem `json:"data"`
	}
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("failed to parse item list: %w", err)
	}

	return wrapper.Data, nil
}

// createItem stores a new item and returns it including the assigned id.
func (c *bwClient) createItem(ctx context.Context, item bwItem) (*bwItem, error) {
	ctx, cancel := context.WithTimeout(ctx, bwRequestTimeout)
	defer cancel()

	data, err := c.do(ctx, http.MethodPost, "/object/item", item)
	if err != nil {
		return nil, err
	}

	var created bwItem
	if err := json.Unmarshal(data, &created); err != nil {
		return nil, fmt.Errorf("failed to parse created item: %w", err)
	}

	if created.ID == "" {
		return nil, fmt.Errorf("bitwarden did not return an item id")
	}

	return &created, nil
}

// sync pulls the latest vault state from the server.
func (c *bwClient) sync(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, bwSyncTimeout)
	defer cancel()

	_, err := c.do(ctx, http.MethodPost, "/sync", nil)
	return err
}
