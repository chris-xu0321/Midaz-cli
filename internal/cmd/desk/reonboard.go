package desk

import (
	"encoding/json"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdReonboard(f *cmdutil.Factory) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "reonboard",
		Short: "Re-submit current radar + playbook to trigger a desk rebuild",
		Long: `Reads your current radar and playbook and re-submits them unchanged,
triggering a desk rebuild. Distinct from the top-level ` + "`midaz onboard`" + `,
which is for first-time setup with user-supplied files.

Subscription-gated in practice because the settings read requires an
active subscription.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk reonboard requires --yes",
					"e.g. midaz desk reonboard --yes")
			}
			c, _, err := cmdutil.AuthedClient(f)
			if err != nil {
				return err
			}

			settingsResp, err := c.Get(opts.Ctx, "/api/desk/settings", nil)
			if err != nil {
				return err
			}
			var settings struct {
				RadarItems []string `json:"radar_items"`
				Playbook   string   `json:"playbook"`
			}
			if err := json.Unmarshal(settingsResp.Body, &settings); err != nil {
				return output.Errorf(output.ExitInternal, "internal",
					"failed to parse desk settings: %s", err)
			}
			if len(settings.RadarItems) == 0 || settings.Playbook == "" {
				return output.ErrWithHint(output.ExitValidation, "not_onboarded",
					"desk has no radar or playbook to re-submit",
					"run 'midaz onboard' to complete initial onboarding")
			}

			onboardResp, err := c.Post(opts.Ctx, "/api/desk/onboard", map[string]any{
				"radar_items": settings.RadarItems,
				"playbook":    settings.Playbook,
			})
			if err != nil {
				return err
			}

			if opts.Raw {
				return output.WriteRaw(opts.Out, onboardResp.Body, opts.Format)
			}
			data, meta, err := cmdutil.NormalizePassthrough(onboardResp.Body)
			if err != nil {
				return output.Errorf(output.ExitInternal, "internal",
					"failed to parse onboard response: %s", err)
			}
			return output.WriteSuccess(opts.Out, data, meta, opts.Format)
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the refresh request")
	return cmd
}
