package client

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chris-xu0321/Midaz-cli/internal/output"
)

func TestGet_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte(`{"status":"ok"}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	resp, err := c.Get(context.Background(), "/api/health", nil)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	if string(resp.Body) != `{"status":"ok"}` {
		t.Errorf("unexpected body: %s", resp.Body)
	}
}

func TestGet_404(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte(`{"error":"Thread not found"}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.Get(context.Background(), "/api/threads/bad-id", nil)
	if err == nil {
		t.Fatal("expected error for 404")
	}

	var exitErr *output.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatal("expected *ExitError")
	}
	if exitErr.Code != output.ExitAPI {
		t.Errorf("expected exit code %d, got %d", output.ExitAPI, exitErr.Code)
	}
	if exitErr.Detail.Code != "not_found" {
		t.Errorf("expected error code 'not_found', got %q", exitErr.Detail.Code)
	}
}

func TestGet_500(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(500)
		w.Write([]byte(`{"error":"internal server error"}`))
	}))
	defer srv.Close()

	c := New(srv.URL)
	_, err := c.Get(context.Background(), "/api/topics", nil)
	if err == nil {
		t.Fatal("expected error for 500")
	}

	var exitErr *output.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatal("expected *ExitError")
	}
	if exitErr.Code != output.ExitAPI {
		t.Errorf("expected exit code %d, got %d", output.ExitAPI, exitErr.Code)
	}
}

func TestGet_ConnectionRefused(t *testing.T) {
	c := New("http://127.0.0.1:1") // port 1 should refuse connections
	_, err := c.Get(context.Background(), "/api/health", nil)
	if err == nil {
		t.Fatal("expected error for connection refused")
	}

	var exitErr *output.ExitError
	if !errors.As(err, &exitErr) {
		t.Fatal("expected *ExitError")
	}
	if exitErr.Code != output.ExitNetwork {
		t.Errorf("expected exit code %d, got %d", output.ExitNetwork, exitErr.Code)
	}
}
