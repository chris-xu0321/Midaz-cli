package workspace

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdSettings(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "Owner-only settings: radar, playbook, Telegram",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/ws/settings",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}
