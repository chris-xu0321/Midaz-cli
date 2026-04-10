package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"
)

// CallbackResult holds the data received from the browser callback.
type CallbackResult struct {
	APIKey        string `json:"api_key"`
	Prefix        string `json:"prefix"`
	WorkspaceID   string `json:"workspace_id"`
	WorkspaceSlug string `json:"workspace_slug"`
	UserEmail     string `json:"user_email"`
}

// StartCallbackServer starts a localhost HTTP server on a random port.
// It returns the port, a channel that receives the callback result,
// and a cleanup function to shut down the server.
func StartCallbackServer(ctx context.Context) (int, <-chan CallbackResult, func()) {
	resultCh := make(chan CallbackResult, 1)

	mux := http.NewServeMux()
	mux.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers for browser POST from auth callback page
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		var result CallbackResult
		if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		fmt.Fprint(w, `<html><body style="background:#0F1523;color:#D6D8D8;font-family:system-ui;display:flex;align-items:center;justify-content:center;min-height:100vh"><p>Authenticated. You can close this tab.</p></body></html>`)

		select {
		case resultCh <- result:
		default:
		}
	})

	// Listen on random port
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		// fallback: return 0 port, closed channel
		close(resultCh)
		return 0, resultCh, func() {}
	}

	port := listener.Addr().(*net.TCPAddr).Port
	server := &http.Server{Handler: mux}

	go func() {
		_ = server.Serve(listener)
	}()

	cleanup := func() {
		shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(shutCtx)
	}

	return port, resultCh, cleanup
}
