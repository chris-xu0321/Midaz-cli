package config

import (
	"slices"

	"github.com/SparkssL/seer-cli/internal/cmdutil"
	cfgpkg "github.com/SparkssL/seer-cli/internal/config"
	"github.com/SparkssL/seer-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get <key>",
		Short: "Get a config value",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"Missing required argument: key",
					"usage: seer-q config get <key>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			key := args[0]

			if !slices.Contains(cfgpkg.ValidKeys, key) {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"Unknown config key: "+key,
					"valid keys: api_url, frontend_url, format")
			}

			cfg, err := f.Config()
			if err != nil {
				return err
			}

			var value string
			switch key {
			case "api_url":
				value = cfg.APIURL
			case "frontend_url":
				value = cfg.FrontendURL
			case "format":
				value = cfg.Format
			}

			flagVal, _ := cmd.Flags().GetString("api-url")
			source := cfgpkg.Source(key, flagVal)

			data := map[string]string{"key": key, "value": value}
			meta := map[string]any{"source": source}
			return output.WriteSuccess(opts.Out, data, meta, opts.Format)
		},
	}
}
