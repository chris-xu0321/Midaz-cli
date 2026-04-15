package subscription

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdPortal(f *cmdutil.Factory) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "portal",
		Short: "Open the Stripe Customer Portal",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"subscription portal requires --yes",
					"run: midaz subscription portal --yes")
			}
			return openStripeURL(cmd, f, "/api/stripe/portal",
				"Opened Stripe Customer Portal in your browser.",
				"Stripe returned no portal URL")
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm opening the billing portal")
	return cmd
}
