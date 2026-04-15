// Package auth hosts the `midaz auth` subcommand tree.
package auth

import (
	"github.com/SparkssL/Midaz-cli/internal/cmd/auth/keys"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAuth builds the `auth` parent command.
func NewCmdAuth(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate and manage API credentials",
		Long: `Authenticate with Midaz and manage API credentials.

Login flow:
  midaz auth login              # Opens browser (default)
  midaz auth login --paste      # Paste a PAT created on the website
  midaz auth login --token sk_… # Provide PAT inline (CI / headless)

Other subcommands let you inspect the current session and manage personal
access tokens (PATs) without logging out.`,
	}
	cmd.AddCommand(newCmdLogin(f))
	cmd.AddCommand(newCmdLogout(f))
	cmd.AddCommand(newCmdStatus(f))
	cmd.AddCommand(newCmdWhoami(f))
	cmd.AddCommand(keys.NewCmdKeys(f))
	return cmd
}
