package claims

import (
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdClaims(f *cmdutil.Factory) *cobra.Command {
	var threadID, sourceID, status, mode string

	cmd := &cobra.Command{
		Use:   "claims",
		Short: "List claims",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			params := url.Values{}
			if threadID != "" {
				params.Set("thread_id", threadID)
			}
			if sourceID != "" {
				params.Set("source_id", sourceID)
			}
			if status != "" {
				params.Set("status", status)
			}
			if mode != "" {
				params.Set("claim_mode", mode)
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/claims",
				Params:    params,
				Normalize: cmdutil.NormalizeBareArray,
			})
		},
	}

	cmd.Flags().StringVar(&threadID, "thread", "", "Filter by thread ID")
	cmd.Flags().StringVar(&sourceID, "source", "", "Filter by source ID")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	cmd.Flags().StringVar(&mode, "mode", "", "Filter by claim mode")
	return cmd
}
