package position

import (
	"net/http"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdClose(f *cmdutil.Factory) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "close <position-id>",
		Short: "Close an open position",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrValidation("usage: midaz desk position close <position-id> --yes")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk position close requires --yes",
					`e.g. midaz desk position close <id> --yes`)
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
				Method:    http.MethodPost,
				Path:      "/api/desks/" + url.PathEscape(slug) + "/positions/" + url.PathEscape(args[0]) + "/close",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm position close")
	return cmd
}
