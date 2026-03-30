package search

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/chris-xu0321/Midaz-cli/internal/cmdutil"
	"github.com/chris-xu0321/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Fuzzy search across topics, threads, assets",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 || args[0] == "" {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"Missing required argument: query",
					"usage: seer-q search <query>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			params := url.Values{"q": {args[0]}}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/search",
				Params:    params,
				Normalize: normalizeSearch,
			})
		},
	}
}

func normalizeSearch(body []byte) (interface{}, map[string]any, error) {
	var resp struct {
		Results []json.RawMessage `json:"results"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, nil, fmt.Errorf("expected search response: %w", err)
	}
	// Re-parse to get proper interface{} for marshaling
	var full map[string]interface{}
	json.Unmarshal(body, &full)
	results := full["results"]

	return results, map[string]any{"count": len(resp.Results)}, nil
}
