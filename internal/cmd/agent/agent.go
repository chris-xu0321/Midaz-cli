package agent

import (
	"slices"
	"strings"

	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/SparkssL/seer-cli/internal/output"
	"github.com/spf13/cobra"
)

var supportedPlatforms = []string{"claude"}

func NewCmdAgent(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "Agent resource management",
	}
	cmd.AddCommand(newCmdInstall(f))
	cmd.AddCommand(newCmdUninstall(f))
	cmd.AddCommand(newCmdDoctor(f))
	return cmd
}

// exactPlatform validates that exactly one positional arg is provided
// and it matches a supported agent platform.
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
		if !slices.Contains(supportedPlatforms, args[0]) {
			return output.ErrValidation("unknown platform: %s (supported: %s)",
				args[0], strings.Join(supportedPlatforms, ", "))
		}
		return nil
	}
}
