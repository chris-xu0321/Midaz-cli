package topics

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdTopics(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "topics",
		Short: "List all topics with thread counts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/topics",
				Normalize: cmdutil.NormalizeBareArray,
			})
		},
	}
}
