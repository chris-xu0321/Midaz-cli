package config

import (
	"slices"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	cfgpkg "github.com/SparkssL/Midaz-cli/internal/config"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdSet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "set <key> <value>",
		Short: "Set a config value",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"Missing required arguments: key and value",
					"usage: seer-q config set <key> <value>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			key, value := args[0], args[1]

			if !slices.Contains(cfgpkg.ValidKeys, key) {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"Unknown config key: "+key,
					"valid keys: api_url, frontend_url, format")
			}

			if err := cfgpkg.SetKey(key, value); err != nil {
				return output.ErrConfig("failed to set config: %s", err)
			}

			data := map[string]string{"key": key, "value": value}
			return output.WriteSuccess(opts.Out, data, nil, opts.Format)
		},
	}
}
