package secrets

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// fakeBw is a stand-in for "bw serve" so the tests never touch the real CLI.
type fakeBw struct {
	server *httptest.Server

	state       string
	items       map[string]bwItem
	lastCreated bwItem
	sawOrigin   bool
	lastSearch  string
	syncCalls   int
}

func newFakeBw(t *testing.T) *fakeBw {
	t.Helper()

	f := &fakeBw{
		state: bwStatusUnlocked,
		items: make(map[string]bwItem),
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		f.note(r)
		writeEnvelope(w, map[string]any{
			"object": "template",
			"template": map[string]string{
				"serverUrl": "https://vault.example.com",
				"userEmail": "user@example.com",
				"status":    f.state,
			},
		})
	})

	mux.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
		f.note(r)
		if f.locked(w) {
			return
		}
		f.syncCalls++
		writeEnvelope(w, map[string]string{"object": "message", "title": "Syncing complete."})
	})

	mux.HandleFunc("/list/object/items", func(w http.ResponseWriter, r *http.Request) {
		f.note(r)
		if f.locked(w) {
			return
		}
		f.lastSearch = r.URL.Query().Get("search")

		list := make([]bwItem, 0, len(f.items))
		for _, item := range f.items {
			list = append(list, item)
		}
		writeEnvelope(w, map[string]any{"object": "list", "data": list})
	})

	mux.HandleFunc("/object/", func(w http.ResponseWriter, r *http.Request) {
		f.note(r)
		if f.locked(w) {
			return
		}

		if r.Method == http.MethodPost {
			var item bwItem
			if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
				writeError(w, http.StatusBadRequest, "bad request body")
				return
			}
			item.ID = "created-id"
			f.lastCreated = item
			f.items[item.ID] = item
			writeEnvelope(w, item)
			return
		}

		// /object/<kind>/<id>
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) != 3 {
			writeError(w, http.StatusNotFound, "Not found.")
			return
		}
		kind, id := parts[1], parts[2]

		item, ok := f.items[id]
		if !ok {
			writeError(w, http.StatusNotFound, "Not found.")
			return
		}

		if kind == "item" {
			writeEnvelope(w, item)
			return
		}

		var value string
		if item.Login != nil {
			switch kind {
			case bitwardenFieldPassword:
				value = item.Login.Password
			case bitwardenFieldUsername:
				value = item.Login.Username
			case bitwardenFieldTotp:
				value = item.Login.Totp
			}
		}
		if kind == bitwardenFieldNotes {
			value = item.Notes
		}

		writeEnvelope(w, map[string]string{"object": "string", "data": value})
	})

	f.server = httptest.NewServer(mux)
	t.Cleanup(f.server.Close)
	return f
}

func (f *fakeBw) note(r *http.Request) {
	if r.Header.Get("Origin") != "" {
		f.sawOrigin = true
	}
}

// locked mimics the CLI, which rejects vault access while locked.
func (f *fakeBw) locked(w http.ResponseWriter) bool {
	switch f.state {
	case bwStatusLocked:
		writeError(w, http.StatusOK, "Vault is locked.")
		return true
	case bwStatusUnauthenticated:
		writeError(w, http.StatusOK, "You are not logged in.")
		return true
	}
	return false
}

func (f *fakeBw) provider() *BitwardenProvider {
	return newBitwardenProviderWithClient(newBwClient(f.server.URL))
}

func (f *fakeBw) addLogin(id, name, username, password string) {
	f.items[id] = bwItem{
		ID:    id,
		Type:  bwTypeLogin,
		Name:  name,
		Login: &bwLogin{Username: username, Password: password},
	}
}

func writeEnvelope(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{"success": true, "data": data})
}

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"success": false, "message": message})
}

func TestBitwardenResolvePassword(t *testing.T) {
	fake := newFakeBw(t)
	fake.addLogin("abc-123", "Web Server", "admin", "hunter2")

	got, err := fake.provider().ResolveSecret("bw://abc-123")
	if err != nil {
		t.Fatalf("ResolveSecret returned error: %v", err)
	}
	if got != "hunter2" {
		t.Errorf("ResolveSecret = %q, want %q", got, "hunter2")
	}
}

func TestBitwardenResolveExplicitField(t *testing.T) {
	fake := newFakeBw(t)
	fake.addLogin("abc-123", "Web Server", "admin", "hunter2")

	got, err := fake.provider().ResolveSecret("bw://abc-123/username")
	if err != nil {
		t.Fatalf("ResolveSecret returned error: %v", err)
	}
	if got != "admin" {
		t.Errorf("ResolveSecret = %q, want %q", got, "admin")
	}
}

