// Package config handles configuration loading for the midaz CLI.
//
// Precedence (highest wins): CLI flags > env vars > config file > defaults
//
// Config file location (via os.UserConfigDir()):
//
//	Windows: %APPDATA%\midaz\config.json
//	macOS:   ~/Library/Application Support/midaz/config.json
//	Linux:   ~/.config/midaz/config.json
//
// Override: MIDAZ_CONFIG_PATH (or legacy SEER_CONFIG_PATH) env var
//
// Env var mapping:
//
//	MIDAZ_API_URL       → api_url      (default: https://www.midaz.xyz)
//	MIDAZ_FRONTEND_URL  → frontend_url (default: https://www.midaz.xyz)
//	MIDAZ_FORMAT        → format       (default: json)
//	MIDAZ_TOKEN         → bearer token for CI/headless auth
//
// Legacy SEER_* env vars are read as a fallback; when both are set the MIDAZ_*
// value wins.
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"slices"
)

// Config holds the resolved CLI configuration.
type Config struct {
	APIURL      string `json:"api_url,omitempty"`
	FrontendURL string `json:"frontend_url,omitempty"`
	Format      string `json:"format,omitempty"`
}

// ValidKeys are the recognized config keys.
var ValidKeys = []string{"api_url", "frontend_url", "format"}

// Defaults returns a Config with default values.
func Defaults() *Config {
	return &Config{
		APIURL:      "https://www.midaz.xyz",
		FrontendURL: "https://www.midaz.xyz",
		Format:      "json",
	}
}

// envFirst returns the first non-empty value from the provided env var names.
func envFirst(names ...string) string {
	for _, n := range names {
		if v := os.Getenv(n); v != "" {
			return v
		}
	}
	return ""
}

// ConfigPath returns the preferred config file path (midaz/config.json).
// If MIDAZ_CONFIG_PATH or legacy SEER_CONFIG_PATH is set, that wins.
// Otherwise: <userConfigDir>/midaz/config.json.
func ConfigPath() string {
	if p := envFirst("MIDAZ_CONFIG_PATH", "SEER_CONFIG_PATH"); p != "" {
		return p
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(".", "midaz", "config.json")
	}
	return filepath.Join(dir, "midaz", "config.json")
}

// LegacyConfigPath returns the legacy `seer/config.json` path, used as a
// read-only fallback for one release to avoid forcing users to re-create
// their config.
func LegacyConfigPath() string {
	if p := os.Getenv("SEER_CONFIG_PATH"); p != "" {
		return p
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "seer", "config.json")
}

// LoadFromFile reads a config JSON file. Returns Defaults() if file does not exist.
func LoadFromFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Defaults(), nil
		}
		return nil, err
	}
	cfg := Defaults()
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Load resolves the full config with precedence: defaults → legacy file → file → env → flags.
// flagAPIURL and flagFormat are from CLI flags (empty string = not set).
func Load(flagAPIURL, flagFormat string) (*Config, error) {
	cfg := Defaults()

	switch data, err := os.ReadFile(ConfigPath()); {
	case err == nil:
		_ = json.Unmarshal(data, cfg)
	case errors.Is(err, os.ErrNotExist):
		if legacy := LegacyConfigPath(); legacy != "" {
			if data, err := os.ReadFile(legacy); err == nil {
				_ = json.Unmarshal(data, cfg)
			}
		}
	default:
		return nil, err
	}

	// Env vars override file (MIDAZ_* preferred, SEER_* deprecated fallback).
	if v := envFirst("MIDAZ_API_URL", "SEER_API_URL"); v != "" {
		cfg.APIURL = v
	}
	if v := envFirst("MIDAZ_FRONTEND_URL", "SEER_FRONTEND_URL"); v != "" {
		cfg.FrontendURL = v
	}
	if v := envFirst("MIDAZ_FORMAT", "SEER_FORMAT"); v != "" {
		cfg.Format = v
	}

	// Flags override env.
	if flagAPIURL != "" {
		cfg.APIURL = flagAPIURL
	}
	if flagFormat != "" {
		cfg.Format = flagFormat
	}

	return cfg, nil
}

// Save writes a config to the preferred config file path. Creates directories as needed.
func Save(cfg *Config) error {
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

// SetKey reads the config file, sets one key, and writes back.
// Creates the file if it doesn't exist.
func SetKey(key, value string) error {
	if !slices.Contains(ValidKeys, key) {
		return errors.New("unknown config key: " + key)
	}

	cfg, err := LoadFromFile(ConfigPath())
	if err != nil {
		cfg = Defaults()
	}

	switch key {
	case "api_url":
		if cfg.APIURL == value {
			return nil
		}
		cfg.APIURL = value
	case "frontend_url":
		if cfg.FrontendURL == value {
			return nil
		}
		cfg.FrontendURL = value
	case "format":
		if cfg.Format == value {
			return nil
		}
		cfg.Format = value
	}

	return Save(cfg)
}

// Source returns where a key's value is coming from: "flag", "env", "file", "legacy_file", or "default".
func Source(key, flagValue string) string {
	if flagValue != "" {
		return "flag"
	}

	var envKeys []string
	switch key {
	case "api_url":
		envKeys = []string{"MIDAZ_API_URL", "SEER_API_URL"}
	case "frontend_url":
		envKeys = []string{"MIDAZ_FRONTEND_URL", "SEER_FRONTEND_URL"}
	case "format":
		envKeys = []string{"MIDAZ_FORMAT", "SEER_FORMAT"}
	}
	for _, k := range envKeys {
		if os.Getenv(k) != "" {
			return "env"
		}
	}

	if fileHasKey(ConfigPath(), key) {
		return "file"
	}
	if legacy := LegacyConfigPath(); legacy != "" && fileHasKey(legacy, key) {
		return "legacy_file"
	}
	return "default"
}

func fileHasKey(path, key string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}
	var fileCfg Config
	if json.Unmarshal(data, &fileCfg) != nil {
		return false
	}
	switch key {
	case "api_url":
		return fileCfg.APIURL != ""
	case "frontend_url":
		return fileCfg.FrontendURL != ""
	case "format":
		return fileCfg.Format != ""
	}
	return false
}
