package subscription

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdStart(f *cmdutil.Factory) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Create a Stripe Checkout session (opens browser)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"subscription start requires --yes",
					"run: midaz subscription start --yes")
			}
			return openStripeURL(cmd, f, "/api/stripe/checkout",
				"Opened Stripe Checkout in your browser.",
				"Stripe returned no checkout URL")
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm creating a Stripe Checkout session")
	return cmd
}
