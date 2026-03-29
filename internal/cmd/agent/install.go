package agent

import (
	"os"
	"path/filepath"

	agentres "github.com/SparkssL/seer-cli/agent"
	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/SparkssL/seer-cli/internal/output"
	"github.com/spf13/cobra"
)

var targets = []struct {
	RelPath string
	Content *[]byte
}{
	{"skills/seer-market/SKILL.md", &agentres.SeerSkillMD},
	{"cmd/seer.md", &agentres.SeerCmdMD},
}

func newCmdInstall(f *cmdutil.Factory) *cobra.Command {
	var workspace string

	cmd := &cobra.Command{
		Use:   "install <platform>",
		Short: "Install agent files to workspace",
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

			var installed []string
			for _, t := range targets {
				dest := filepath.Join(workspace, ".claude", t.RelPath)
				if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
					return output.ErrConfig("failed to create directory: %s", err)
				}
				if err := os.WriteFile(dest, *t.Content, 0644); err != nil {
					return output.ErrConfig("failed to write file: %s", err)
				}
				installed = append(installed, ".claude/"+t.RelPath)
			}

			data := map[string]any{"installed": installed}
			meta := map[string]any{"workspace": workspace}
			return output.WriteSuccess(opts.Out, data, meta, opts.Format)
		},
	}

	cmd.Flags().StringVar(&workspace, "workspace", "", "Target workspace (default: current directory)")
	return cmd
}
