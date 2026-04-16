// Package auth manages PAT (personal access token) credentials for the midaz CLI.
//
// Credentials live in a single JSON file (auth.json) at
// <userConfigDir>/midaz/auth.json with 0600 permissions. The file carries a
// profile-shaped envelope so future `--profile` support is a pure additive.
//
// MIDAZ_TOKEN env var is honored as an override: when set, it is treated as the
// current profile's api_key and bypasses the file entirely. Useful for CI.
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Credentials captures a single profile's auth state.
type Credentials struct {
	APIKey     string `json:"api_key"`
	DeskID     string `json:"desk_id,omitempty"`
	DeskSlug   string `json:"desk_slug,omitempty"`
	UserEmail  string `json:"user_email,omitempty"`
	UserID     string `json:"user_id,omitempty"`
	VerifiedAt string `json:"verified_at,omitempty"`
	Label      string `json:"label,omitempty"`
}

// UnmarshalJSON accepts both the canonical desk_* fields and the legacy
// workspace_* fields written by older CLI versions. Re-saves as canonical
// on the next SetCurrent.
func (c *Credentials) UnmarshalJSON(b []byte) error {
	var raw struct {
		APIKey        string `json:"api_key"`
		DeskID        string `json:"desk_id"`
		DeskSlug      string `json:"desk_slug"`
		WorkspaceID   string `json:"workspace_id"`
		WorkspaceSlug string `json:"workspace_slug"`
		UserEmail     string `json:"user_email"`
		UserID        string `json:"user_id"`
		VerifiedAt    string `json:"verified_at"`
		Label         string `json:"label"`
	}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	c.APIKey = raw.APIKey
	c.DeskID = raw.DeskID
	if c.DeskID == "" {
		c.DeskID = raw.WorkspaceID
	}
	c.DeskSlug = raw.DeskSlug
	if c.DeskSlug == "" {
		c.DeskSlug = raw.WorkspaceSlug
	}
	c.UserEmail = raw.UserEmail
	c.UserID = raw.UserID
	c.VerifiedAt = raw.VerifiedAt
	c.Label = raw.Label
	return nil
}

// Store is the on-disk auth file shape.
type Store struct {
	Current  string                  `json:"current"`
	Profiles map[string]*Credentials `json:"profiles"`
}

// DefaultProfile is the name used when --profile is not set.
const DefaultProfile = "default"

// AuthPath returns the canonical auth file path.
func AuthPath() string {
	if p := os.Getenv("MIDAZ_AUTH_PATH"); p != "" {
		return p
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(".", "midaz", "auth.json")
	}
	return filepath.Join(dir, "midaz", "auth.json")
}

// Load reads the auth store from disk. Returns an empty store when the file
// doesn't exist.
func Load() (*Store, error) {
	data, err := os.ReadFile(AuthPath())
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Store{Current: DefaultProfile, Profiles: map[string]*Credentials{}}, nil
		}
		return nil, fmt.Errorf("read auth file: %w", err)
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("parse auth file: %w", err)
	}
	if s.Profiles == nil {
		s.Profiles = map[string]*Credentials{}
	}
	if s.Current == "" {
		s.Current = DefaultProfile
	}
	return &s, nil
}

// Save writes the auth store to disk with 0600 permissions. Directory is
// created with 0700 to keep cross-user access tight.
func Save(s *Store) error {
	path := AuthPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("create auth dir: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("encode auth file: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return fmt.Errorf("write auth file: %w", err)
	}
	return nil
}

// Current returns the credentials for the currently-active profile, honoring
// MIDAZ_TOKEN overrides. Returns nil, nil when the user is not logged in
// (callers should surface ErrAuth).
func Current(profile string) (*Creds, error) {
	if token := os.Getenv("MIDAZ_TOKEN"); token != "" {
		return &Creds{
			Credentials: &Credentials{APIKey: token, Label: "MIDAZ_TOKEN"},
			Profile:     "env",
		}, nil
	}
	s, err := Load()
	if err != nil {
		return nil, err
	}
	name := profile
	if name == "" {
		name = s.Current
	}
	c, ok := s.Profiles[name]
	if !ok || c == nil || c.APIKey == "" {
		return nil, nil
	}
	return &Creds{Credentials: c, Profile: name}, nil
}

// Creds bundles credentials with the profile name they came from.
type Creds struct {
	*Credentials
	Profile string
}

// SetCurrent writes credentials to the named profile (or DefaultProfile if
// empty) and updates the "current" pointer.
func SetCurrent(profile string, c *Credentials) error {
	if profile == "" {
		profile = DefaultProfile
	}
	s, err := Load()
	if err != nil {
		return err
	}
	if c.VerifiedAt == "" {
		c.VerifiedAt = time.Now().UTC().Format(time.RFC3339)
	}
	s.Profiles[profile] = c
	s.Current = profile
	return Save(s)
}

// Clear removes the named profile. If it was the current profile, current
// reverts to DefaultProfile (which may not exist — caller's problem).
func Clear(profile string) error {
	if profile == "" {
		profile = DefaultProfile
	}
	s, err := Load()
	if err != nil {
		return err
	}
	delete(s.Profiles, profile)
	if s.Current == profile {
		s.Current = DefaultProfile
	}
	return Save(s)
}

// MaskKey returns a display-safe prefix of the PAT (first 11 chars, e.g.
// "sk_12345678..."). Empty string when key is empty.
func MaskKey(key string) string {
	if len(key) < 12 {
		return ""
	}
	return key[:11] + "…"
}

// WriteMetadata appends auth state to an io.Writer (masked). Used by doctor.
func WriteMetadata(w io.Writer, c *Creds) {
	if c == nil {
		fmt.Fprintln(w, "auth: not logged in")
		return
	}
	fmt.Fprintf(w, "auth: %s @ %s (%s)\n",
		NonEmpty(c.UserEmail, "unknown"),
		NonEmpty(c.DeskSlug, c.DeskID),
		MaskKey(c.APIKey),
	)
}

// NonEmpty returns a if non-empty, otherwise b.
func NonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}
