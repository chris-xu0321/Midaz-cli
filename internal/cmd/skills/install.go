package skills

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	embedskills "github.com/SparkssL/Midaz-cli/skills"
	"github.com/spf13/cobra"
)

// target describes an agent skill directory.
type target struct {
	Name     string // "claude" or "codex"
	RootDir  string // e.g. ~/.claude
	SkillDir string // e.g. ~/.claude/skills
}

// installEntry records a single file action.
type installEntry struct {
	Skill  string `json:"skill"`
	Target string `json:"target"`
	Path   string `json:"path"`
	Action string `json:"action"` // "created", "updated", "skipped", "dry-run"
}

var validTargets = []string{"auto", "claude", "codex", "all"}

func newCmdInstall(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "install [auto|claude|codex|all]",
		Short: "Install skills to agent directories",
		Long: `Install the embedded skills into agent skill directories.

Targets:
  auto    Install to detected agent directories only (default)
  claude  Install to Claude Code skill directory
  codex   Install to Codex skill directory
  all     Install to all known agent directories`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)

			targetName := "auto"
			if len(args) > 0 {
				targetName = args[0]
			}
			if !isValidTarget(targetName) {
				return output.ErrValidation("unknown target %q, expected one of: auto, claude, codex, all", targetName)
			}

			yes, _ := cmd.Flags().GetBool("yes")
			force, _ := cmd.Flags().GetBool("force")
			dryRun, _ := cmd.Flags().GetBool("dry-run")
			skillDir, _ := cmd.Flags().GetString("skill-dir")

			if !dryRun && !yes {
				return output.ErrWithHint(
					output.ExitValidation,
					"confirmation_required",
					"skills install requires --yes flag",
					"run: midaz skills install "+targetName+" --yes",
				)
			}

			return runInstall(opts, targetName, force, dryRun, skillDir)
		},
	}

	cmd.Flags().Bool("yes", false, "Skip confirmation (required for non-dry-run)")
	cmd.Flags().Bool("force", false, "Overwrite existing skill files")
	cmd.Flags().Bool("dry-run", false, "Print what would be installed without writing")
	cmd.Flags().String("skill-dir", "", "Custom skill directory (bypasses target resolution)")

	return cmd
}

func runInstall(opts *cmdutil.RunOpts, targetName string, force, dryRun bool, skillDir string) error {
	// If --skill-dir is set, use it directly as a single target.
	if skillDir != "" {
		entries, err := installSkillsTo("custom", skillDir, force, dryRun)
		if err != nil {
			return err
		}
		data := map[string]any{
			"detected":  []string{"custom"},
			"installed": entries,
		}
		meta := countActions(entries)
		meta["targets"] = []string{"custom"}
		return output.WriteSuccess(opts.Out, data, meta, opts.Format)
	}

	targets := resolveTargets(targetName)

	var detected []string
	for _, t := range targets {
		detected = append(detected, t.Name)
	}
	if detected == nil {
		detected = []string{}
	}

	var allEntries []installEntry
	for _, t := range targets {
		entries, err := installSkillsTo(t.Name, t.SkillDir, force, dryRun)
		if err != nil {
			return err
		}
		allEntries = append(allEntries, entries...)
	}
	if allEntries == nil {
		allEntries = []installEntry{}
	}

	data := map[string]any{
		"detected":  detected,
		"installed": allEntries,
	}
	meta := countActions(allEntries)

	targetNames := make([]string, len(targets))
	for i, t := range targets {
		targetNames[i] = t.Name
	}
	meta["targets"] = targetNames

	if targetName == "auto" && len(targets) == 0 {
		meta["hint"] = "No agent directories detected. Run one of:\n  midaz skills install claude --yes\n  midaz skills install codex --yes\n  midaz skills install all --yes"
	}

	// Clean up legacy seer-* skill directories replaced by midaz-*.
	removed := []string{}
	if !dryRun {
		for _, t := range targets {
			for _, legacy := range []string{"seer-shared", "seer-market", "seer-api-explorer"} {
				legacyPath := filepath.Join(t.SkillDir, legacy)
				if dirExists(legacyPath) {
					if err := os.RemoveAll(legacyPath); err == nil {
						removed = append(removed, legacyPath)
					}
				}
			}
		}
	}
	if len(removed) > 0 {
		meta["removed_legacy"] = removed
	}

	return output.WriteSuccess(opts.Out, data, meta, opts.Format)
}

