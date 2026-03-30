package agent

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	agentres "github.com/chris-xu0321/Midaz-cli/agent"
	skillsdata "github.com/chris-xu0321/Midaz-cli/skills"
	"github.com/chris-xu0321/Midaz-cli/internal/cmdutil"
	"github.com/chris-xu0321/Midaz-cli/internal/config"
	"github.com/chris-xu0321/Midaz-cli/internal/output"
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

	// Verify skill files exist and match embedded content
	skills, err := discoverSkills()
	if err != nil {
		t.Fatalf("discovery failed: %v", err)
	}
	for _, skill := range skills {
		dest := filepath.Join(workspace, ".claude", "skills", skill.Dir, "SKILL.md")
		data, err := os.ReadFile(dest)
		if err != nil {
			t.Fatalf("skill file not found: %s", skill.Dir)
		}
		expected, _ := skillsdata.FS.ReadFile(skill.Path)
		if !bytes.Equal(data, expected) {
			t.Errorf("content mismatch: %s", skill.Dir)
		}
	}

	// Verify cmd/seer.md exists and matches
	cmdDest := filepath.Join(workspace, ".claude", cmdTarget.RelPath)
	data, err := os.ReadFile(cmdDest)
	if err != nil {
		t.Fatalf("cmd file not found: %s", cmdTarget.RelPath)
	}
	if !bytes.Equal(data, *cmdTarget.Content) {
		t.Errorf("content mismatch: %s", cmdTarget.RelPath)
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

	// Verify skill files removed
	skills, _ := discoverSkills()
	for _, skill := range skills {
		dest := filepath.Join(workspace, ".claude", "skills", skill.Dir, "SKILL.md")
		if _, err := os.Stat(dest); !os.IsNotExist(err) {
			t.Errorf("skill file should be removed: %s", skill.Dir)
		}
	}

	// Verify cmd file removed
	cmdDest := filepath.Join(workspace, ".claude", cmdTarget.RelPath)
	if _, err := os.Stat(cmdDest); !os.IsNotExist(err) {
		t.Errorf("cmd file should be removed: %s", cmdTarget.RelPath)
	}
}

func TestDoctorCommand_Installed(t *testing.T) {
	workspace := t.TempDir()

	// Install first
	var installOut, installErr bytes.Buffer
	installF := testFactory(&installOut, &installErr)
	installCmd := newCmdInstall(installF)
	installCmd.SetOut(&installOut)
	installCmd.SetErr(&installErr)
	installCmd.SetArgs([]string{"claude", "--workspace", workspace})
	installCmd.Execute()

	// Doctor
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

	// File checks: N discovered skills + 1 cmd, should all pass
	skills, _ := discoverSkills()
	expectedFiles := len(skills) + 1
	for i := 0; i < expectedFiles && i < len(checks); i++ {
		check := checks[i].(map[string]any)
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

	// File checks should all fail on empty workspace
	skills, _ := discoverSkills()
	expectedFiles := len(skills) + 1
	for i := 0; i < expectedFiles && i < len(checks); i++ {
		check := checks[i].(map[string]any)
		if check["status"] != "fail" {
			t.Errorf("expected fail for %s on empty workspace, got %s", check["file"], check["status"])
		}
	}
}

func TestDiscoverSkills(t *testing.T) {
	skills, err := discoverSkills()
	if err != nil {
		t.Fatalf("discovery failed: %v", err)
	}
	if len(skills) != 3 {
		t.Fatalf("expected 3 skills, got %d", len(skills))
	}

	// Verify each discovered skill has valid frontmatter
	for _, skill := range skills {
		if skill.Name == "" {
			t.Errorf("skill %s has empty name", skill.Dir)
		}
		if skill.Version == "" {
			t.Errorf("skill %s has empty version", skill.Dir)
		}
		if skill.Name != skill.Dir {
			t.Errorf("skill name %q does not match dir %q", skill.Name, skill.Dir)
		}

		// Verify content is non-empty in FS
		data, err := skillsdata.FS.ReadFile(skill.Path)
		if err != nil {
			t.Errorf("skill not found in FS: %s", skill.Path)
		}
		if len(data) == 0 {
			t.Errorf("skill should not be empty: %s", skill.Path)
		}
	}
}

func TestParseFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		wantFM  *frontmatter
	}{
		{
			name:  "valid",
			input: "---\nname: test-skill\nversion: 1.0.0\ndescription: A test skill\n---\n# Content",
			wantFM: &frontmatter{Name: "test-skill", Version: "1.0.0", Description: "A test skill"},
		},
		{
			name:    "missing name",
			input:   "---\nversion: 1.0.0\n---\n# Content",
			wantErr: true,
		},
		{
			name:    "missing version",
			input:   "---\nname: test\n---\n# Content",
			wantErr: true,
		},
		{
			name:    "no frontmatter block",
			input:   "# Just content\nNo frontmatter here",
			wantErr: true,
		},
		{
			name:    "unclosed block",
			input:   "---\nname: test\nversion: 1.0.0\n# No closing marker",
			wantErr: true,
		},
		{
			name:  "metadata line ignored",
			input: "---\nname: test-skill\nversion: 1.0.0\nmetadata: {\"requires\":{\"bins\":[\"seer-q\"]}}\n---\n# Content with name: fake",
			wantFM: &frontmatter{Name: "test-skill", Version: "1.0.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, err := parseFrontmatter([]byte(tt.input))
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fm.Name != tt.wantFM.Name {
				t.Errorf("name: got %q, want %q", fm.Name, tt.wantFM.Name)
			}
			if fm.Version != tt.wantFM.Version {
				t.Errorf("version: got %q, want %q", fm.Version, tt.wantFM.Version)
			}
			if tt.wantFM.Description != "" && fm.Description != tt.wantFM.Description {
				t.Errorf("description: got %q, want %q", fm.Description, tt.wantFM.Description)
			}
		})
	}
}

func TestEmbeddedCmdContent(t *testing.T) {
	if len(agentres.SeerCmdMD) == 0 {
		t.Error("SeerCmdMD should not be empty")
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
