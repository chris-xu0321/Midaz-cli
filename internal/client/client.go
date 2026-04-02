// Package client provides the HTTP client for communicating with the Seer API.
package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/SparkssL/Midaz-cli/internal/output"
)

// Response holds a parsed HTTP response.
type Response struct {
	StatusCode int
	Body       []byte
}

// Client is the Seer API HTTP client.
type Client struct {
	APIURL     string
	HTTPClient *http.Client
}

// New creates a Client with the given base URL and a 30-second timeout.
func New(apiURL string) *Client {
	return &Client{
		APIURL: strings.TrimRight(apiURL, "/"),
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Get makes a GET request to the API. Path should start with "/".
// Query params are appended if non-nil. Returns Response on 2xx,
// or a classified *output.ExitError on failure.
func (c *Client) Get(ctx context.Context, path string, params url.Values) (*Response, error) {
	u := c.APIURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, output.ErrNetwork("failed to create request: %s", err)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, classifyConnError(err, c.APIURL)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, output.ErrNetwork("failed to read response: %s", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return &Response{StatusCode: resp.StatusCode, Body: body}, nil
	}

	return nil, classifyHTTPError(resp.StatusCode, body, path)
}

// classifyConnError maps connection-level errors to ExitError.
func classifyConnError(err error, apiURL string) *output.ExitError {
	// Check for timeout
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return output.ErrWithHint(output.ExitNetwork, "timeout",
			fmt.Sprintf("Request timed out to %s", apiURL),
			"check your network connection or increase timeout")
	}

	// Connection refused, DNS failure, etc.
	return output.ErrWithHint(output.ExitNetwork, "network",
		fmt.Sprintf("Cannot connect to Seer API at %s", apiURL),
		"check your API URL with: seer-q config get api_url")
}

// classifyHTTPError maps HTTP status codes to ExitError.
func classifyHTTPError(status int, body []byte, path string) *output.ExitError {
	// Try to extract error message from API response
	msg := extractAPIMessage(body)

	switch {
	case status == 404:
		if msg == "" {
			msg = fmt.Sprintf("Not found: %s", path)
		}
		return output.ErrWithHint(output.ExitAPI, "not_found", msg, "")
	case status >= 400 && status < 500:
		if msg == "" {
			msg = fmt.Sprintf("API error %d: %s", status, path)
		}
		return output.ErrAPI("api", "%s", msg)
	default: // 5xx
		if msg == "" {
			msg = fmt.Sprintf("API server error %d: %s", status, path)
		}
		return output.ErrAPI("api", "%s", msg)
	}
}

// extractAPIMessage tries to pull an error message from the API JSON response.
func extractAPIMessage(body []byte) string {
	var parsed struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(body, &parsed) == nil && parsed.Error != "" {
		return parsed.Error
	}
	return ""
}
