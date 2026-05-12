package usage

import (
	"encoding/json"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdUsage builds the usage command tree. The root behaves as the
// top-level summary (/api/usage); `by-run` drills into a single run.
func NewCmdUsage(f *cmdutil.Factory) *cobra.Command {
	var since string

	cmd := &cobra.Command{
		Use:   "usage",
		Short: "Token usage and cost summary",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			params := url.Values{}
			if since != "" {
				params.Set("since", since)
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/usage",
				Params:    params,
				Normalize: normalizeUsage,
			})
		},
	}

	cmd.Flags().StringVar(&since, "since", "24h", "Time period (e.g., 24h, 7d)")
	cmd.AddCommand(newCmdByRun(f))
	return cmd
}

func newCmdByRun(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "by-run <run_id>",
		Short: "Per-pipeline-run token usage breakdown",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrValidation("usage: midaz usage by-run <run_id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/usage/by-run/" + url.PathEscape(args[0]),
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}

func normalizeUsage(body []byte) (any, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}

	viewURL := cmdutil.ExtractViewURL(rawMap)

	sinceVal := cmdutil.UnmarshalString(rawMap["since"])
	totalCalls := cmdutil.UnmarshalInt(rawMap["total_calls"])

	// Assistant usage is a sibling totals block alongside pipeline `by_stage`.
	// Surface its top-line cost + call count in meta so a pretty-print line
	// distinguishes desk-assistant spend from pipeline spend at a glance.
	var assistantCalls int
	var assistantCost float64
	if raw, ok := rawMap["assistant"]; ok {
		var nested struct {
			Total struct {
				CallCount    int     `json:"call_count"`
				TotalCostUSD float64 `json:"total_cost_usd"`
			} `json:"total"`
		}
		_ = json.Unmarshal(raw, &nested)
		assistantCalls = nested.Total.CallCount
		assistantCost = nested.Total.TotalCostUSD
	}

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{
		"total_calls": totalCalls,
	}
	if viewURL != "" {
		meta["view_url"] = viewURL
	}
	if sinceVal != "" {
		meta["since"] = sinceVal
	}
	if assistantCalls > 0 || assistantCost > 0 {
		meta["assistant_calls"] = assistantCalls
		meta["assistant_cost_usd"] = assistantCost
	}

	return data, meta, nil
}
