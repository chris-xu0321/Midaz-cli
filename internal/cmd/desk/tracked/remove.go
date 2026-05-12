package tracked

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdRemove(f *cmdutil.Factory) *cobra.Command {
	var (
		items string
		yes   bool
	)
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove asset ids from the tracked set",
		Long: `Read current tracked assets, drop the given ids (case-insensitive),
and PATCH the remainder back.

Example:
  midaz desk tracked-assets remove --items "OLD,STALE" --yes`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk tracked-assets remove requires --yes",
					`e.g. midaz desk tracked-assets remove --items "TLT" --yes`)
			}
			parsed := parseAssetList(items)
			if len(parsed) == 0 {
				return output.ErrValidation("--items must list at least one asset id")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			c, err := f.Client()
			if err != nil {
				return err
			}
			current, err := fetchTracked(cmd.Context(), c)
			if err != nil {
				return err
			}
			remaining := removeFrom(current, parsed)
			if len(remaining) == 0 {
				return output.ErrValidation("refusing to clear tracked-assets — server requires at least one tracked asset; use `desk tracked-assets set` if you really want to overwrite")
			}
			resp, err := pushTracked(cmd.Context(), c, remaining)
			if err != nil {
				return err
			}
			meta := map[string]any{
				"removed":       parsed,
				"tracked_count": len(remaining),
			}
			return output.WriteSuccess(opts.Out, resp, meta, opts.Format)
		},
	}
	cmd.Flags().StringVar(&items, "items", "", "Comma-separated asset ids to remove")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm remove")
	return cmd
}
