package auth

import (
	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdLogout(f *cmdutil.Factory) *cobra.Command {
	var revoke bool
	cmd := &cobra.Command{
		Use:   "logout",
		Short: "Clear locally stored credentials",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			profile := opts.Profile
			if profile == "" {
				profile = auth.DefaultProfile
			}

			if err := auth.Clear(profile); err != nil {
				return output.ErrWithHint(output.ExitInternal, "internal", err.Error(), "")
			}
			data := map[string]any{
				"profile": profile,
				"revoked": false,
			}
			meta := map[string]any{"message": "Signed out."}
			if revoke {
				meta["hint"] = "to revoke the key server-side, run: midaz auth keys list; then: midaz auth keys revoke <id> --yes"
			}
			return output.WriteSuccess(opts.Out, data, meta, opts.Format)
		},
	}
	cmd.Flags().BoolVar(&revoke, "revoke", false, "Show server-side revoke instructions in the output")
	return cmd
}
