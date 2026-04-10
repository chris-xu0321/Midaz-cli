// Package logout implements the `seer-q logout` command.
package logout

import (
	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdLogout creates the logout command.
func NewCmdLogout(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "logout",
		Short: "Clear stored Seer credentials",
		Long:  "Removes the locally stored API key and workspace info.",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)

			if err := auth.Clear(); err != nil {
				return output.ErrConfig("failed to clear credentials: %s", err)
			}

			result := map[string]any{
				"ok":      true,
				"message": "Logged out. Credentials removed.",
			}
			return output.WriteSuccess(opts.Out, result, nil, opts.Format)
		},
	}
}
