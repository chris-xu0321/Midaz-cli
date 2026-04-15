package subscription

import (
	"encoding/json"

	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdSubscription(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "subscription",
		Short: "Subscribe, manage billing, and check status",
	}
	cmd.AddCommand(newCmdStatus(f))
	cmd.AddCommand(newCmdStart(f))
	cmd.AddCommand(newCmdPortal(f))
	return cmd
}

// openStripeURL POSTs to a Stripe session endpoint, opens the returned URL in
// the user's browser, and writes the success envelope. Shared by `start` and
// `portal` since both have identical plumbing and differ only in path/messaging.
func openStripeURL(cmd *cobra.Command, f *cmdutil.Factory, path, message, missingURLErr string) error {
	opts := cmdutil.ResolveRunOpts(cmd, f)
	if _, err := cmdutil.RequireAuth(f); err != nil {
		return err
	}
	c, err := f.Client()
	if err != nil {
		return err
	}
	resp, err := c.Post(cmd.Context(), path, nil)
	if err != nil {
		return err
	}
	var body struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(resp.Body, &body); err != nil || body.URL == "" {
		return output.ErrAPI("api", "%s", missingURLErr)
	}
	_ = auth.OpenBrowser(body.URL)
	return output.WriteSuccess(opts.Out,
		map[string]any{"url": body.URL},
		map[string]any{"message": message, "view_url": body.URL},
		opts.Format)
}
