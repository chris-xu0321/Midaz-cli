package sources

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSources(f *cmdutil.Factory) *cobra.Command {
	var decision string
	var tier int

	cmd := &cobra.Command{
		Use:   "sources",
		Short: "List ingested sources",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			params := url.Values{}
			if decision != "" {
				params.Set("decision", decision)
			}
			if tier > 0 {
				params.Set("tier", fmt.Sprintf("%d", tier))
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/sources",
				Params:    params,
				Normalize: normalizeSources,
			})
		},
	}

	cmd.Flags().StringVar(&decision, "decision", "", "Filter by gate decision")
	cmd.Flags().IntVar(&tier, "tier", 0, "Filter by source tier")
	return cmd
}

func normalizeSources(body []byte) (interface{}, map[string]any, error) {
	m, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}

	viewURL := cmdutil.ExtractViewURL(m)

	// Unwrap items
	itemsRaw, ok := m["items"]
	if !ok {
		return nil, nil, fmt.Errorf("expected 'items' key in sources response")
	}

	var items interface{}
	if err := json.Unmarshal(itemsRaw, &items); err != nil {
		return nil, nil, err
	}

	meta := map[string]any{"count": cmdutil.CountArray(itemsRaw)}
	if viewURL != "" {
		meta["view_url"] = viewURL
	}

	return items, meta, nil
}
