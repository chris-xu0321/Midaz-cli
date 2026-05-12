package assistant

import (
	"net/url"
	"strconv"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdEvents(f *cmdutil.Factory) *cobra.Command {
	var (
		after string
		limit int
	)
	cmd := &cobra.Command{
		Use:   "events",
		Short: "Recent assistant alert/inbox events for the current desk",
		Long: `Owner-only, subscription-gated. Surfaces recent Gate A
notifications and inbox events so chat assistants can include fresh
context without storing transcripts.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			params := url.Values{}
			if after != "" {
				params.Set("after", after)
			}
			if limit > 0 {
				params.Set("limit", strconv.Itoa(limit))
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/assistant/events",
				Params:    params,
				Normalize: normalizeEvents,
			})
		},
	}
	cmd.Flags().StringVar(&after, "after", "", "Only return events after this ISO-8601 timestamp")
	cmd.Flags().IntVar(&limit, "limit", 0, "Max events to return (server default 20, max 50)")
	return cmd
}

func normalizeEvents(body []byte) (any, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}
	count := cmdutil.CountArray(rawMap["events"])
	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}
	return data, map[string]any{"event_count": count}, nil
}
