// Package cmdutil provides shared infrastructure for midaz commands:
// Factory (dependency injection), IOStreams, RunOpts, and API command helpers.
package cmdutil

import (
	"context"
	"io"

	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/client"
	"github.com/SparkssL/Midaz-cli/internal/config"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// IOStreams holds the standard I/O writers for the CLI.
type IOStreams struct {
	Out    io.Writer // stdout — success envelopes
	ErrOut io.Writer // stderr — error envelopes
}

// Factory provides shared dependencies to commands via lazy initialization,
// so commands that don't need network/config (e.g. version, schema) aren't
// broken by a missing or malformed config file.
type Factory struct {
	IOStreams *IOStreams
	Config    func() (*config.Config, error)
	Client    func() (*client.Client, error)
	Auth      func() (*auth.Creds, error) // may return nil creds without error when unauthenticated
}

// RunOpts holds explicit runtime options for command execution.
type RunOpts struct {
	Ctx     context.Context
	Format  string
	Raw     bool
	Out     io.Writer
	ErrOut  io.Writer
	Profile string
}

// ResolveRunOpts reads --format, --raw, and --profile flags from the cobra command.
func ResolveRunOpts(cmd *cobra.Command, f *Factory) *RunOpts {
	format, _ := cmd.Flags().GetString("format")

	if !cmd.Flags().Changed("format") {
		if cfg, err := f.Config(); err == nil && cfg.Format != "" {
			format = cfg.Format
		}
	}
	if format == "" {
		format = "json"
	}

	raw, _ := cmd.Flags().GetBool("raw")
	profile, _ := cmd.Flags().GetString("profile")

	return &RunOpts{
		Ctx:     cmd.Context(),
		Format:  format,
		Raw:     raw,
		Out:     f.IOStreams.Out,
		ErrOut:  f.IOStreams.ErrOut,
		Profile: profile,
	}
}

// RequireAuth loads credentials or returns an ErrAuth. Use for commands that
// must be authenticated.
func RequireAuth(f *Factory) (*auth.Creds, error) {
	c, err := f.Auth()
	if err != nil {
		return nil, output.ErrAuth(err.Error(), "")
	}
	if c == nil {
		return nil, output.ErrAuth("not logged in", "")
	}
	return c, nil
}

// AuthedClient returns a client with the Authorization header pre-set. Returns
// ErrAuth if the user isn't logged in.
func AuthedClient(f *Factory) (*client.Client, *auth.Creds, error) {
	creds, err := RequireAuth(f)
	if err != nil {
		return nil, nil, err
	}
	base, err := f.Client()
	if err != nil {
		return nil, nil, err
	}
	return base.WithToken(creds.APIKey), creds, nil
}
