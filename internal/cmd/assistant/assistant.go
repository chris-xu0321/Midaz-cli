// Package assistant hosts `midaz assistant …` — desk assistant
// surface area. Today only `events` is exposed; the conversational
// `/chat` endpoint is interactive and stays out of the CLI for now.
package assistant

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdAssistant builds the `assistant` parent command.
func NewCmdAssistant(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "assistant",
		Short: "Desk assistant (events feed)",
	}
	cmd.AddCommand(newCmdEvents(f))
	return cmd
}
