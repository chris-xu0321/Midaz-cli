package agent

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"

	skillsdata "github.com/SparkssL/seer-cli/skills"
	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/SparkssL/seer-cli/internal/output"
	"github.com/spf13/cobra"
)

type fileCheck struct {
	File    string `json:"file"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func newCmdDoctor(f *cmdutil.Factory) *cobra.Command {
	var workspace string

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check agent file status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)

			if workspace == "" {
				var err error
				workspace, err = os.Getwd()
				if err != nil {
					return output.ErrConfig("failed to get working directory: %s", err)
				}
			}

			var checks []fileCheck

			// Check installed skills (discovered from tree)
			skills, err := discoverSkills()
			if err != nil {
				return output.ErrConfig("skill discovery failed: %s", err)
			}

			for _, skill := range skills {
				relPath := ".claude/skills/" + skill.Dir + "/SKILL.md"
				dest := filepath.Join(workspace, relPath)

				source, err := skillsdata.FS.ReadFile(skill.Path)
				if err != nil {
					checks = append(checks, fileCheck{relPath, "fail", "embedded source not found"})
					continue
				}

				data, err := os.ReadFile(dest)
				if err != nil {
					checks = append(checks, fileCheck{relPath, "fail", "not found"})
					continue
				}
				if bytes.Equal(data, source) {
					msg := "up to date"
					if fm, fmErr := parseFrontmatter(data); fmErr == nil && fm.Version != "" {
						msg = "v" + fm.Version + ", up to date"
					}
					checks = append(checks, fileCheck{relPath, "pass", msg})
				} else {
					msg := "exists but differs from source"
					if fm, fmErr := parseFrontmatter(data); fmErr == nil && fm.Version != "" {
						msg = "v" + fm.Version + ", differs from source"
					}
					checks = append(checks, fileCheck{relPath, "warn", msg})
				}
			}

			// Check Claude command file
			{
				relPath := ".claude/" + cmdTarget.RelPath
				dest := filepath.Join(workspace, relPath)

				data, err := os.ReadFile(dest)
				if err != nil {
					checks = append(checks, fileCheck{relPath, "fail", "not found"})
				} else if bytes.Equal(data, *cmdTarget.Content) {
					checks = append(checks, fileCheck{relPath, "pass", "up to date"})
				} else {
					checks = append(checks, fileCheck{relPath, "warn", "exists but differs from source"})
				}
			}

			// Check seer-q binary on PATH
			if path, err := exec.LookPath("seer-q"); err == nil {
				checks = append(checks, fileCheck{"seer-q", "pass", path})
			} else {
				checks = append(checks, fileCheck{"seer-q", "warn", "not found on PATH"})
			}

			result := map[string]any{"checks": checks}
			meta := map[string]any{"workspace": workspace, "deprecated": true}
			return output.WriteSuccess(opts.Out, result, meta, opts.Format)
		},
	}

	cmd.Flags().StringVar(&workspace, "workspace", "", "Target workspace (default: current directory)")
	return cmd
}
