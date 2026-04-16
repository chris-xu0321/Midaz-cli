// Package cli defines the root cobra command and registers all subcommands.
package cli

import (
	"errors"
	"io"
	"os"
	"sync"

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
		Use:   "midaz",
		Short: "Midaz market intelligence CLI",
		Long: `Midaz CLI — authenticate, manage your desk, and query the market intelligence graph.

INSTALL:
    curl -fsSL https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.sh | sh
    (Windows: irm .../install.ps1 | iex)

    Or via npm: npm install -g @midaz/cli

    Full setup: https://github.com/SparkssL/Midaz-cli#installation

Run 'midaz auth login' to sign in, then 'midaz onboard' to set up your radar.`,
		Version: build.Version,
	}
	rootCmd.SilenceErrors = true
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		cmd.SilenceUsage = true
	}

	rootCmd.PersistentFlags().String("format", "json", "Output format: json or pretty")
	rootCmd.PersistentFlags().Bool("raw", false, "Bypass envelope — write raw API response to stdout")
	rootCmd.PersistentFlags().String("api-url", "", "Override API base URL")
	rootCmd.PersistentFlags().String("profile", "", "Auth profile to use (default: current)")

	f := &cmdutil.Factory{IOStreams: ios}

	var (
		cfgOnce sync.Once
		cfgVal  *config.Config
		cfgErr  error
	)
	f.Config = func() (*config.Config, error) {
		cfgOnce.Do(func() {
			flagAPIURL, _ := rootCmd.PersistentFlags().GetString("api-url")
			cfgVal, cfgErr = config.Load(flagAPIURL, "")
		})
		return cfgVal, cfgErr
	}

	var (
		authOnce sync.Once
		authVal  *auth.Creds
		authErr  error
	)
	f.Auth = func() (*auth.Creds, error) {
		authOnce.Do(func() {
			profile, _ := rootCmd.PersistentFlags().GetString("profile")
			authVal, authErr = auth.Current(profile)
		})
		return authVal, authErr
	}

	var (
		clientOnce sync.Once
		clientVal  *client.Client
		clientErr  error
	)
	f.Client = func() (*client.Client, error) {
		clientOnce.Do(func() {
			cfg, err := f.Config()
			if err != nil {
				clientErr = err
				return
			}
			c := client.New(cfg.APIURL)
			if creds, err := f.Auth(); err == nil && creds != nil && creds.APIKey != "" {
				c = c.WithToken(creds.APIKey)
			}
			clientVal = c
		})
		return clientVal, clientErr
	}

	schema.LoadSchemaData = func() []schema.CommandInfo {
		data := make([]schema.CommandInfo, len(registry.Commands))
		for i, def := range registry.Commands {
			argNames := make([]string, len(def.Args))
			for j, a := range def.Args {
				argNames[j] = a.Name
			}
			flagNames := make([]string, len(def.Flags))
			for j, fl := range def.Flags {
				flagNames[j] = "--" + fl.Name
			}
			data[i] = schema.CommandInfo{
				Name:        def.Name,
				Description: def.Description,
				Args:        argNames,
				Flags:       flagNames,
				Endpoints:   def.Endpoints,
			}
		}
		return data
	}

	for _, def := range registry.Commands {
		rootCmd.AddCommand(def.NewCmd(f))
	}

	if err := rootCmd.Execute(); err != nil {
		return handleRootError(ios.ErrOut, err)
	}
	return 0
}

func handleRootError(errOut io.Writer, err error) int {
	var exitErr *output.ExitError
	if errors.As(err, &exitErr) {
		output.WriteErrorEnvelope(errOut, exitErr)
		return exitErr.Code
	}
	wrapped := &output.ExitError{
		Code:   output.ExitInternal,
		Detail: &output.ErrDetail{Code: "internal", Message: err.Error()},
	}
	output.WriteErrorEnvelope(errOut, wrapped)
	return output.ExitInternal
}
