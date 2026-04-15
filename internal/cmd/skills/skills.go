// Package skills provides the `midaz skills` command group for managing
// skill directories on the local machine.
package skills

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSkills(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "skills",
		Short: "Manage embedded agent skills",
		Long: `Manage the skill bundles that ship with the midaz binary.

The installer (curl/PowerShell/npm) only places the binary on your machine.
Use "midaz skills install" to unpack skills into your agent directories
(Claude Code, Codex, or a custom path).`,
	}

	cmd.AddCommand(newCmdInstall(f))
	return cmd
}
