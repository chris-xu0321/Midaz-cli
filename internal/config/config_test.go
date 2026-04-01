package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaults(t *testing.T) {
	cfg := Defaults()
	if cfg.APIURL != "https://www.midaz.xyz" {
		t.Errorf("expected default APIURL, got %q", cfg.APIURL)
	}
	if cfg.FrontendURL != "https://www.midaz.xyz" {
		t.Errorf("expected default FrontendURL, got %q", cfg.FrontendURL)
	}
	if cfg.Format != "json" {
		t.Errorf("expected default format 'json', got %q", cfg.Format)
	}
}

func TestLoadFromFile_Missing(t *testing.T) {
	cfg, err := LoadFromFile("/nonexistent/path/config.json")
	if err != nil {
		t.Fatal(err)
	}
	// Should return defaults when file is missing
	if cfg.APIURL != "https://www.midaz.xyz" {
		t.Errorf("expected default APIURL for missing file, got %q", cfg.APIURL)
	}
}

func TestSetKeyAndLoad(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")
	t.Setenv("SEER_CONFIG_PATH", cfgPath)

	// Set a key
	if err := SetKey("api_url", "http://test:9000"); err != nil {
		t.Fatal(err)
	}

	// Load and verify
	cfg, err := LoadFromFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.APIURL != "http://test:9000" {
		t.Errorf("expected 'http://test:9000', got %q", cfg.APIURL)
	}
}

func TestLoad_Precedence(t *testing.T) {
	tmpDir := t.TempDir()
	cfgPath := filepath.Join(tmpDir, "config.json")
	t.Setenv("SEER_CONFIG_PATH", cfgPath)

	// Write a config file with format=pretty
	if err := SetKey("format", "pretty"); err != nil {
		t.Fatal(err)
	}

	// Env var should override file
	t.Setenv("SEER_FORMAT", "json")
	cfg, err := Load("", "")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Format != "json" {
		t.Errorf("env should override file: expected 'json', got %q", cfg.Format)
	}

	// Flag should override env
	cfg, err = Load("", "pretty")
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Format != "pretty" {
		t.Errorf("flag should override env: expected 'pretty', got %q", cfg.Format)
	}
}

func TestConfigPath_EnvOverride(t *testing.T) {
	t.Setenv("SEER_CONFIG_PATH", "/custom/path.json")
	if ConfigPath() != "/custom/path.json" {
		t.Errorf("expected custom path, got %q", ConfigPath())
	}
}

func TestConfigPath_Default(t *testing.T) {
	t.Setenv("SEER_CONFIG_PATH", "")
	path := ConfigPath()
	if path == "" {
		t.Error("expected non-empty default config path")
	}
}

func TestSource(t *testing.T) {
	// Flag source
	if s := Source("api_url", "http://flag"); s != "flag" {
		t.Errorf("expected 'flag', got %q", s)
	}

	// Env source
	t.Setenv("SEER_API_URL", "http://env")
	if s := Source("api_url", ""); s != "env" {
		t.Errorf("expected 'env', got %q", s)
	}
	os.Unsetenv("SEER_API_URL")

	// Default source (no file, no env)
	t.Setenv("SEER_CONFIG_PATH", "/nonexistent/config.json")
	if s := Source("api_url", ""); s != "default" {
		t.Errorf("expected 'default', got %q", s)
	}
}
