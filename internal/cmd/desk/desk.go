// Package desk hosts `midaz desk …` commands.
package desk

import (
	"github.com/SparkssL/Midaz-cli/internal/cmd/desk/playbook"
	"github.com/SparkssL/Midaz-cli/internal/cmd/desk/radar"
	"github.com/SparkssL/Midaz-cli/internal/cmd/desk/telegram"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdDesk builds the `desk` parent command.
func NewCmdDesk(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "desk",
		Short: "Manage your desk (radar, playbook, sharing, Telegram)",
		Long: `Inspect and configure your Midaz desk.

All write commands require --yes; they mirror the frontend toggles at
/desk/settings.`,
	}
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdSettings(f))
	cmd.AddCommand(newCmdView(f))
	cmd.AddCommand(newCmdShare(f))
	cmd.AddCommand(newCmdRegenerate(f))
	cmd.AddCommand(newCmdReonboard(f))
	cmd.AddCommand(radar.NewCmdRadar(f))
	cmd.AddCommand(playbook.NewCmdPlaybook(f))
	cmd.AddCommand(telegram.NewCmdTelegram(f))
	return cmd
}
