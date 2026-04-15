package radar

import (
	"strings"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// newCmdRemove builds `midaz workspace radar remove`.
//
// Exactly one of --index (1-based) or --match (case-insensitive substring)
// selects the item to drop. With --match, ambiguous matches error unless
// --first is passed.
func newCmdRemove(f *cmdutil.Factory) *cobra.Command {
	var (
		index int
		match string
		first bool
		yes   bool
	)
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove an item from the radar by index or substring match",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"workspace radar remove requires --yes",
					"e.g. --index 1 --yes")
			}
			indexSet := cmd.Flags().Changed("index")
			matchSet := strings.TrimSpace(match) != ""
			if indexSet == matchSet {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"specify exactly one of --index or --match",
					"e.g. --index 2 --yes, or --match AAPL --yes")
			}

			c, _, err := cmdutil.AuthedClient(f)
			if err != nil {
				return err
			}

			items, err := fetchCurrentItems(opts.Ctx, c)
			if err != nil {
				return err
			}
			if len(items) == 0 {
				return output.ErrValidation("radar is empty — nothing to remove")
			}

			var removeAt int
			var removed string
			if indexSet {
				if index < 1 || index > len(items) {
					return output.ErrValidation("--index %d is out of range (1..%d)", index, len(items))
				}
				removeAt = index - 1
				removed = items[removeAt]
			} else {
				needle := strings.ToLower(strings.TrimSpace(match))
				matches := []int{}
				for i, it := range items {
					if strings.Contains(strings.ToLower(it), needle) {
						matches = append(matches, i)
					}
				}
				if len(matches) == 0 {
					return output.ErrValidation("no radar item matches %q", match)
				}
				if len(matches) > 1 && !first {
					return output.ErrWithHint(output.ExitValidation, "ambiguous_match",
						"multiple radar items match the substring — refine --match or pass --first",
						"run 'midaz workspace radar get' to see all items")
				}
				removeAt = matches[0]
				removed = items[removeAt]
			}

			newItems := append(items[:removeAt:removeAt], items[removeAt+1:]...)
			resp, err := pushItems(opts.Ctx, c, newItems)
			if err != nil {
				return err
			}

			meta := map[string]any{
				"removed": removed,
				"count":   len(newItems),
			}
			return output.WriteSuccess(opts.Out, resp, meta, opts.Format)
		},
	}
	cmd.Flags().IntVar(&index, "index", 0, "1-based index of the radar item to remove")
	cmd.Flags().StringVar(&match, "match", "", "Case-insensitive substring — removes the first matching item")
	cmd.Flags().BoolVar(&first, "first", false, "With --match, remove the first match even if more than one item matches")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the removal")
	return cmd
}
