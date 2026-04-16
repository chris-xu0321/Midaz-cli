package desk

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdRegenerate(f *cmdutil.Factory) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "regenerate",
		Short: "Trigger a manual personal-desk rebuild",
		Long: `Enqueue a manual rebuild of your personal desk.

Owner-only and subscription-gated. Returns { status: "queued" } on
success, or 409 with a refresh_id if a rebuild is already in flight.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk regenerate requires --yes",
					"e.g. midaz desk regenerate --yes")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "POST",
				Path:      "/api/desk/personal-desk/regenerate",
				Body:      nil,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the refresh request")
	return cmd
}
