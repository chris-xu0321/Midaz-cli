// Package workspace hosts `midaz workspace …` commands.
package workspace

import (
	"github.com/SparkssL/Midaz-cli/internal/cmd/workspace/playbook"
	"github.com/SparkssL/Midaz-cli/internal/cmd/workspace/radar"
	"github.com/SparkssL/Midaz-cli/internal/cmd/workspace/telegram"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdWorkspace builds the `workspace` parent command.
func NewCmdWorkspace(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Manage your workspace (radar, playbook, sharing, Telegram)",
		Long: `Inspect and configure your Midaz workspace.

All write commands require --yes; they mirror the frontend toggles at
/workspace/settings.`,
	}
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdSettings(f))
	cmd.AddCommand(newCmdView(f))
	cmd.AddCommand(newCmdShare(f))
	cmd.AddCommand(radar.NewCmdRadar(f))
	cmd.AddCommand(playbook.NewCmdPlaybook(f))
	cmd.AddCommand(telegram.NewCmdTelegram(f))
	return cmd
}
