// Package tracked hosts `midaz desk tracked-assets` (get/set/add/remove).
//
// The tracked-asset list is the L4 asset scope: every tracked asset
// without an open position becomes a Monitoring card on the desk read.
// Radar edits no longer change this scope — use these verbs to control
// what L4 monitors.
//
// Backed by PATCH /api/desk/tracked-assets, which enqueues an L4
// rebuild with reason=asset_scope_edit. The settings response advertises
// the valid universe via `asset_universe[]` so we can echo a count.
package tracked

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/SparkssL/Midaz-cli/internal/client"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdTracked builds the `desk tracked-assets` subcommand tree.
func NewCmdTracked(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "tracked-assets",
		Aliases: []string{"tracked"},
		Short:   "Manage the L4 asset scope (tracked assets that become Monitoring cards)",
	}
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdSet(f))
	cmd.AddCommand(newCmdAdd(f))
	cmd.AddCommand(newCmdRemove(f))
	return cmd
}

// fetchTracked returns the current tracked_asset_ids from GET /api/desk/settings.
func fetchTracked(ctx context.Context, c *client.Client) ([]string, error) {
	resp, err := c.Get(ctx, "/api/desk/settings", nil)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		TrackedAssetIDs []string `json:"tracked_asset_ids"`
	}
	if err := json.Unmarshal(resp.Body, &parsed); err != nil {
		return nil, output.Errorf(output.ExitInternal, "internal",
			"failed to parse /api/desk/settings response: %s", err)
	}
	return parsed.TrackedAssetIDs, nil
}

// pushTracked writes the full tracked-assets list back via PATCH.
func pushTracked(ctx context.Context, c *client.Client, ids []string) (map[string]any, error) {
	resp, err := c.Patch(ctx, "/api/desk/tracked-assets", map[string]any{"tracked_asset_ids": ids})
	if err != nil {
		return nil, err
	}
	var out map[string]any
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, &out); err != nil {
			return nil, output.Errorf(output.ExitInternal, "internal",
				"failed to parse /api/desk/tracked-assets response: %s", err)
		}
	}
	return out, nil
}

// parseAssetList splits a comma-separated string into deduped uppercase
// asset IDs. Whitespace is trimmed. Empty entries are dropped.
func parseAssetList(raw string) []string {
	if raw == "" {
		return nil
	}
	seen := map[string]struct{}{}
	out := []string{}
	for _, part := range strings.Split(raw, ",") {
		s := strings.ToUpper(strings.TrimSpace(part))
		if s == "" {
			continue
		}
		if _, ok := seen[s]; ok {
			continue
		}
		seen[s] = struct{}{}
		out = append(out, s)
	}
	return out
}

// mergeUnique merges b into a (uppercased), preserving order of a then
// new entries from b. Duplicates are dropped.
func mergeUnique(a, b []string) []string {
	seen := map[string]struct{}{}
	out := make([]string, 0, len(a)+len(b))
	for _, x := range a {
		u := strings.ToUpper(strings.TrimSpace(x))
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		out = append(out, u)
	}
	for _, x := range b {
		u := strings.ToUpper(strings.TrimSpace(x))
		if u == "" {
			continue
		}
		if _, ok := seen[u]; ok {
			continue
		}
		seen[u] = struct{}{}
		out = append(out, u)
	}
	return out
}

// removeFrom returns a copy of `from` with any element in `drop` (case-
// insensitive) removed.
func removeFrom(from, drop []string) []string {
	skip := map[string]struct{}{}
	for _, x := range drop {
		skip[strings.ToUpper(strings.TrimSpace(x))] = struct{}{}
	}
	out := make([]string, 0, len(from))
	for _, x := range from {
		if _, hit := skip[strings.ToUpper(strings.TrimSpace(x))]; hit {
			continue
		}
		out = append(out, x)
	}
	return out
}

