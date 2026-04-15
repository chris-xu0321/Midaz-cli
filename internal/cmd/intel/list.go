package intel

import (
	"fmt"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdList(f *cmdutil.Factory) *cobra.Command {
	var limit int
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recent intel items",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			params := url.Values{}
			if limit > 0 {
				params.Set("limit", fmt.Sprintf("%d", limit))
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/intel",
				Params:    params,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().IntVar(&limit, "limit", 0, "Maximum items (1-200, default 50)")
	return cmd
}
