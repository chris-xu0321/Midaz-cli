package tracked

import (
	"os"
	"strings"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdSet(f *cmdutil.Factory) *cobra.Command {
	var (
		items    string
		fromFile string
		yes      bool
	)
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Replace the tracked asset list with the given comma-separated set",
		Long: `Overwrite the desk's tracked asset list.

Each id is uppercased and deduped. The server validates every id against
the active asset universe and rejects unknown ones with 400.

Examples:
  midaz desk tracked-assets set --items "NVDA,GLD,US10Y" --yes
  midaz desk tracked-assets set --from-file watchlist.txt --yes`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk tracked-assets set requires --yes",
					`e.g. midaz desk tracked-assets set --items "NVDA,GLD" --yes`)
			}
			if items == "" && fromFile == "" {
				return output.ErrValidation("provide --items or --from-file")
			}
			if items != "" && fromFile != "" {
				return output.ErrValidation("--items and --from-file are mutually exclusive")
			}
			raw := items
			if fromFile != "" {
				data, err := os.ReadFile(fromFile)
				if err != nil {
					return output.Errorf(output.ExitConfig, "config",
						"failed to read --from-file %s: %s", fromFile, err)
				}
				// File may use newlines or commas — normalize both into commas.
				raw = strings.ReplaceAll(string(data), "\n", ",")
			}
			parsed := parseAssetList(raw)
			if len(parsed) == 0 {
				return output.ErrValidation("no asset ids parsed from --items / --from-file")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "PATCH",
				Path:      "/api/desk/tracked-assets",
				Body:      map[string]any{"tracked_asset_ids": parsed},
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&items, "items", "", "Comma-separated asset ids (e.g. \"NVDA,GLD,US10Y\")")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Read asset ids from a file (comma- or newline-separated)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm overwrite")
	return cmd
}
