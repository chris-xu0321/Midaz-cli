package market

import (
	"encoding/json"

	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdMarket(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "market",
		Short: "Global regime + all topics with thread counts",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/market",
				Normalize: normalizeMarket,
			})
		},
	}
}

func normalizeMarket(body []byte) (interface{}, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}

	viewURL := cmdutil.ExtractViewURL(rawMap)

	topicCount := 0
	if raw, ok := rawMap["topics"]; ok {
		var arr []json.RawMessage
		json.Unmarshal(raw, &arr)
		topicCount = len(arr)
	}

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{"topic_count": topicCount}
	if viewURL != "" {
		meta["view_url"] = viewURL
	}

	return data, meta, nil
}
