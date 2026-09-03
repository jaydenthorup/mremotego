package secrets

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	// bwReadyTimeout bounds how long we wait for "bw serve" to accept requests.
	// The CLI is a Node application and can take a few seconds to boot.
	bwReadyTimeout = 20 * time.Second
	// bwReadyPollInterval is how often readiness is probed.
	bwReadyPollInterval = 150 * time.Millisecond
	// bwStopGrace is how long Stop waits for the process to disappear.
	bwStopGrace = 3 * time.Second
	// bwStderrLimit caps the amount of stderr kept for diagnostics.
	bwStderrLimit = 4096
)

// bwServer is a running "bw serve" child process.
type bwServer struct {
	cmd     *exec.Cmd
	port    int
	baseURL string
	stderr  *boundedBuffer
	exited  chan struct{}

	stopOnce sync.Once
}

// boundedBuffer keeps only the last n bytes written to it. The CLI can be
// chatty and only the tail is needed to explain a failed start.
type boundedBuffer struct {
	mu    sync.Mutex
	limit int
	data  []byte
}

func newBoundedBuffer(limit int) *boundedBuffer {
	return &boundedBuffer{limit: limit}
}

func (b *boundedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.data = append(b.data, p...)
	if len(b.data) > b.limit {
		b.data = b.data[len(b.data)-b.limit:]
	}
	return len(p), nil
}

func (b *boundedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return strings.TrimSpace(string(b.data))
}

// findFreePort asks the operating system for an unused loopback port.
func findFreePort() (int, error) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer listener.Close()

	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("unexpected listener address type %T", listener.Addr())
	}
	return addr.Port, nil
}

// startBwServe launches "bw serve" bound to loopback on a free port.
//
// The child inherits the environment, which is how BW_SESSION reaches it: the
// user unlocks the vault once in their shell and every request made through
// this server is then authorised.
func startBwServe(bwPath string) (*bwServer, error) {
	port, err := findFreePort()
	if err != nil {
		return nil, fmt.Errorf("failed to reserve a local port: %w", err)
	}

	args := []string{"serve", "--hostname", "127.0.0.1", "--port", strconv.Itoa(port)}

	var cmd *exec.Cmd
	if isBatchFile(bwPath) {
		// The npm distribution installs bw as a .cmd shim, which cannot be
		// executed directly.
		cmd = exec.Command("cmd", append([]string{"/c", bwPath}, args...)...)
	} else {
		cmd = exec.Command(bwPath, args...)
	}

	cmd.Env = os.Environ()
	cmd.Stdout = io.Discard
	stderr := newBoundedBuffer(bwStderrLimit)
	cmd.Stderr = stderr
	hideConsoleWindow(cmd)
	configureChildProcess(cmd)

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start bw serve: %w", err)
	}

	server := &bwServer{
		cmd:     cmd,
		port:    port,
		baseURL: fmt.Sprintf("http://127.0.0.1:%d", port),
		stderr:  stderr,
		exited:  make(chan struct{}),
	}

	// Reap the child and let readiness polling notice an early exit.
	go func() {
		_ = cmd.Wait()
		close(server.exited)
	}()

	adoptChildProcess(cmd)

	return server, nil
}

func isBatchFile(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasSuffix(lower, ".cmd") || strings.HasSuffix(lower, ".bat")
}

// Stop terminates the server. It is safe to call more than once.
func (s *bwServer) Stop() {
	if s == nil {
		return
	}

	s.stopOnce.Do(func() {
		if s.cmd.Process == nil {
			return
		}

		stopProcess(s.cmd)

		select {
		case <-s.exited:
		case <-time.After(bwStopGrace):
			_ = s.cmd.Process.Kill()
		}
	})
}

// Exited reports whether the server process has already terminated.
func (s *bwServer) Exited() bool {
	select {
	case <-s.exited:
		return true
	default:
		return false
	}
}

// StderrTail returns the tail of the stderr produced by the server, so that a
// failed start can be explained to the user.
func (s *bwServer) StderrTail() string {
	if s == nil || s.stderr == nil {
		return ""
	}
	return s.stderr.String()
}

// waitForReady polls the status endpoint until the server answers, the process
// exits or the deadline passes. A locked or unauthenticated vault still counts
// as ready: the server is up, it simply has nothing to hand out yet.
func waitForReady(ctx context.Context, client *bwClient, exited <-chan struct{}) error {
	ctx, cancel := context.WithTimeout(ctx, bwReadyTimeout)
	defer cancel()

	ticker := time.NewTicker(bwReadyPollInterval)
	defer ticker.Stop()

	var lastErr error
	for {
		_, err := client.status(ctx)
		if err == nil || isVaultStateError(err) {
			return nil
		}
		lastErr = err

		select {
		case <-exited:
			return fmt.Errorf("bw serve exited before it became ready")
		case <-ctx.Done():
			if lastErr != nil {
				return fmt.Errorf("bw serve did not become ready: %w", lastErr)
			}
			return fmt.Errorf("bw serve did not become ready within %s", bwReadyTimeout)
		case <-ticker.C:
		}
	}
}

// isVaultStateError reports whether err means the server works but the vault is
// not usable, as opposed to the server not being up yet.
func isVaultStateError(err error) bool {
	return errors.Is(err, ErrVaultLocked) || errors.Is(err, ErrNotAuthenticated)
}
