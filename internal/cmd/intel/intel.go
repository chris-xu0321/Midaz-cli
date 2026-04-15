// Package intel hosts `midaz intel …` (list, push, rm).
package intel

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdIntel builds the intel command tree.
func NewCmdIntel(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intel",
		Short: "Push, list, and delete private intel (notes / sources)",
	}
	cmd.AddCommand(newCmdList(f))
	cmd.AddCommand(newCmdPush(f))
	cmd.AddCommand(newCmdRm(f))
	return cmd
}
