package workspace

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdShare(f *cmdutil.Factory) *cobra.Command {
	var (
		on  bool
		off bool
		yes bool
	)
	cmd := &cobra.Command{
		Use:   "share",
		Short: "Toggle workspace sharing (public /w/<id> page)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if on == off {
				return output.ErrValidation("specify exactly one of --on or --off")
			}
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"workspace share requires --yes", "e.g. midaz workspace share --on --yes")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "PATCH",
				Path:      "/api/ws",
				Body:      map[string]any{"shared": on},
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().BoolVar(&on, "on", false, "Enable sharing")
	cmd.Flags().BoolVar(&off, "off", false, "Disable sharing")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the change")
	return cmd
}
