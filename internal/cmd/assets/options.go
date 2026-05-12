package assets

import (
	"net/url"
	"strconv"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdOptions(f *cmdutil.Factory) *cobra.Command {
	var maxDays int
	cmd := &cobra.Command{
		Use:   "options <asset_id>",
		Short: "Options-market context: term structure, skew, positioning, surface",
		Long: `Pull the options surface and derived context for an asset.

Server sources options from Massive's options snapshot; the response
carries provider, recency, contract count, ATM IV at standard tenors,
term structure (front vs. mid IV), skew (put/call IV gap), and the top
strikes by open interest. Returns ok=false with empty arrays when the
asset has no options coverage.`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrValidation("usage: midaz assets options <asset_id> [--max-days N]")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			params := url.Values{}
			if maxDays > 0 {
				params.Set("max_days", strconv.Itoa(maxDays))
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/assets/" + url.PathEscape(args[0]) + "/options-context",
				Params:    params,
				Normalize: normalizeOptions,
			})
		},
	}
	cmd.Flags().IntVar(&maxDays, "max-days", 0, "Max expiration window in days (server clamps to 7–180; default 90 when unset)")
	return cmd
}

// normalizeOptions lifts provider/recency/contract_count into meta so the
// pretty-print line summarizes coverage at a glance.
func normalizeOptions(body []byte) (any, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}
	provider := cmdutil.UnmarshalString(rawMap["provider"])
	recency := cmdutil.UnmarshalString(rawMap["recency"])
	contracts := cmdutil.UnmarshalInt(rawMap["contract_count"])
	expirations := cmdutil.UnmarshalInt(rawMap["expiration_count"])
	proxyFor := cmdutil.UnmarshalString(rawMap["proxy_for"])

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{
		"provider":         provider,
		"recency":          recency,
		"contract_count":   contracts,
		"expiration_count": expirations,
	}
	if proxyFor != "" {
		meta["proxy_for"] = proxyFor
	}
	return data, meta, nil
}