func TestBitwardenResolveLockedVault(t *testing.T) {
	fake := newFakeBw(t)
	fake.addLogin("abc-123", "Web Server", "admin", "hunter2")
	fake.state = bwStatusLocked

	_, err := fake.provider().ResolveSecret("bw://abc-123")
	if !errors.Is(err, ErrVaultLocked) {
		t.Fatalf("ResolveSecret error = %v, want ErrVaultLocked", err)
	}
}

func TestBitwardenResolveUnauthenticated(t *testing.T) {
	fake := newFakeBw(t)
	fake.state = bwStatusUnauthenticated

	_, err := fake.provider().ResolveSecret("bw://abc-123")
	if !errors.Is(err, ErrNotAuthenticated) {
		t.Fatalf("ResolveSecret error = %v, want ErrNotAuthenticated", err)
	}
}

func TestBitwardenResolveUnknownItem(t *testing.T) {
	fake := newFakeBw(t)

	_, err := fake.provider().ResolveSecret("bw://missing")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("ResolveSecret error = %v, want ErrNotFound", err)
	}
}

func TestBitwardenIsAuthenticated(t *testing.T) {
	tests := map[string]bool{
		bwStatusUnlocked:        true,
		bwStatusLocked:          false,
		bwStatusUnauthenticated: false,
	}

	for state, want := range tests {
		fake := newFakeBw(t)
		fake.state = state

		if got := fake.provider().IsAuthenticated(); got != want {
			t.Errorf("IsAuthenticated with status %q = %v, want %v", state, got, want)
		}
	}
}

func TestBitwardenCreateItem(t *testing.T) {
	fake := newFakeBw(t)

	reference, err := fake.provider().CreateItem(CreateItemRequest{
		Title:    "Web Server",
		Username: "admin",
		Password: "hunter2",
		URI:      "ssh://web.example.com",
	})
	if err != nil {
		t.Fatalf("CreateItem returned error: %v", err)
	}
	if reference != "bw://created-id" {
		t.Errorf("CreateItem = %q, want %q", reference, "bw://created-id")
	}

	created := fake.lastCreated
	if created.Type != bwTypeLogin {
		t.Errorf("created item type = %d, want %d", created.Type, bwTypeLogin)
	}
	if created.Name != "Web Server" {
		t.Errorf("created item name = %q, want %q", created.Name, "Web Server")
	}
	if created.Login == nil || created.Login.Password != "hunter2" || created.Login.Username != "admin" {
		t.Fatalf("created item login = %+v, want admin/hunter2", created.Login)
	}
	if len(created.Login.URIs) != 1 || created.Login.URIs[0].URI != "ssh://web.example.com" {
		t.Errorf("created item URIs = %+v, want one ssh:// entry", created.Login.URIs)
	}
}

func TestBitwardenCreateItemRequiresTitle(t *testing.T) {
	fake := newFakeBw(t)

	if _, err := fake.provider().CreateItem(CreateItemRequest{Password: "hunter2"}); err == nil {
		t.Error("expected CreateItem to reject an empty title")
	}
}

