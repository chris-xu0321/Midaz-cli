// Package assets hosts `midaz assets …` (list, get, timeline).
package assets

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAssets builds the assets command tree.
func NewCmdAssets(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "List and inspect assets (bias, contributions, timeline)",
	}
	cmd.AddCommand(newCmdList(f))
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdTimeline(f))
	cmd.AddCommand(newCmdOptions(f))
	return cmd
}
