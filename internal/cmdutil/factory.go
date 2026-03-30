// Package cmdutil provides shared infrastructure for seer-q commands:
// Factory (dependency injection), IOStreams, RunOpts, and API command helpers.
package cmdutil

import (
	"context"
	"io"

	"github.com/chris-xu0321/Midaz-cli/internal/client"
	"github.com/chris-xu0321/Midaz-cli/internal/config"
	"github.com/spf13/cobra"
)

// IOStreams holds the standard I/O writers for the CLI.
type IOStreams struct {
	Out    io.Writer // stdout — success envelopes
	ErrOut io.Writer // stderr — error envelopes
}

// Factory provides shared dependencies to commands via lazy initialization.
// Local commands (version, schema) never call Config/Client, so a bad config
// file cannot break them.
type Factory struct {
	IOStreams *IOStreams
	Config   func() (*config.Config, error)  // lazy — only called by commands that need it
	Client   func() (*client.Client, error)  // lazy — only called by API commands
}

// RunOpts holds explicit runtime options for command execution.
// Format defaults to "json" via cobra flag — no config dependency.
type RunOpts struct {
	Ctx    context.Context
	Format string    // "json" or "pretty"
	Raw    bool      // bypass envelope, write raw API response
	Out    io.Writer // stdout
	ErrOut io.Writer // stderr
}

// ResolveRunOpts reads --format and --raw flags from the cobra command.
// If --format was not explicitly passed, falls back to config's format value.
// If config fails (local commands with bad config file), defaults to "json".
func ResolveRunOpts(cmd *cobra.Command, f *Factory) *RunOpts {
	format, _ := cmd.Flags().GetString("format")

	// If --format was not explicitly passed, try config
	if !cmd.Flags().Changed("format") {
		if cfg, err := f.Config(); err == nil && cfg.Format != "" {
			format = cfg.Format
		}
		// If config fails or format is empty, keep cobra default "json"
	}

	if format == "" {
		format = "json"
	}

	raw, _ := cmd.Flags().GetBool("raw")

	return &RunOpts{
		Ctx:    cmd.Context(),
		Format: format,
		Raw:    raw,
		Out:    f.IOStreams.Out,
		ErrOut: f.IOStreams.ErrOut,
	}
}
