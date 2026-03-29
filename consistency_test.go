package main

import (
	"strings"
	"testing"

	agentres "github.com/SparkssL/seer-cli/agent"
	"github.com/SparkssL/seer-cli/internal/registry"
)

// TestSkillCoversRegistryCommands asserts that every command in the registry
// is mentioned in the SKILL.md file. Prevents drift.
func TestSkillCoversRegistryCommands(t *testing.T) {
	skillContent := string(agentres.SeerSkillMD)

	for _, cmd := range registry.Commands {
		pattern := "seer-q " + cmd.Name
		if !strings.Contains(skillContent, pattern) {
			t.Errorf("SKILL.md does not mention command %q (expected pattern: %q)", cmd.Name, pattern)
		}
	}
}

// TestCmdCoversRegistryCommands asserts that every command in the registry
// is mentioned in seer.md (the slash command file). Prevents drift.
func TestCmdCoversRegistryCommands(t *testing.T) {
	cmdContent := string(agentres.SeerCmdMD)

	for _, cmd := range registry.Commands {
		pattern := "seer-q " + cmd.Name
		if !strings.Contains(cmdContent, pattern) {
			t.Errorf("seer.md does not mention command %q (expected pattern: %q)", cmd.Name, pattern)
		}
	}
}

// TestSkillHasRequiredSections verifies structural completeness.
func TestSkillHasRequiredSections(t *testing.T) {
	skillContent := string(agentres.SeerSkillMD)

	requiredSections := []string{
		"## Response Format",
		"## Command Reference",
		"## Query Strategy",
		"## Key Response Fields",
		"## Examples",
	}

	for _, section := range requiredSections {
		if !strings.Contains(skillContent, section) {
			t.Errorf("SKILL.md missing required section: %q", section)
		}
	}
}
