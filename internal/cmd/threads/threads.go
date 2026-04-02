package threads

import (
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdThreads(f *cmdutil.Factory) *cobra.Command {
	var topicID, status string

	cmd := &cobra.Command{
		Use:   "threads",
		Short: "List threads",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			params := url.Values{}
			if topicID != "" {
				params.Set("topic_id", topicID)
			}
			if status != "" {
				params.Set("status", status)
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/threads",
				Params:    params,
				Normalize: cmdutil.NormalizeBareArray,
			})
		},
	}

	cmd.Flags().StringVar(&topicID, "topic", "", "Filter by topic ID")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status (active/weakening/divided/resolved)")
	return cmd
}
