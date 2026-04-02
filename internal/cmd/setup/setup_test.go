package setup

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/config"
	"github.com/SparkssL/Midaz-cli/skills"
)

// testFactoryResult wraps a Factory with captured stdout/stderr.
type testFactoryResult struct {
	*cmdutil.Factory
	stdout *bytes.Buffer
	stderr *bytes.Buffer
}

func testFactory() *testFactoryResult {
	var stdout, stderr bytes.Buffer
	f := &cmdutil.Factory{
		IOStreams: &cmdutil.IOStreams{Out: &stdout, ErrOut: &stderr},
		Config:   func() (*config.Config, error) { return config.Defaults(), nil },
	}
	return &testFactoryResult{Factory: f, stdout: &stdout, stderr: &stderr}
}

// parseEnvelope decodes a JSON envelope and returns the top-level map.
func parseEnvelope(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("invalid JSON envelope: %v\nraw: %s", err, data)
	}
	return env
}

func TestSetupRequiresYes(t *testing.T) {
	// Without --yes and without --dry-run, the command should fail
	tf := testFactory()
	cmd := NewCmdSetup(tf.Factory)
	cmd.SetArgs([]string{"all"})
	cmd.SilenceUsage = true

	err := cmd.Execute()
	if err == nil {
		t.Fatal("expected error without --yes flag")
	}
	if err.Error() != "setup requires --yes flag" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestSetupDryRunNoYes(t *testing.T) {
	// --dry-run should work without --yes
	dir := t.TempDir()
	tf := testFactory()
	cmd := NewCmdSetup(tf.Factory)
	cmd.SetArgs([]string{"all", "--dry-run", "--skill-dir", dir})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("dry-run should succeed without --yes: %v", err)
	}

	// Verify no files were actually written
	entries, _ := os.ReadDir(dir)
	if len(entries) > 0 {
		t.Errorf("dry-run should not write files, found %d entries", len(entries))
	}
}

func TestSetupAllCreatesSkills(t *testing.T) {
	dir := t.TempDir()
	tf := testFactory()
	cmd := NewCmdSetup(tf.Factory)
	cmd.SetArgs([]string{"all", "--yes", "--skill-dir", dir})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("setup all --yes failed: %v", err)
	}

	// Verify all 3 skill directories were created
	expectedSkills := []string{"seer-shared", "seer-market", "seer-api-explorer"}
	for _, skill := range expectedSkills {
		path := filepath.Join(dir, skill, "SKILL.md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected skill file %s to exist", path)
		}
	}

	// Parse stdout envelope
	env := parseEnvelope(t, tf.stdout.Bytes())
	if env["ok"] != true {
		t.Error("expected ok=true")
	}
	meta := env["meta"].(map[string]any)
	if meta["created"].(float64) != 3 {
		t.Errorf("expected 3 created, got %v", meta["created"])
	}
}

func TestSetupForceOverwrites(t *testing.T) {
	dir := t.TempDir()

	// Create a skill file with different content
	skillDir := filepath.Join(dir, "seer-shared")
	os.MkdirAll(skillDir, 0755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("old content"), 0644)

	// Without --force: should skip
	tf1 := testFactory()
	cmd1 := NewCmdSetup(tf1.Factory)
	cmd1.SetArgs([]string{"all", "--yes", "--skill-dir", dir})
	if err := cmd1.Execute(); err != nil {
		t.Fatalf("setup without force failed: %v", err)
	}

	env1 := parseEnvelope(t, tf1.stdout.Bytes())
	meta1 := env1["meta"].(map[string]any)
	if meta1["skipped"].(float64) < 1 {
		t.Error("expected at least 1 skipped without --force")
	}

	// Verify old content preserved
	content, _ := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if string(content) != "old content" {
		t.Error("without --force, existing file should not be overwritten")
	}

	// With --force: should update
	tf2 := testFactory()
	cmd2 := NewCmdSetup(tf2.Factory)
	cmd2.SetArgs([]string{"all", "--yes", "--force", "--skill-dir", dir})
	if err := cmd2.Execute(); err != nil {
		t.Fatalf("setup with force failed: %v", err)
	}

	env2 := parseEnvelope(t, tf2.stdout.Bytes())
	meta2 := env2["meta"].(map[string]any)
	if meta2["updated"].(float64) < 1 {
		t.Error("expected at least 1 updated with --force")
	}

	// Verify content was overwritten
	content, _ = os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if string(content) == "old content" {
		t.Error("with --force, existing file should be overwritten")
	}
}

func TestSetupAutoEnvDetection(t *testing.T) {
	// Use a temp home so no real dirs are detected
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)

	// Set CLAUDECODE=1 to simulate running inside Claude Code
	t.Setenv("CLAUDECODE", "1")
	// Clear AGENT to avoid codex detection
	t.Setenv("AGENT", "")

	known := []target{
		{Name: "claude", RootDir: filepath.Join(tmpHome, ".claude"), SkillDir: filepath.Join(tmpHome, ".claude", "skills")},
		{Name: "codex", RootDir: filepath.Join(tmpHome, ".codex"), SkillDir: filepath.Join(tmpHome, ".codex", "skills")},
	}

	detected := detectTargets(known)
	if len(detected) != 1 {
		t.Fatalf("expected 1 detected target, got %d", len(detected))
	}
	if detected[0].Name != "claude" {
		t.Errorf("expected claude, got %s", detected[0].Name)
	}
}

func TestSetupAutoNoTargets(t *testing.T) {
	// Use a temp home with no agent dirs and no env signals
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("CLAUDECODE", "")
	t.Setenv("AGENT", "")

	tf := testFactory()
	cmd := NewCmdSetup(tf.Factory)
	cmd.SetArgs([]string{"auto", "--yes"})

	err := cmd.Execute()
	if err != nil {
		t.Fatalf("auto with no targets should succeed: %v", err)
	}

	env := parseEnvelope(t, tf.stdout.Bytes())
	data := env["data"].(map[string]any)
	detected := data["detected"].([]any)
	installed := data["installed"].([]any)

	if len(detected) != 0 {
		t.Errorf("expected 0 detected, got %d", len(detected))
	}
	if len(installed) != 0 {
		t.Errorf("expected 0 installed, got %d", len(installed))
	}

	meta := env["meta"].(map[string]any)
	if _, ok := meta["hint"]; !ok {
		t.Error("expected hint in meta when no targets detected")
	}
}

func TestSetupSkillContentMatches(t *testing.T) {
	dir := t.TempDir()
	tf := testFactory()
	cmd := NewCmdSetup(tf.Factory)
	cmd.SetArgs([]string{"all", "--yes", "--skill-dir", dir})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	// Compare each installed file against the embedded FS
	err := fs.WalkDir(skills.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		embedded, _ := fs.ReadFile(skills.FS, path)
		installed, readErr := os.ReadFile(filepath.Join(dir, path))
		if readErr != nil {
			t.Errorf("installed file missing: %s", path)
			return nil
		}
		if string(embedded) != string(installed) {
			t.Errorf("content mismatch for %s", path)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk error: %v", err)
	}
}
