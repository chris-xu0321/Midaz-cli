// Package delta hosts `midaz delta`.
package delta

import (
	"fmt"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdDelta builds the `delta` top-level command.
func NewCmdDelta(f *cmdutil.Factory) *cobra.Command {
	var hours int
	cmd := &cobra.Command{
		Use:   "delta",
		Short: "Recent claims + theses + topics from the last N hours",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			params := url.Values{}
			if hours > 0 {
				params.Set("hours", fmt.Sprintf("%d", hours))
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/delta",
				Params:    params,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().IntVar(&hours, "hours", 12, "Lookback hours (1-168)")
	return cmd
}
