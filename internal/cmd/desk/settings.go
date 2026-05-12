package desk

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func newCmdSettings(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "settings",
		Short: "Owner-only settings: radar, playbook, Telegram, tracked-assets",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/desk/settings",
				Normalize: normalizeSettings,
			})
		},
	}
}

// normalizeSettings surfaces the asset-universe size and tracked-asset
// count in meta so callers can sanity-check the L4 scope without
// scrolling the full settings blob.
func normalizeSettings(body []byte) (any, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}
	universeSize := cmdutil.CountArray(rawMap["asset_universe"])
	trackedCount := cmdutil.CountArray(rawMap["tracked_asset_ids"])
	radarItems := cmdutil.CountArray(rawMap["radar_items"])
	openPositions := cmdutil.CountArray(rawMap["positions"])

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{
		"universe_size":  universeSize,
		"tracked_count":  trackedCount,
		"radar_items":    radarItems,
		"open_positions": openPositions,
	}
	return data, meta, nil
}
