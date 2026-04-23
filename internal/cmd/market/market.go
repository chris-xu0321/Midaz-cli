package market

import (
	"encoding/json"

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
			base := ""
			if cfg, err := f.Config(); err == nil && cfg != nil {
				base = cfg.FrontendURL
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/market",
				Normalize: normalizeMarket(base),
			})
		},
	}
}

func normalizeMarket(frontendBase string) cmdutil.NormalizeFn {
	return func(body []byte) (any, map[string]any, error) {
		rawMap, err := cmdutil.ParseMap(body)
		if err != nil {
			return nil, nil, err
		}

		viewURL := cmdutil.ExtractViewURL(rawMap)
		driverCount := cmdutil.CountArray(rawMap["drivers"])

		if frontendBase != "" {
			if enriched, ok := injectDriverURLs(rawMap["drivers"], frontendBase); ok {
				rawMap["drivers"] = enriched
			}
		}

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
}

// injectDriverURLs walks the drivers[] array, adding `view_url` to each item
// that doesn't already have one. Returns the re-marshalled array and ok=true
// on success; ok=false leaves the caller with the original raw message.
func injectDriverURLs(raw json.RawMessage, frontendBase string) (json.RawMessage, bool) {
	if len(raw) == 0 {
		return raw, false
	}
	var arr []map[string]json.RawMessage
	if err := json.Unmarshal(raw, &arr); err != nil {
		return raw, false
	}
	for _, item := range arr {
		if _, has := item["view_url"]; has {
			continue
		}
		id := cmdutil.UnmarshalString(item["id"])
		if id == "" {
			continue
		}
		u, err := json.Marshal(frontendBase + "/market-read?driver=" + id)
		if err != nil {
			continue
		}
		item["view_url"] = u
	}
	out, err := json.Marshal(arr)
	if err != nil {
		return raw, false
	}
	return out, true
}
