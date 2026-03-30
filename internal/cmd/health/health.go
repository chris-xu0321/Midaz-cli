package health

import (
	"github.com/chris-xu0321/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdHealth(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "health",
		Short: "API health check",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/health",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}
