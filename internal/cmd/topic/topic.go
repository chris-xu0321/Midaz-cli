package topic

import (
	"encoding/json"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdTopic(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "topic <id>",
		Short: "Topic detail + threads",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"Missing required argument: id",
					"usage: seer-q topic <id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/topics/" + url.PathEscape(args[0]),
				Normalize: normalizeTopic,
			})
		},
	}
}

// topicMeta extracts only the fields needed for meta computation.
type topicMeta struct {
	ViewURL string            `json:"view_url"`
	Threads []json.RawMessage `json:"threads"`
}

func normalizeTopic(body []byte) (interface{}, map[string]any, error) {
	// Dual-parse: raw map for data, typed struct for meta
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}

	var tm topicMeta
	if err := json.Unmarshal(body, &tm); err != nil {
		return nil, nil, err
	}

	// Remove view_url from data (moved to meta)
	delete(rawMap, "view_url")

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{
		"thread_count": len(tm.Threads),
	}
	if tm.ViewURL != "" {
		meta["view_url"] = tm.ViewURL
	}

	return data, meta, nil
}

