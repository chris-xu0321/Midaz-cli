// Package theses hosts `midaz theses` (list).
package theses

import (
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdTheses lists theses (market arguments), with optional filters.
func NewCmdTheses(f *cmdutil.Factory) *cobra.Command {
	var topicID, status string
	cmd := &cobra.Command{
		Use:   "theses",
		Short: "List theses (market arguments)",
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
				Path:      "/api/theses",
				Params:    params,
				Normalize: cmdutil.NormalizeBareArray,
			})
		},
	}
	cmd.Flags().StringVar(&topicID, "topic", "", "Filter by topic ID")
	cmd.Flags().StringVar(&status, "status", "", "Filter by status (active/weakening/divided/resolved)")
	return cmd
}
