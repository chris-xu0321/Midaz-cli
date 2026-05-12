package position

import (
	"net/http"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	var (
		direction string
		thesis    string
		yes       bool
	)
	cmd := &cobra.Command{
		Use:   "update <position-id>",
		Short: "Update an open position's bias direction or entry thesis",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrValidation("usage: midaz desk position update <position-id> [--direction] [--thesis] --yes")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk position update requires --yes",
					`e.g. midaz desk position update <id> --direction short --yes`)
			}
			body := map[string]any{}
			if direction != "" {
				if direction != "long" && direction != "short" {
					return output.ErrValidation("--direction must be long or short")
				}
				body["bias_direction"] = direction
			}
			if thesis != "" {
				if len([]rune(thesis)) > 1200 {
					return output.ErrValidation("--thesis must be ≤1200 characters")
				}
				body["entry_thesis"] = thesis
			}
			if len(body) == 0 {
				return output.ErrValidation("provide --direction and/or --thesis")
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
				Method:    http.MethodPatch,
				Path:      "/api/desks/" + url.PathEscape(slug) + "/positions/" + url.PathEscape(args[0]),
				Body:      body,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&direction, "direction", "", "New bias direction: long or short")
	cmd.Flags().StringVar(&thesis, "thesis", "", "New entry thesis (≤1200 chars)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm position update")
	return cmd
}
