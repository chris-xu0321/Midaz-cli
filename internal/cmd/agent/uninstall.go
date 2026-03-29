package agent

import (
	"os"
	"path/filepath"

	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/SparkssL/seer-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdUninstall(f *cmdutil.Factory) *cobra.Command {
	var workspace string

	cmd := &cobra.Command{
		Use:   "uninstall <platform>",
		Short: "Remove installed agent files",
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

			var removed []string
			for _, t := range targets {
				dest := filepath.Join(workspace, ".claude", t.RelPath)
				if err := os.Remove(dest); err != nil && !os.IsNotExist(err) {
					return output.ErrConfig("failed to remove file: %s", err)
				}
				removed = append(removed, ".claude/"+t.RelPath)
			}

			data := map[string]any{"removed": removed}
			meta := map[string]any{"workspace": workspace}
			return output.WriteSuccess(opts.Out, data, meta, opts.Format)
		},
	}

	cmd.Flags().StringVar(&workspace, "workspace", "", "Target workspace (default: current directory)")
	return cmd
}
