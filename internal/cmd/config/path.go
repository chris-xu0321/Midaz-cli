package config

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	cfgpkg "github.com/SparkssL/Midaz-cli/internal/config"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdPath(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "path",
		Short: "Show config file path",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			data := map[string]string{"path": cfgpkg.ConfigPath()}
			return output.WriteSuccess(opts.Out, data, nil, opts.Format)
		},
	}
}
