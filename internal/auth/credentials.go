// Package auth manages CLI authentication credentials.
//
// Credentials are stored in ~/.config/seer/credentials.json (mode 0600).
// The file contains the Seer PAT (sk_...) and workspace metadata.
package auth

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// Credentials holds the stored authentication state.
type Credentials struct {
	APIKey        string `json:"api_key"`
	WorkspaceID   string `json:"workspace_id,omitempty"`
	WorkspaceSlug string `json:"workspace_slug,omitempty"`
	UserEmail     string `json:"user_email,omitempty"`
}

// CredentialsPath returns the path to the credentials file.
func CredentialsPath() string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(".", "seer", "credentials.json")
	}
	return filepath.Join(dir, "seer", "credentials.json")
}

// Load reads credentials from disk. Returns nil if file does not exist.
func Load() (*Credentials, error) {
	data, err := os.ReadFile(CredentialsPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, err
	}
	var creds Credentials
	if err := json.Unmarshal(data, &creds); err != nil {
		return nil, err
	}
	return &creds, nil
}

// Save writes credentials to disk with restricted permissions.
func Save(creds *Credentials) error {
	path := CredentialsPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0600)
}

// Clear deletes the credentials file.
func Clear() error {
	err := os.Remove(CredentialsPath())
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

// Exists returns true if a credentials file exists with a non-empty API key.
func Exists() bool {
	creds, err := Load()
	return err == nil && creds != nil && creds.APIKey != ""
}

// ResolveToken returns the auth token from (in priority order):
//  1. SEER_API_KEY env var
//  2. config api_key (passed in)
//  3. credentials.json
//  4. "" (no auth)
func ResolveToken(configAPIKey string) string {
	if v := os.Getenv("SEER_API_KEY"); v != "" {
		return v
	}
	if configAPIKey != "" {
		return configAPIKey
	}
	creds, err := Load()
	if err == nil && creds != nil && creds.APIKey != "" {
		return creds.APIKey
	}
	return ""
}
