package usage

import (
	"net/url"

	"github.com/chris-xu0321/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdUsage(f *cmdutil.Factory) *cobra.Command {
	var since string

	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Token usage and cost summary",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			params := url.Values{}
			if since != "" {
				params.Set("since", since)
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/usage",
				Params:    params,
				Normalize: normalizeUsage,
			})
		},
	}

	cmd.Flags().StringVar(&since, "since", "24h", "Time period (e.g., 24h, 7d)")
	return cmd
}

func normalizeUsage(body []byte) (interface{}, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}

	viewURL := cmdutil.ExtractViewURL(rawMap)

	// Extract meta fields (keep since and total_calls in data too)
	sinceVal := cmdutil.UnmarshalString(rawMap["since"])
	totalCalls := cmdutil.UnmarshalInt(rawMap["total_calls"])

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{
		"total_calls": totalCalls,
	}
	if viewURL != "" {
		meta["view_url"] = viewURL
	}
	if sinceVal != "" {
		meta["since"] = sinceVal
	}

	return data, meta, nil
}
