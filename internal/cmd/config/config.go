package config

import (
	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdConfig(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configuration management",
	}
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdSet(f))
	cmd.AddCommand(newCmdList(f))
	cmd.AddCommand(newCmdPath(f))
	return cmd
}
