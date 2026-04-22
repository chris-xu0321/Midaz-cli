// Package klines hosts `midaz klines` and `midaz klines <asset>`
// for price/candlestick time-series data.
package klines

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdKlines builds the klines command tree. The root lists assets that
// have kline data; a positional argument drills into one asset's history.
func NewCmdKlines(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "klines [asset_id]",
		Short: "List assets with kline data, or show history for one",
		Long: `Without arguments, lists assets that have kline (candlestick) coverage.
With an asset id, returns that asset's { history, latest } series.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if len(args) == 0 {
				return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
					Path:      "/api/klines",
					Normalize: cmdutil.NormalizePassthrough,
				})
			}
			return runGet(f, opts, args[0])
		},
	}
	return cmd
}