func TestBitwardenListLoginItemsSkipsOtherTypes(t *testing.T) {
	fake := newFakeBw(t)
	fake.addLogin("login-1", "Web Server", "admin", "hunter2")
	fake.items["note-1"] = bwItem{ID: "note-1", Type: 2, Name: "Secure Note"}

	items, err := fake.provider().ListLoginItems("")
	if err != nil {
		t.Fatalf("ListLoginItems returned error: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("ListLoginItems returned %d items, want 1", len(items))
	}
	if items[0].ID != "login-1" || items[0].Username != "admin" {
		t.Errorf("ListLoginItems = %+v, want the login item", items[0])
	}
}

func TestBitwardenListLoginItemsPassesSearch(t *testing.T) {
	fake := newFakeBw(t)

	if _, err := fake.provider().ListLoginItems("web server"); err != nil {
		t.Fatalf("ListLoginItems returned error: %v", err)
	}
	if fake.lastSearch != "web server" {
		t.Errorf("search term = %q, want %q", fake.lastSearch, "web server")
	}
}

func TestBitwardenSync(t *testing.T) {
	fake := newFakeBw(t)

	if err := fake.provider().Sync(); err != nil {
		t.Fatalf("Sync returned error: %v", err)
	}
	if fake.syncCalls != 1 {
		t.Errorf("sync called %d times, want 1", fake.syncCalls)
	}
}

// bw serve rejects any request carrying an Origin header, so the client must
// never send one.
func TestBitwardenClientSendsNoOriginHeader(t *testing.T) {
	fake := newFakeBw(t)
	fake.addLogin("abc-123", "Web Server", "admin", "hunter2")

	provider := fake.provider()
	if _, err := provider.ResolveSecret("bw://abc-123"); err != nil {
		t.Fatalf("ResolveSecret returned error: %v", err)
	}
	if _, err := provider.ListLoginItems(""); err != nil {
		t.Fatalf("ListLoginItems returned error: %v", err)
	}

	if fake.sawOrigin {
		t.Error("client sent an Origin header, which bw serve would reject")
	}
}

func TestBitwardenDisabledProviderReportsMissingCLI(t *testing.T) {
	provider := &BitwardenProvider{enabled: false}

	_, err := provider.ResolveSecret("bw://abc-123")
	if err == nil {
		t.Fatal("expected an error when the CLI is missing")
	}
	if !strings.Contains(err.Error(), "not installed") {
		t.Errorf("error = %q, want it to mention the missing CLI", err)
	}

	if provider.IsAuthenticated() {
		t.Error("a disabled provider must not report itself as authenticated")
	}
}

func TestBitwardenIsReference(t *testing.T) {
	provider := &BitwardenProvider{}

	tests := map[string]bool{
		"bw://abc":                     true,
		"bw://abc/username":            true,
		"op://Private/Server/password": false,
		"hunter2":                      false,
		"":                             false,
	}

	for value, want := range tests {
		if got := provider.IsReference(value); got != want {
			t.Errorf("IsReference(%q) = %v, want %v", value, got, want)
		}
	}
}

func TestBitwardenResolveRejectsForeignReference(t *testing.T) {
	fake := newFakeBw(t)

	if _, err := fake.provider().ResolveSecret("op://Private/Server/password"); err == nil {
		t.Error("expected a 1Password reference to be rejected")
	}
}

func TestBitwardenStatusReported(t *testing.T) {
	fake := newFakeBw(t)
	fake.state = bwStatusLocked

	status, err := fake.provider().Status()
	if err != nil {
		t.Fatalf("Status returned error: %v", err)
	}
	if status != bwStatusLocked {
		t.Errorf("Status = %q, want %q", status, bwStatusLocked)
	}
}

func TestWaitForReadyAcceptsLockedVault(t *testing.T) {
	fake := newFakeBw(t)
	fake.state = bwStatusLocked

	// A locked vault still means the server is up and usable once unlocked.
	if err := waitForReady(context.Background(), newBwClient(fake.server.URL), make(chan struct{})); err != nil {
		t.Errorf("waitForReady returned error for a locked vault: %v", err)
	}
}

func TestWaitForReadyDetectsExit(t *testing.T) {
	// A server that never answers, plus an already exited process.
	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer dead.Close()

	exited := make(chan struct{})
	close(exited)

	err := waitForReady(context.Background(), newBwClient(dead.URL), exited)
	if err == nil {
		t.Fatal("expected an error when the process has exited")
	}
	if !strings.Contains(err.Error(), "exited") {
		t.Errorf("error = %q, want it to mention the exit", err)
	}
}

func TestWaitForReadyTimesOut(t *testing.T) {
	dead := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer dead.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 400*time.Millisecond)
	defer cancel()

	if err := waitForReady(ctx, newBwClient(dead.URL), make(chan struct{})); err == nil {
		t.Fatal("expected waitForReady to give up")
	}
}

func TestFindFreePort(t *testing.T) {
	port, err := findFreePort()
	if err != nil {
		t.Fatalf("findFreePort returned error: %v", err)
	}
	if port <= 0 || port > 65535 {
		t.Errorf("findFreePort = %d, want a valid TCP port", port)
	}
}

func TestBoundedBufferKeepsTail(t *testing.T) {
	buf := newBoundedBuffer(8)
	if _, err := buf.Write([]byte("abcdefghijkl")); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	if got := buf.String(); got != "efghijkl" {
		t.Errorf("String = %q, want %q", got, "efghijkl")
	}
}

func TestIsBatchFile(t *testing.T) {
	tests := map[string]bool{
		`C:\npm\bw.cmd`:   true,
		`C:\npm\bw.CMD`:   true,
		`C:\tools\bw.bat`: true,
		`C:\tools\bw.exe`: false,
		"/usr/bin/bw":     false,
	}

	for path, want := range tests {
		if got := isBatchFile(path); got != want {
			t.Errorf("isBatchFile(%q) = %v, want %v", path, got, want)
		}
	}
}
