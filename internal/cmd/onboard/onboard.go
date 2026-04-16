// Package onboard hosts `midaz onboard` (status, generate, complete).
package onboard

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdOnboard builds the onboard command tree.
func NewCmdOnboard(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "onboard",
		Short: "Complete desk onboarding (radar + playbook)",
		Long: `Onboarding sets your radar (watchlist) and playbook (trading rules)
for the first time. Two paths:

  midaz onboard generate --mode guided --from-file input.json --yes
      Lets the server draft radar + playbook from a structured prompt.

  midaz onboard complete --radar radar.md --playbook playbook.md --yes
      Commits caller-supplied radar and playbook content directly.`,
	}
	cmd.AddCommand(newCmdStatus(f))
	cmd.AddCommand(newCmdGenerate(f))
	cmd.AddCommand(newCmdComplete(f))
	return cmd
}
