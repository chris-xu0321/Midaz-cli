package main

import (
	"io/fs"
	"strings"
	"testing"

	agentres "github.com/chris-xu0321/Midaz-cli/agent"
	skillsdata "github.com/chris-xu0321/Midaz-cli/skills"
	"github.com/chris-xu0321/Midaz-cli/internal/registry"
)

// allSkillContent returns the combined content of all SKILL.md files
// discovered via the skills/*/SKILL.md tree pattern.
func allSkillContent(t *testing.T) string {
	t.Helper()

	matches, err := fs.Glob(skillsdata.FS, "*/SKILL.md")
	if err != nil {
		t.Fatalf("failed to glob skills: %v", err)
	}
	if len(matches) == 0 {
		t.Fatal("no skills found in embedded FS")
	}

	var combined strings.Builder
	for _, p := range matches {
		data, err := skillsdata.FS.ReadFile(p)
		if err != nil {
			t.Fatalf("failed to read %s from skills.FS: %v", p, err)
		}
		combined.Write(data)
		combined.WriteByte('\n')
	}
	return combined.String()
}

// TestSkillCoversRegistryCommands asserts that every command in the registry
// is mentioned in at least one SKILL.md file. Prevents drift.
func TestSkillCoversRegistryCommands(t *testing.T) {
	content := allSkillContent(t)

	for _, cmd := range registry.Commands {
		// Skip deprecated commands — they should not be taught to agents
		if strings.Contains(cmd.Description, "[deprecated]") {
			continue
		}
		pattern := "seer-q " + cmd.Name
		if !strings.Contains(content, pattern) {
			t.Errorf("skills do not mention command %q (expected pattern: %q)", cmd.Name, pattern)
		}
	}
}

// TestCmdCoversRegistryCommands asserts that every command in the registry
// is mentioned in seer.md (the slash command file). Prevents drift.
func TestCmdCoversRegistryCommands(t *testing.T) {
	cmdContent := string(agentres.SeerCmdMD)

	for _, cmd := range registry.Commands {
		// Skip deprecated commands — they should not be taught to agents
		if strings.Contains(cmd.Description, "[deprecated]") {
			continue
		}
		pattern := "seer-q " + cmd.Name
		if !strings.Contains(cmdContent, pattern) {
			t.Errorf("seer.md does not mention command %q (expected pattern: %q)", cmd.Name, pattern)
		}
	}
}

// TestSkillsHaveRequiredSections verifies structural completeness across the skill set.
func TestSkillsHaveRequiredSections(t *testing.T) {
	content := allSkillContent(t)

	requiredSections := []string{
		"## Response Format",     // seer-shared
		"## Command Reference",   // seer-market
		"## Query Strategy",      // seer-market
		"## Key Response Fields", // seer-market
		"## Examples",            // seer-market
	}

	for _, section := range requiredSections {
		if !strings.Contains(content, section) {
			t.Errorf("skills missing required section: %q", section)
		}
	}
}
