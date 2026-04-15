package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// DefaultRelayOrigin is the expected Origin header on callback POSTs. The
// browser-relay page must be served from this origin.
var DefaultRelayOrigin = "https://www.midaz.xyz"

// LoopbackResult is what the /cli-auth relay POSTs back to the local server.
type LoopbackResult struct {
	APIKey        string `json:"api_key"`
	WorkspaceID   string `json:"workspace_id"`
	WorkspaceSlug string `json:"workspace_slug"`
	UserEmail     string `json:"user_email"`
	UserID        string `json:"user_id"`
	Nonce         string `json:"nonce"`
}

// Listener owns a one-shot HTTP server that awaits the relay page's POST and
// then self-closes.
type Listener struct {
	Port   int
	Nonce  string
	Origin string // Expected Origin header; default DefaultRelayOrigin

	listener net.Listener
	result   chan *LoopbackResult
	errCh    chan error
}

// NewListener binds 127.0.0.1 to a random port and generates a fresh nonce.
func NewListener() (*Listener, error) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("bind loopback: %w", err)
	}
	nonceBytes := make([]byte, 32)
	if _, err := rand.Read(nonceBytes); err != nil {
		ln.Close()
		return nil, fmt.Errorf("generate nonce: %w", err)
	}
	return &Listener{
		Port:     ln.Addr().(*net.TCPAddr).Port,
		Nonce:    hex.EncodeToString(nonceBytes),
		Origin:   DefaultRelayOrigin,
		listener: ln,
		result:   make(chan *LoopbackResult, 1),
		errCh:    make(chan error, 1),
	}, nil
}

// CallbackURL returns the local URL the relay page must POST to.
func (l *Listener) CallbackURL() string {
	return fmt.Sprintf("http://127.0.0.1:%d/callback", l.Port)
}

// Serve starts the HTTP server and blocks the caller until either (a) a valid
// callback POST arrives, (b) the context is canceled, or (c) timeout elapses.
func (l *Listener) Serve(ctx context.Context, timeout time.Duration) (*LoopbackResult, error) {
	mux := http.NewServeMux()
	mux.HandleFunc("/callback", l.handleCallback)
	mux.HandleFunc("/", l.handleRoot)

	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		if err := server.Serve(l.listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			select {
			case l.errCh <- err:
			default:
			}
		}
	}()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	timer := time.NewTimer(timeout)
	defer timer.Stop()

	select {
	case r := <-l.result:
		return r, nil
	case err := <-l.errCh:
		return nil, err
	case <-timer.C:
		return nil, fmt.Errorf("login timed out after %s", timeout)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Close releases the listener. Idempotent.
func (l *Listener) Close() {
	if l.listener != nil {
		l.listener.Close()
	}
}

func (l *Listener) handleRoot(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "midaz loopback is waiting for /callback", http.StatusNotFound)
}

func (l *Listener) handleCallback(w http.ResponseWriter, r *http.Request) {
	// CORS preflight — mirror the relay origin so the browser can POST.
	if r.Method == http.MethodOptions {
		l.writeCORS(w, r)
		w.WriteHeader(http.StatusNoContent)
		return
	}
	l.writeCORS(w, r)

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := l.verifyOrigin(r); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	var res LoopbackResult
	if err := json.Unmarshal(body, &res); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	if subtle.ConstantTimeCompare([]byte(res.Nonce), []byte(l.Nonce)) != 1 {
		http.Error(w, "nonce mismatch", http.StatusForbidden)
		return
	}
	if !strings.HasPrefix(res.APIKey, "sk_") {
		http.Error(w, "invalid api_key shape", http.StatusBadRequest)
		return
	}

	// The /cli-auth page ignores the response body (only checks response.ok),
	// but we still return JSON so any curl-based debugging gets a well-formed
	// payload.
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"ok":true}`)

	select {
	case l.result <- &res:
	default:
	}
}

func (l *Listener) writeCORS(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")
	if origin == l.Origin {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Vary", "Origin")
	}
}

func (l *Listener) verifyOrigin(r *http.Request) error {
	origin := r.Header.Get("Origin")
	if origin == "" {
		// Some user-agents (curl, Postman) skip Origin. Reject to force the
		// intended browser-relay path.
		return errors.New("missing Origin header")
	}
	if origin != l.Origin {
		return fmt.Errorf("unexpected origin %q", origin)
	}
	return nil
}

