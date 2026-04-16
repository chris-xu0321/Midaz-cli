package desk

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdView(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Personal market read (subscription-gated)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/desk/view",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}
