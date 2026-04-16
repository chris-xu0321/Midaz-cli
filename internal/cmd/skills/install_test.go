package skills

import (
	"bytes"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/config"
	embedskills "github.com/SparkssL/Midaz-cli/skills"
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
		Config:    func() (*config.Config, error) { return config.Defaults(), nil },
	}
	return &testFactoryResult{Factory: f, stdout: &stdout, stderr: &stderr}
}

// runSkills executes `midaz skills <args...>` against a fresh command tree.
func runSkills(tf *testFactoryResult, args ...string) error {
	cmd := NewCmdSkills(tf.Factory)
	cmd.SetArgs(args)
	cmd.SilenceUsage = true
	cmd.SetOut(tf.stdout)
	cmd.SetErr(tf.stderr)
	return cmd.Execute()
}

func parseEnvelope(t *testing.T, data []byte) map[string]any {
	t.Helper()
	var env map[string]any
	if err := json.Unmarshal(data, &env); err != nil {
		t.Fatalf("invalid JSON envelope: %v\nraw: %s", err, data)
	}
	return env
}

func TestInstallRequiresYes(t *testing.T) {
	tf := testFactory()
	err := runSkills(tf, "install", "all")
	if err == nil {
		t.Fatal("expected error without --yes flag")
	}
	if err.Error() != "skills install requires --yes flag" {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestInstallDryRunNoYes(t *testing.T) {
	dir := t.TempDir()
	tf := testFactory()
	if err := runSkills(tf, "install", "all", "--dry-run", "--skill-dir", dir); err != nil {
		t.Fatalf("dry-run should succeed without --yes: %v", err)
	}

	entries, _ := os.ReadDir(dir)
	if len(entries) > 0 {
		t.Errorf("dry-run should not write files, found %d entries", len(entries))
	}
}

func TestInstallAllCreatesSkills(t *testing.T) {
	dir := t.TempDir()
	tf := testFactory()
	if err := runSkills(tf, "install", "all", "--yes", "--skill-dir", dir); err != nil {
		t.Fatalf("skills install all --yes failed: %v", err)
	}

	expectedSkills := []string{"midaz-shared", "midaz-market", "midaz-api-explorer", "midaz-account", "midaz-desk"}
	for _, skill := range expectedSkills {
		path := filepath.Join(dir, skill, "SKILL.md")
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected skill file %s to exist", path)
		}
	}

	env := parseEnvelope(t, tf.stdout.Bytes())
	if env["ok"] != true {
		t.Error("expected ok=true")
	}
	meta := env["meta"].(map[string]any)
	if meta["created"].(float64) != float64(len(expectedSkills)) {
		t.Errorf("expected %d created, got %v", len(expectedSkills), meta["created"])
	}
}

func TestInstallForceOverwrites(t *testing.T) {
	dir := t.TempDir()

	skillDir := filepath.Join(dir, "midaz-shared")
	os.MkdirAll(skillDir, 0755)
	os.WriteFile(filepath.Join(skillDir, "SKILL.md"), []byte("old content"), 0644)

	tf1 := testFactory()
	if err := runSkills(tf1, "install", "all", "--yes", "--skill-dir", dir); err != nil {
		t.Fatalf("install without force failed: %v", err)
	}

	env1 := parseEnvelope(t, tf1.stdout.Bytes())
	meta1 := env1["meta"].(map[string]any)
	if meta1["skipped"].(float64) < 1 {
		t.Error("expected at least 1 skipped without --force")
	}

	content, _ := os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if string(content) != "old content" {
		t.Error("without --force, existing file should not be overwritten")
	}

	tf2 := testFactory()
	if err := runSkills(tf2, "install", "all", "--yes", "--force", "--skill-dir", dir); err != nil {
		t.Fatalf("install with force failed: %v", err)
	}

	env2 := parseEnvelope(t, tf2.stdout.Bytes())
	meta2 := env2["meta"].(map[string]any)
	if meta2["updated"].(float64) < 1 {
		t.Error("expected at least 1 updated with --force")
	}

	content, _ = os.ReadFile(filepath.Join(skillDir, "SKILL.md"))
	if string(content) == "old content" {
		t.Error("with --force, existing file should be overwritten")
	}
}

func TestInstallAutoEnvDetection(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("CLAUDECODE", "1")
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

func TestInstallAutoNoTargets(t *testing.T) {
	tmpHome := t.TempDir()
	t.Setenv("HOME", tmpHome)
	t.Setenv("USERPROFILE", tmpHome)
	t.Setenv("CLAUDECODE", "")
	t.Setenv("AGENT", "")

	tf := testFactory()
	if err := runSkills(tf, "install", "auto", "--yes"); err != nil {
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

func TestInstallSkillContentMatches(t *testing.T) {
	dir := t.TempDir()
	tf := testFactory()
	if err := runSkills(tf, "install", "all", "--yes", "--skill-dir", dir); err != nil {
		t.Fatalf("install failed: %v", err)
	}

	err := fs.WalkDir(embedskills.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		embedded, _ := fs.ReadFile(embedskills.FS, path)
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
