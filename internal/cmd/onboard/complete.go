package onboard

import (
	"os"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdComplete(f *cmdutil.Factory) *cobra.Command {
	var (
		radarFile    string
		playbookFile string
		yes          bool
	)
	cmd := &cobra.Command{
		Use:   "complete",
		Short: "Commit caller-supplied radar + playbook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"onboard complete requires --yes",
					"run: midaz onboard complete --radar radar.md --playbook playbook.md --yes")
			}
			if radarFile == "" || playbookFile == "" {
				return output.ErrValidation("--radar and --playbook are required")
			}
			radar, err := os.ReadFile(radarFile)
			if err != nil {
				return output.ErrConfig("cannot read %s: %s", radarFile, err)
			}
			playbook, err := os.ReadFile(playbookFile)
			if err != nil {
				return output.ErrConfig("cannot read %s: %s", playbookFile, err)
			}
			if len(playbook) > 20_000 {
				return output.ErrValidation("playbook exceeds 20000 chars")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			body := map[string]any{
				"radar":    string(radar),
				"playbook": string(playbook),
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "POST",
				Path:      "/api/desk/onboard",
				Body:      body,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&radarFile, "radar", "", "Path to radar Markdown file (required)")
	cmd.Flags().StringVar(&playbookFile, "playbook", "", "Path to playbook Markdown file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm onboarding")
	return cmd
}
