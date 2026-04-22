package market

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdMarket(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "market",
		Short: "Global regime + drivers + thesis memberships (composite)",
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

func normalizeMarket(body []byte) (any, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}

	viewURL := cmdutil.ExtractViewURL(rawMap)
	driverCount := cmdutil.CountArray(rawMap["drivers"])

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{"driver_count": driverCount}
	if viewURL != "" {
		meta["view_url"] = viewURL
	}

	return data, meta, nil
}
