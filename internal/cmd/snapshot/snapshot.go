package snapshot

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/chris-xu0321/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSnapshot(f *cmdutil.Factory) *cobra.Command {
	var history bool
	var limit int

	cmd := &cobra.Command{
		Use:   "snapshot",
		Short: "Global regime snapshot",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)

			if history {
				params := url.Values{}
				if limit > 0 {
					params.Set("limit", fmt.Sprintf("%d", limit))
				}
				return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
					Path:   "/api/global/snapshots",
					Params: params,
					Normalize: func(body []byte) (interface{}, map[string]any, error) {
						var arr []json.RawMessage
						if err := json.Unmarshal(body, &arr); err != nil {
							return nil, nil, err
						}
						var data interface{}
						json.Unmarshal(body, &data)
						effectiveLimit := limit
						if effectiveLimit == 0 {
							effectiveLimit = 10
						}
						return data, map[string]any{
							"count": len(arr),
							"limit": effectiveLimit,
						}, nil
					},
				})
			}

			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/global/snapshot",
				Normalize: normalizeSnapshot,
			})
		},
	}

	cmd.Flags().BoolVar(&history, "history", false, "Show snapshot history")
	cmd.Flags().IntVar(&limit, "limit", 0, "Limit history count (default 10)")
	return cmd
}

func normalizeSnapshot(body []byte) (interface{}, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}

	viewURL := cmdutil.ExtractViewURL(rawMap)

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{}
	if viewURL != "" {
		meta["view_url"] = viewURL
	}

	return data, meta, nil
}
