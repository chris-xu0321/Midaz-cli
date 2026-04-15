package keys

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdCreate(f *cmdutil.Factory) *cobra.Command {
	var (
		label string
		yes   bool
	)
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new PAT (shown once)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"auth keys create requires --yes", "run: midaz auth keys create --label <name> --yes")
			}
			if label == "" {
				return output.ErrValidation("--label is required")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "POST",
				Path:      "/api/app/api-keys",
				Body:      map[string]any{"label": label},
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&label, "label", "", "Human-readable name for the key (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm PAT creation")
	return cmd
}
