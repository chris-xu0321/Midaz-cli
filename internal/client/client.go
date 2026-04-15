// Package client provides the HTTP client for communicating with the Midaz API.
package client

import (
	"bytes"
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

// Client is the Midaz API HTTP client.
type Client struct {
	APIURL     string
	Token      string // PAT (sk_...) or JWT — attached as Authorization: Bearer
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

// WithToken returns a shallow copy of the client with Token set.
func (c *Client) WithToken(token string) *Client {
	cp := *c
	cp.Token = token
	return &cp
}

// Get makes a GET request to the API.
func (c *Client) Get(ctx context.Context, path string, params url.Values) (*Response, error) {
	return c.do(ctx, http.MethodGet, path, params, nil)
}

// Post makes a POST request with an optional JSON body.
func (c *Client) Post(ctx context.Context, path string, body any) (*Response, error) {
	return c.doJSON(ctx, http.MethodPost, path, body)
}

// Patch makes a PATCH request with an optional JSON body.
func (c *Client) Patch(ctx context.Context, path string, body any) (*Response, error) {
	return c.doJSON(ctx, http.MethodPatch, path, body)
}

// Delete makes a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) (*Response, error) {
	return c.do(ctx, http.MethodDelete, path, nil, nil)
}

func (c *Client) doJSON(ctx context.Context, method, path string, body any) (*Response, error) {
	var buf io.Reader
	if body != nil {
		raw, err := json.Marshal(body)
		if err != nil {
			return nil, output.Errorf(output.ExitInternal, "internal", "failed to encode request body: %s", err)
		}
		buf = bytes.NewReader(raw)
	}
	return c.do(ctx, method, path, nil, buf)
}

func (c *Client) do(ctx context.Context, method, path string, params url.Values, body io.Reader) (*Response, error) {
	u := c.APIURL + path
	if len(params) > 0 {
		u += "?" + params.Encode()
	}

	req, err := http.NewRequestWithContext(ctx, method, u, body)
	if err != nil {
		return nil, output.ErrNetwork("failed to create request: %s", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, classifyConnError(err, c.APIURL)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, output.ErrNetwork("failed to read response: %s", err)
	}

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return &Response{StatusCode: resp.StatusCode, Body: respBody}, nil
	}

	return nil, classifyHTTPError(resp.StatusCode, respBody, path, c.Token != "")
}

// classifyConnError maps connection-level errors to ExitError.
func classifyConnError(err error, apiURL string) *output.ExitError {
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return output.ErrWithHint(output.ExitNetwork, "timeout",
			fmt.Sprintf("Request timed out to %s", apiURL),
			"check your network connection or increase timeout")
	}
	return output.ErrWithHint(output.ExitNetwork, "network",
		fmt.Sprintf("Cannot connect to Midaz API at %s", apiURL),
		"check your API URL with: midaz config get api_url")
}

// classifyHTTPError maps HTTP status codes to ExitError.
// hasToken indicates whether the request carried an Authorization header; it
// affects the hint we surface for 401.
func classifyHTTPError(status int, body []byte, path string, hasToken bool) *output.ExitError {
	msg := extractAPIMessage(body)

	switch {
	case status == 401:
		if msg == "" {
			msg = "Not authenticated"
		}
		hint := "run 'midaz auth login' to authenticate"
		if hasToken {
			hint = "your credentials are invalid or expired — run 'midaz auth login' again"
		}
		return output.ErrAuth(msg, hint)
	case status == 402:
		if msg == "" {
			msg = "Active subscription required"
		}
		return output.ErrSubscription(msg, "")
	case status == 403:
		if msg == "" {
			msg = fmt.Sprintf("Forbidden: %s", path)
		}
		return output.ErrWithHint(output.ExitAPI, "forbidden", msg,
			"you may lack the required role (owner-only endpoint) or haven't redeemed an invitation code")
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
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if json.Unmarshal(body, &parsed) == nil {
		if parsed.Error != "" {
			return parsed.Error
		}
		if parsed.Message != "" {
			return parsed.Message
		}
	}
	return ""
}
