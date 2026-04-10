// Package config handles configuration loading for the seer-q CLI.
//
// Precedence (highest wins): CLI flags > env vars > config file > defaults
//
// Config file location (via os.UserConfigDir()):
//
//	Windows: %APPDATA%\seer\config.json
//	macOS:   ~/Library/Application Support/seer/config.json
//	Linux:   ~/.config/seer/config.json
//
// Override: SEER_CONFIG_PATH env var
//
// Env var mapping:
//
//	SEER_API_URL       → api_url      (default: https://www.midaz.xyz)
//	SEER_FRONTEND_URL  → frontend_url (default: https://www.midaz.xyz)
//	SEER_FORMAT        → format       (default: json)
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
	APIKey      string `json:"api_key,omitempty"`
}

// ValidKeys are the recognized config keys.
var ValidKeys = []string{"api_url", "frontend_url", "format", "api_key"}

// Defaults returns a Config with default values.
func Defaults() *Config {
	return &Config{
		APIURL:      "https://www.midaz.xyz",
		FrontendURL: "https://www.midaz.xyz",
		Format:      "json",
	}
}

// ConfigPath returns the config file path.
// Uses SEER_CONFIG_PATH env var if set, otherwise os.UserConfigDir()/seer/config.json.
func ConfigPath() string {
	if p := os.Getenv("SEER_CONFIG_PATH"); p != "" {
		return p
	}
	dir, err := os.UserConfigDir()
	if err != nil {
		return filepath.Join(".", "seer", "config.json")
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

// Load resolves the full config with precedence: defaults → file → env → flags.
// flagAPIURL and flagFormat are from CLI flags (empty string = not set).
func Load(flagAPIURL, flagFormat string) (*Config, error) {
	cfg, err := LoadFromFile(ConfigPath())
	if err != nil {
		return nil, err
	}

	// Env vars override file
	if v := os.Getenv("SEER_API_URL"); v != "" {
		cfg.APIURL = v
	}
	if v := os.Getenv("SEER_FRONTEND_URL"); v != "" {
		cfg.FrontendURL = v
	}
	if v := os.Getenv("SEER_FORMAT"); v != "" {
		cfg.Format = v
	}
	if v := os.Getenv("SEER_API_KEY"); v != "" {
		cfg.APIKey = v
	}

	// Flags override env
	if flagAPIURL != "" {
		cfg.APIURL = flagAPIURL
	}
	if flagFormat != "" {
		cfg.Format = flagFormat
	}

	return cfg, nil
}

// Save writes a config to the config file path. Creates directories as needed.
func Save(cfg *Config) error {
	path := ConfigPath()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
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
		cfg.APIURL = value
	case "frontend_url":
		cfg.FrontendURL = value
	case "format":
		cfg.Format = value
	case "api_key":
		cfg.APIKey = value
	}

	return Save(cfg)
}

// Source returns where a key's value is coming from: "flag", "env", "file", or "default".
func Source(key, flagValue string) string {
	if flagValue != "" {
		return "flag"
	}

	envKey := ""
	switch key {
	case "api_url":
		envKey = "SEER_API_URL"
	case "frontend_url":
		envKey = "SEER_FRONTEND_URL"
	case "format":
		envKey = "SEER_FORMAT"
	case "api_key":
		envKey = "SEER_API_KEY"
	}
	if envKey != "" && os.Getenv(envKey) != "" {
		return "env"
	}

	// Check if file has a non-default value
	path := ConfigPath()
	if data, err := os.ReadFile(path); err == nil {
		var fileCfg Config
		if json.Unmarshal(data, &fileCfg) == nil {
			switch key {
			case "api_url":
				if fileCfg.APIURL != "" {
					return "file"
				}
			case "frontend_url":
				if fileCfg.FrontendURL != "" {
					return "file"
				}
			case "format":
				if fileCfg.Format != "" {
					return "file"
				}
			case "api_key":
				if fileCfg.APIKey != "" {
					return "file"
				}
			}
		}
	}

	return "default"
}