func resolveTargets(targetName string) []target {
	home := userHomeDir()
	known := []target{
		{Name: "claude", RootDir: filepath.Join(home, ".claude"), SkillDir: filepath.Join(home, ".claude", "skills")},
		{Name: "codex", RootDir: filepath.Join(home, ".codex"), SkillDir: filepath.Join(home, ".codex", "skills")},
	}

	switch targetName {
	case "all":
		return known
	case "claude":
		return []target{known[0]}
	case "codex":
		return []target{known[1]}
	case "auto":
		return detectTargets(known)
	}
	return nil
}

func detectTargets(known []target) []target {
	var detected []target
	seen := map[string]bool{}

	if os.Getenv("CLAUDECODE") == "1" {
		for _, t := range known {
			if t.Name == "claude" {
				detected = append(detected, t)
				seen["claude"] = true
			}
		}
	}
	if os.Getenv("AGENT") == "codex" {
		for _, t := range known {
			if t.Name == "codex" {
				detected = append(detected, t)
				seen["codex"] = true
			}
		}
	}

	for _, t := range known {
		if !seen[t.Name] && dirExists(t.RootDir) {
			detected = append(detected, t)
		}
	}

	return detected
}

func installSkillsTo(targetName, skillDir string, force, dryRun bool) ([]installEntry, error) {
	var entries []installEntry

	err := fs.WalkDir(embedskills.FS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		destPath := filepath.Join(skillDir, path)

		content, readErr := fs.ReadFile(embedskills.FS, path)
		if readErr != nil {
			return readErr
		}

		action := resolveAction(destPath, content, force)

		if dryRun {
			entries = append(entries, installEntry{
				Skill:  skillNameFromPath(path),
				Target: targetName,
				Path:   destPath,
				Action: "dry-run:" + action,
			})
			return nil
		}

		if action == "skipped" {
			entries = append(entries, installEntry{
				Skill:  skillNameFromPath(path),
				Target: targetName,
				Path:   destPath,
				Action: "skipped",
			})
			return nil
		}

		if mkErr := os.MkdirAll(filepath.Dir(destPath), 0755); mkErr != nil {
			return mkErr
		}
		if writeErr := os.WriteFile(destPath, content, 0644); writeErr != nil {
			return writeErr
		}

		entries = append(entries, installEntry{
			Skill:  skillNameFromPath(path),
			Target: targetName,
			Path:   destPath,
			Action: action,
		})
		return nil
	})

	return entries, err
}

func resolveAction(destPath string, newContent []byte, force bool) string {
	existing, err := os.ReadFile(destPath)
	if err != nil {
		return "created"
	}
	if bytes.Equal(existing, newContent) {
		return "skipped"
	}
	if force {
		return "updated"
	}
	return "skipped"
}

func countActions(entries []installEntry) map[string]any {
	created, updated, skipped := 0, 0, 0
	for _, e := range entries {
		switch e.Action {
		case "created":
			created++
		case "updated":
			updated++
		default:
			skipped++
		}
	}
	return map[string]any{
		"created": created,
		"updated": updated,
		"skipped": skipped,
	}
}

func skillNameFromPath(path string) string {
	dir := filepath.Dir(path)
	if dir == "." {
		return path
	}
	return filepath.ToSlash(dir)
}

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		if p := os.Getenv("USERPROFILE"); p != "" {
			return p
		}
	}
	if p := os.Getenv("HOME"); p != "" {
		return p
	}
	h, _ := os.UserHomeDir()
	return h
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func isValidTarget(name string) bool {
	return slices.Contains(validTargets, name)
}
