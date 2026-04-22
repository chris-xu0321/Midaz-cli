package desk

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdRefresh(f *cmdutil.Factory) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "refresh",
		Short: "Trigger a full pipeline refresh for your desk",
		Long: `Enqueue a full pipeline refresh — rebuilds the whole market view from
source ingestion through to your personal desk.

Use ` + "`regenerate`" + ` instead when you only want to rebuild your personal
desk against the existing market refresh; ` + "`refresh`" + ` rebuilds the
market itself first, which is slower and shared across all desks.

Owner-only and subscription-gated. Returns { status: "queued" }.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk refresh requires --yes",
					"e.g. midaz desk refresh --yes")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "POST",
				Path:      "/api/desk/refresh",
				Body:      nil,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the refresh request")
	return cmd
}
