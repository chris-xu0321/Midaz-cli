package assets

import (
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdThesis(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "thesis <asset_id> <thesis_id>",
		Short: "Drill-down: one (asset, thesis) pair",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return output.ErrValidation("usage: midaz assets thesis <asset_id> <thesis_id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			path := "/api/assets/" + url.PathEscape(args[0]) + "/theses/" + url.PathEscape(args[1])
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      path,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}
