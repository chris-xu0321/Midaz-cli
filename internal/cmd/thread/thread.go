package thread

import (
	"encoding/json"
	"net/url"

	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/SparkssL/seer-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdThread(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "thread <id>",
		Short: "Thread detail + claims + market links",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"Missing required argument: id",
					"usage: seer-q thread <id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/threads/" + url.PathEscape(args[0]),
				Normalize: normalizeThread,
			})
		},
	}
}

// threadMeta extracts only the fields needed for meta computation.
type threadMeta struct {
	ViewURL            string            `json:"view_url"`
	TopicURL           string            `json:"topic_url"`
	Claims             []json.RawMessage `json:"claims"`
	MarketLinks        []json.RawMessage `json:"market_links"`
	SupportingCount    int               `json:"supporting_count"`
	ContradictingCount int               `json:"contradicting_count"`
}

func normalizeThread(body []byte) (interface{}, map[string]any, error) {
	// Dual-parse: raw map for data, typed struct for meta
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}

	var tm threadMeta
	if err := json.Unmarshal(body, &tm); err != nil {
		return nil, nil, err
	}

	// Remove meta-only and computed fields from data
	delete(rawMap, "view_url")
	delete(rawMap, "topic_url")
	delete(rawMap, "has_market_link")
	delete(rawMap, "market_link_count")
	// Keep supporting_count and contradicting_count in data (duplicated in meta)

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{
		"claim_count":        len(tm.Claims),
		"supporting_count":   tm.SupportingCount,
		"contradicting_count": tm.ContradictingCount,
		"market_link_count":  len(tm.MarketLinks),
	}
	if tm.ViewURL != "" {
		meta["view_url"] = tm.ViewURL
	}
	if tm.TopicURL != "" {
		meta["topic_url"] = tm.TopicURL
	}

	return data, meta, nil
}
