// Package thesis hosts `midaz thesis <id>` — thesis detail + claims + links.
package thesis

import (
	"encoding/json"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdThesis returns the `thesis` cobra command.
func NewCmdThesis(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "thesis <id>",
		Short: "Thesis detail + claims + market links",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"Missing required argument: id",
					"usage: midaz thesis <id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/theses/" + url.PathEscape(args[0]),
				Normalize: normalize,
			})
		},
	}
}

type thesisMeta struct {
	ViewURL            string            `json:"view_url"`
	TopicURL           string            `json:"topic_url"`
	Claims             []json.RawMessage `json:"claims"`
	MarketLinks        []json.RawMessage `json:"market_links"`
	SupportingCount    int               `json:"supporting_count"`
	ContradictingCount int               `json:"contradicting_count"`
}

func normalize(body []byte) (interface{}, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}
	var tm thesisMeta
	if err := json.Unmarshal(body, &tm); err != nil {
		return nil, nil, err
	}
	delete(rawMap, "view_url")
	delete(rawMap, "topic_url")
	delete(rawMap, "has_market_link")
	delete(rawMap, "market_link_count")
	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}
	meta := map[string]any{
		"claim_count":         len(tm.Claims),
		"supporting_count":    tm.SupportingCount,
		"contradicting_count": tm.ContradictingCount,
		"market_link_count":   len(tm.MarketLinks),
	}
	if tm.ViewURL != "" {
		meta["view_url"] = tm.ViewURL
	}
	if tm.TopicURL != "" {
		meta["topic_url"] = tm.TopicURL
	}
	return data, meta, nil
}
