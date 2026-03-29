package agent

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	agentres "github.com/SparkssL/seer-cli/agent"
	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/SparkssL/seer-cli/internal/config"
	"github.com/SparkssL/seer-cli/internal/output"
)

func testFactory(out, errOut *bytes.Buffer) *cmdutil.Factory {
	return &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{Out: out, ErrOut: errOut},
		Config:   func() (*config.Config, error) { return config.Defaults(), nil },
	}
}

func TestInstallCommand(t *testing.T) {
	workspace := t.TempDir()
	var stdout, stderr bytes.Buffer
	f := testFactory(&stdout, &stderr)

	cmd := newCmdInstall(f)
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"claude", "--workspace", workspace})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("install failed: %v", err)
	}

	// Verify files exist
	for _, tgt := range targets {
		dest := filepath.Join(workspace, ".claude", tgt.RelPath)
		data, err := os.ReadFile(dest)
		if err != nil {
			t.Fatalf("file not found: %s", tgt.RelPath)
		}
		if !bytes.Equal(data, *tgt.Content) {
			t.Errorf("content mismatch: %s", tgt.RelPath)
		}
	}

	// Verify JSON envelope
	var result map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON output: %v\n%s", err, stdout.String())
	}
	if result["ok"] != true {
		t.Error("expected ok=true")
	}
}

func TestUninstallCommand(t *testing.T) {
	workspace := t.TempDir()
	var stdout, stderr bytes.Buffer
	f := testFactory(&stdout, &stderr)

	// Install first
	installCmd := newCmdInstall(f)
	installCmd.SetOut(&bytes.Buffer{})
	installCmd.SetErr(&bytes.Buffer{})
	installCmd.SetArgs([]string{"claude", "--workspace", workspace})
	installCmd.Execute()

	// Uninstall
	uninstallCmd := newCmdUninstall(f)
	uninstallCmd.SetOut(&stdout)
	uninstallCmd.SetErr(&stderr)
	uninstallCmd.SetArgs([]string{"claude", "--workspace", workspace})

	if err := uninstallCmd.Execute(); err != nil {
		t.Fatalf("uninstall failed: %v", err)
	}

	// Verify files removed
	for _, tgt := range targets {
		dest := filepath.Join(workspace, ".claude", tgt.RelPath)
		if _, err := os.Stat(dest); !os.IsNotExist(err) {
			t.Errorf("file should be removed: %s", tgt.RelPath)
		}
	}
}

func TestDoctorCommand_Installed(t *testing.T) {
	workspace := t.TempDir()

	// Install first (use separate buffers)
	var installOut, installErr bytes.Buffer
	installF := testFactory(&installOut, &installErr)
	installCmd := newCmdInstall(installF)
	installCmd.SetOut(&installOut)
	installCmd.SetErr(&installErr)
	installCmd.SetArgs([]string{"claude", "--workspace", workspace})
	installCmd.Execute()

	// Doctor (fresh buffers)
	var stdout, stderr bytes.Buffer
	f := testFactory(&stdout, &stderr)
	doctorCmd := newCmdDoctor(f)
	doctorCmd.SetOut(&stdout)
	doctorCmd.SetErr(&stderr)
	doctorCmd.SetArgs([]string{"--workspace", workspace})

	if err := doctorCmd.Execute(); err != nil {
		t.Fatalf("doctor failed: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
		t.Fatalf("invalid JSON: %v\nraw output: %s", err, stdout.String())
	}

	data := result["data"].(map[string]any)
	checks := data["checks"].([]any)

	// File checks should pass (first 2 entries are file checks)
	for _, c := range checks[:2] {
		check := c.(map[string]any)
		if check["status"] != "pass" {
			t.Errorf("expected pass for %s, got %s: %s", check["file"], check["status"], check["message"])
		}
	}
}

func TestDoctorCommand_Empty(t *testing.T) {
	workspace := t.TempDir()
	var stdout, stderr bytes.Buffer
	f := testFactory(&stdout, &stderr)

	doctorCmd := newCmdDoctor(f)
	doctorCmd.SetOut(&stdout)
	doctorCmd.SetErr(&stderr)
	doctorCmd.SetArgs([]string{"--workspace", workspace})

	if err := doctorCmd.Execute(); err != nil {
		t.Fatalf("doctor failed: %v", err)
	}

	var result map[string]any
	json.Unmarshal(stdout.Bytes(), &result)
	data := result["data"].(map[string]any)
	checks := data["checks"].([]any)

	for _, c := range checks[:2] {
		check := c.(map[string]any)
		if check["status"] != "fail" {
			t.Errorf("expected fail for %s on empty workspace, got %s", check["file"], check["status"])
		}
	}
}

func TestInstallCommand_NoPlatform(t *testing.T) {
	var stdout, stderr bytes.Buffer
	f := testFactory(&stdout, &stderr)

	cmd := newCmdInstall(f)
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for missing platform")
	}
	var exitErr *output.ExitError
	if ok := isExitError(err, &exitErr); !ok || exitErr.Code != output.ExitValidation {
		t.Errorf("expected validation error (exit 2), got: %v", err)
	}
}

func TestInstallCommand_ExtraArgs(t *testing.T) {
	var stdout, stderr bytes.Buffer
	f := testFactory(&stdout, &stderr)

	cmd := newCmdInstall(f)
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"claude", "extra"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for extra args")
	}
}

func TestInstallCommand_BadPlatform(t *testing.T) {
	var stdout, stderr bytes.Buffer
	f := testFactory(&stdout, &stderr)

	cmd := newCmdInstall(f)
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs([]string{"cursor"})

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error for unknown platform")
	}
}

func TestExtractVersion(t *testing.T) {
	data := []byte("---\nname: test\nversion: 0.2.0\n---\ncontent")
	if v := extractVersion(data); v != "0.2.0" {
		t.Errorf("expected '0.2.0', got %q", v)
	}
	if v := extractVersion([]byte("no frontmatter")); v != "" {
		t.Errorf("expected empty, got %q", v)
	}
}

func TestEmbeddedContent(t *testing.T) {
	if len(agentres.SeerSkillMD) == 0 {
		t.Error("SeerSkillMD should not be empty")
	}
	if len(agentres.SeerCmdMD) == 0 {
		t.Error("SeerCmdMD should not be empty")
	}
}

func TestExactPlatform(t *testing.T) {
	v := exactPlatform()
	if err := v(nil, []string{"claude"}); err != nil {
		t.Errorf("claude should be valid: %v", err)
	}
	if err := v(nil, []string{}); err == nil {
		t.Error("empty args should fail")
	}
	if err := v(nil, []string{"cursor"}); err == nil {
		t.Error("unknown platform should fail")
	}
	if err := v(nil, []string{"claude", "extra"}); err == nil {
		t.Error("extra args should fail")
	}
}

// isExitError checks if err (or a wrapped error) is *output.ExitError.
func isExitError(err error, target **output.ExitError) bool {
	if e, ok := err.(*output.ExitError); ok {
		*target = e
		return true
	}
	return false
}
