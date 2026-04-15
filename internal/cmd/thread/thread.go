package thread

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdThread is a deprecated alias of `midaz thesis`.
// Kept for one release to avoid breaking existing agent scripts.
func NewCmdThread(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:        "thread <id>",
		Short:      "Deprecated alias of `thesis`",
		Deprecated: "use `midaz thesis`",
		Hidden:     true,
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
			fmt.Fprintln(opts.ErrOut, "note: `thread` is deprecated — use `midaz thesis` instead.")
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/theses/" + url.PathEscape(args[0]),
				Normalize: threadNormalize,
			})
		},
	}
}

type threadMeta struct {
	ViewURL            string            `json:"view_url"`
	TopicURL           string            `json:"topic_url"`
	Claims             []json.RawMessage `json:"claims"`
	MarketLinks        []json.RawMessage `json:"market_links"`
	SupportingCount    int               `json:"supporting_count"`
	ContradictingCount int               `json:"contradicting_count"`
}

func threadNormalize(body []byte) (interface{}, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}
	var tm threadMeta
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
