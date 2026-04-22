package assets

import (
	"net/url"
	"strconv"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdTimeline(f *cmdutil.Factory) *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "timeline <asset_id>",
		Short: "Asset event timeline (claims, deltas, signal changes)",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrValidation("usage: midaz assets timeline <asset_id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			params := url.Values{}
			if limit > 0 {
				params.Set("limit", strconv.Itoa(limit))
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/assets/" + url.PathEscape(args[0]) + "/timeline",
				Params:    params,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 0, "Max timeline entries (server default if unset)")
	return cmd
}
