package onboard

import (
	"encoding/json"
	"os"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdGenerate(f *cmdutil.Factory) *cobra.Command {
	var (
		mode     string
		fromFile string
		yes      bool
	)
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "AI-generate radar + playbook and commit them",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"onboard generate requires --yes",
					"run: midaz onboard generate --mode guided --from-file input.json --yes")
			}
			if mode != "guided" && mode != "freeform" {
				return output.ErrValidation("--mode must be guided or freeform")
			}
			if fromFile == "" {
				return output.ErrValidation("--from-file is required")
			}
			raw, err := os.ReadFile(fromFile)
			if err != nil {
				return output.ErrConfig("cannot read %s: %s", fromFile, err)
			}
			var payload map[string]any
			if err := json.Unmarshal(raw, &payload); err != nil {
				return output.ErrValidation("--from-file must be valid JSON: %s", err)
			}
			payload["mode"] = mode
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "POST",
				Path:      "/api/desk/onboard/generate",
				Body:      payload,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&mode, "mode", "", "guided | freeform (required)")
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Path to JSON input matching the mode schema (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm onboarding")
	return cmd
}
