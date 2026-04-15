// Package keys exposes `midaz auth keys` — list / create / revoke PATs.
package keys

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdKeys builds the `auth keys` parent command.
func NewCmdKeys(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keys",
		Short: "Manage personal access tokens (PATs)",
	}
	cmd.AddCommand(newCmdList(f))
	cmd.AddCommand(newCmdCreate(f))
	cmd.AddCommand(newCmdRevoke(f))
	return cmd
}
