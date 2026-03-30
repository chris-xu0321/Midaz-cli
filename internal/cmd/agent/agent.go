package agent

import (
	"fmt"
	"slices"

	"github.com/chris-xu0321/Midaz-cli/internal/cmdutil"
	"github.com/chris-xu0321/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// supportedBridgePlatforms lists platforms with embedded compatibility bridge assets.
// Other platforms should use the agent ecosystem's skill installer.
var supportedBridgePlatforms = []string{"claude"}

func NewCmdAgent(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Agent resource management [deprecated — prefer skill installer]",
	}
	cmd.AddCommand(newCmdInstall(f))
	cmd.AddCommand(newCmdUninstall(f))
	cmd.AddCommand(newCmdDoctor(f))
	return cmd
}

// exactPlatform validates that exactly one positional arg is provided
// and it matches a supported bridge platform.
func exactPlatform() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		name := "install"
		if cmd != nil {
			name = cmd.Name()
		}
		if len(args) < 1 {
			return output.ErrWithHint(output.ExitValidation, "validation",
				"Missing platform argument",
				"usage: seer-q agent "+name+" claude")
		}
		if len(args) > 1 {
			return output.ErrWithHint(output.ExitValidation, "validation",
				"Too many arguments",
				"usage: seer-q agent "+name+" claude")
		}
		if !slices.Contains(supportedBridgePlatforms, args[0]) {
			return output.ErrWithHint(output.ExitValidation, "validation",
				fmt.Sprintf("platform %q is not supported by the compatibility bridge", args[0]),
				"install skills through your agent platform's skill installer instead.\n  See: https://github.com/SparkssL/seer-skills")
		}
		return nil
	}
}
