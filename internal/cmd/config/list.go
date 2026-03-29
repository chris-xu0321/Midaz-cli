package config

import (
	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/SparkssL/seer-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all config values",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)

			cfg, err := f.Config()
			if err != nil {
				return err
			}

			data := map[string]string{
				"api_url":      cfg.APIURL,
				"frontend_url": cfg.FrontendURL,
				"format":       cfg.Format,
			}
			return output.WriteSuccess(opts.Out, data, nil, opts.Format)
		},
	}
}
