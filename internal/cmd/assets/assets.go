// Package assets hosts `midaz assets …` (list, get, thesis).
package assets

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAssets builds the assets command tree.
func NewCmdAssets(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assets",
		Short: "List and inspect assets with thesis links",
	}
	cmd.AddCommand(newCmdList(f))
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdThesis(f))
	return cmd
}
