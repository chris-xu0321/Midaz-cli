// Package cli defines the root cobra command and registers all subcommands.
package cli

import (
	"errors"
	"io"
	"os"

	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/build"
	"github.com/SparkssL/Midaz-cli/internal/client"
	"github.com/SparkssL/Midaz-cli/internal/cmd/schema"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/config"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/SparkssL/Midaz-cli/internal/registry"
	"github.com/spf13/cobra"
)

// Execute runs the root command and returns the process exit code.
func Execute() int {
	ios := &cmdutil.IOStreams{
		Out:    os.Stdout,
		ErrOut: os.Stderr,
	}

	rootCmd := &cobra.Command{
		Use:   "seer-q",
		Short: "Seer market intelligence CLI",
		Long: `Seer market intelligence CLI.

INSTALL:
    curl -fsSL https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.sh | sh
    (Windows: irm .../install.ps1 | iex)

    Or via npm: npm install -g @midaz/cli && npx skills add SparkssL/Midaz-cli -y -g

    Full setup: https://github.com/SparkssL/Midaz-cli#installation`,
		Version: build.Version,
	}
	rootCmd.SilenceErrors = true
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
	}

	// Global persistent flags
	rootCmd.PersistentFlags().String("format", "json", "Output format: json or pretty")
	rootCmd.PersistentFlags().Bool("raw", false, "Bypass envelope — write raw API response to stdout")
	rootCmd.PersistentFlags().String("api-url", "", "Override API base URL")

	// Build factory with lazy config and client
	f := &cmdutil.Factory{
		IOStreams: ios,
	}

	// Lazy config — only called by commands that need it
	f.Config = func() (*config.Config, error) {
		flagAPIURL, _ := rootCmd.PersistentFlags().GetString("api-url")
		return config.Load(flagAPIURL, "")
	}

	// Lazy client — only called by API commands
	f.Client = func() (*client.Client, error) {
		cfg, err := f.Config()
		if err != nil {
			return nil, err
		}
		c := client.New(cfg.APIURL)
		// Inject auth token: env > config > credentials
		c.AuthToken = auth.ResolveToken(cfg.APIKey)
		return c, nil
	}

	// Populate schema data from registry (breaks the import cycle)
	schemaData := make([]schema.CommandInfo, len(registry.Commands))
	for i, def := range registry.Commands {
		argNames := make([]string, len(def.Args))
		for j, a := range def.Args {
			argNames[j] = a.Name
		}
		flagNames := make([]string, len(def.Flags))
		for j, fl := range def.Flags {
			flagNames[j] = "--" + fl.Name
		}
		schemaData[i] = schema.CommandInfo{
			Name:        def.Name,
			Description: def.Description,
			Args:        argNames,
			Flags:       flagNames,
			Endpoints:   def.Endpoints,
		}
	}
	schema.SchemaData = schemaData

	// Register all commands from the registry
	for _, def := range registry.Commands {
		rootCmd.AddCommand(def.NewCmd(f))
	}

	if err := rootCmd.Execute(); err != nil {
		return handleRootError(ios.ErrOut, err)
	}
	return 0
}

// handleRootError converts errors to exit codes and writes error envelopes.
// Known *ExitError → use its code. Unknown errors → exit 1 (internal).
func handleRootError(errOut io.Writer, err error) int {
	var exitErr *output.ExitError
	if errors.As(err, &exitErr) {
		output.WriteErrorEnvelope(errOut, exitErr)
		return exitErr.Code
	}

	// Unknown error (cobra internals, unexpected) → exit 1
	wrapped := &output.ExitError{
		Code:   output.ExitInternal,
		Detail: &output.ErrDetail{Code: "internal", Message: err.Error()},
	}
	output.WriteErrorEnvelope(errOut, wrapped)
	return output.ExitInternal
}
