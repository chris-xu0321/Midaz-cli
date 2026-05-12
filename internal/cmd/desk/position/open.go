package position

import (
	"net/http"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdOpen(f *cmdutil.Factory) *cobra.Command {
	var (
		asset     string
		direction string
		thesis    string
		yes       bool
	)
	cmd := &cobra.Command{
		Use:   "open",
		Short: "Open a position (asset + bias direction + entry thesis)",
		Long: `Open a position on your desk.

Each position records the asset you're trading, the bias direction
(long or short), and a short entry thesis (≤1200 chars). Server side
enforces one open position per asset — a second open against the same
asset returns 409.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk position open requires --yes",
					`e.g. midaz desk position open --asset NVDA --direction long --thesis "..." --yes`)
			}
			if asset == "" {
				return output.ErrValidation("--asset is required (e.g. --asset NVDA)")
			}
			if direction != "long" && direction != "short" {
				return output.ErrValidation("--direction must be long or short")
			}
			if thesis == "" {
				return output.ErrValidation("--thesis is required")
			}
			if len([]rune(thesis)) > 1200 {
				return output.ErrValidation("--thesis must be ≤1200 characters")
			}
			creds, err := cmdutil.RequireAuth(f)
			if err != nil {
				return err
			}
			slug, err := resolveDeskSlug(cmd.Context(), f, creds)
			if err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method: http.MethodPost,
				Path:   "/api/desks/" + url.PathEscape(slug) + "/positions",
				Body: map[string]any{
					"asset":          asset,
					"bias_direction": direction,
					"entry_thesis":   thesis,
				},
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&asset, "asset", "", "Asset ID (e.g. NVDA)")
	cmd.Flags().StringVar(&direction, "direction", "", "Bias direction: long or short")
	cmd.Flags().StringVar(&thesis, "thesis", "", "Entry thesis (≤1200 chars)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm position open")
	return cmd
}
