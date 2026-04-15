package assets

import (
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdList(f *cmdutil.Factory) *cobra.Command {
	var (
		tier string
		bias string
	)
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List assets (with optional filters)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			params := url.Values{}
			if tier != "" {
				params.Set("tier", tier)
			}
			if bias != "" {
				params.Set("bias", bias)
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/assets",
				Params:    params,
				Normalize: cmdutil.NormalizeBareArray,
			})
		},
	}
	cmd.Flags().StringVar(&tier, "tier", "", "Filter by tier (1 or 2)")
	cmd.Flags().StringVar(&bias, "bias", "", "Filter by bias (bullish/bearish/neutral/mixed)")
	return cmd
}
