package tracked

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdAdd(f *cmdutil.Factory) *cobra.Command {
	var (
		items string
		yes   bool
	)
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add asset ids to the tracked set (merge, dedupe)",
		Long: `Read current tracked assets, merge in the new ids, and PATCH back.

Example:
  midaz desk tracked-assets add --items "AMD,QQQ" --yes`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk tracked-assets add requires --yes",
					`e.g. midaz desk tracked-assets add --items "AMD,QQQ" --yes`)
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
			merged := mergeUnique(current, parsed)
			resp, err := pushTracked(cmd.Context(), c, merged)
			if err != nil {
				return err
			}
			meta := map[string]any{
				"added":         parsed,
				"tracked_count": len(merged),
			}
			return output.WriteSuccess(opts.Out, resp, meta, opts.Format)
		},
	}
	cmd.Flags().StringVar(&items, "items", "", "Comma-separated asset ids to add")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm add")
	return cmd
}
