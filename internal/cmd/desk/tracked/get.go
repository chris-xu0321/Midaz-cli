package tracked

import (
	"encoding/json"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Show the current tracked asset list and the valid universe",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/desk/settings",
				Normalize: normalizeTrackedGet,
			})
		},
	}
}

// normalizeTrackedGet projects /api/desk/settings down to the tracked-
// asset slice + asset universe (so callers don't see the whole settings
// blob when they only asked for tracked assets).
func normalizeTrackedGet(body []byte) (any, map[string]any, error) {
	var parsed struct {
		TrackedAssetIDs []string `json:"tracked_asset_ids"`
		AssetUniverse   []struct {
			AssetID string   `json:"asset_id"`
			Name    string   `json:"name"`
			Aliases []string `json:"aliases"`
		} `json:"asset_universe"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, nil, output.Errorf(output.ExitInternal, "internal",
			"failed to parse /api/desk/settings: %s", err)
	}
	data := map[string]any{
		"tracked_asset_ids": parsed.TrackedAssetIDs,
		"asset_universe":    parsed.AssetUniverse,
	}
	meta := map[string]any{
		"tracked_count": len(parsed.TrackedAssetIDs),
		"universe_size": len(parsed.AssetUniverse),
	}
	return data, meta, nil
}
