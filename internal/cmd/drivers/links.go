package drivers

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdDriverLinks returns `midaz driver-links` — list causal edges between drivers.
func NewCmdDriverLinks(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "driver-links",
		Short: "List causal edges between drivers (sphere/radar graph)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/driver-links",
				Normalize: cmdutil.NormalizeBareArray,
			})
		},
	}
}
