package agent

import (
	"os"
	"path/filepath"

	agentres "github.com/SparkssL/seer-cli/agent"
	skillsdata "github.com/SparkssL/seer-cli/skills"
	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/SparkssL/seer-cli/internal/output"
	"github.com/spf13/cobra"
)

// Claude-specific command file (not a skill)
var cmdTarget = struct {
	RelPath string
	Content *[]byte
}{
	"cmd/seer.md", &agentres.SeerCmdMD,
}

func newCmdInstall(f *cmdutil.Factory) *cobra.Command {
	var workspace string

	cmd := &cobra.Command{
		Use:   "install <platform>",
		Short: "Install agent files to workspace (compatibility bridge)",
		Args:  exactPlatform(),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)

			if workspace == "" {
				var err error
				workspace, err = os.Getwd()
				if err != nil {
					return output.ErrConfig("failed to get working directory: %s", err)
				}
			}

			skills, err := discoverSkills()
			if err != nil {
				return output.ErrConfig("skill discovery failed: %s", err)
			}

			var installed []string

			// Install skills discovered from skills/*/SKILL.md tree
			for _, skill := range skills {
				content, err := skillsdata.FS.ReadFile(skill.Path)
				if err != nil {
					return output.ErrConfig("failed to read embedded skill %s: %s", skill.Dir, err)
				}
				dest := filepath.Join(workspace, ".claude", "skills", skill.Dir, "SKILL.md")
				if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
					return output.ErrConfig("failed to create directory: %s", err)
				}
				if err := os.WriteFile(dest, content, 0644); err != nil {
					return output.ErrConfig("failed to write file: %s", err)
				}
				installed = append(installed, ".claude/skills/"+skill.Dir+"/SKILL.md")
			}

			// Install Claude command file
			dest := filepath.Join(workspace, ".claude", cmdTarget.RelPath)
			if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
				return output.ErrConfig("failed to create directory: %s", err)
			}
			if err := os.WriteFile(dest, *cmdTarget.Content, 0644); err != nil {
				return output.ErrConfig("failed to write file: %s", err)
			}
			installed = append(installed, ".claude/"+cmdTarget.RelPath)

			data := map[string]any{"installed": installed}
			meta := map[string]any{"workspace": workspace, "deprecated": true}
			return output.WriteSuccess(opts.Out, data, meta, opts.Format)
		},
	}

	cmd.Flags().StringVar(&workspace, "workspace", "", "Target workspace (default: current directory)")
	return cmd
}
