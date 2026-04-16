package snapshot

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdSnapshot(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "snapshot",
		Short: "Global regime snapshot",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/global",
				Normalize: normalizeSnapshot,
			})
		},
	}
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
