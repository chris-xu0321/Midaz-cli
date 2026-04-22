// Package thesis hosts `midaz thesis <id>` — thesis detail + claims + links.
package thesis

import (
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdThesis returns the `thesis` cobra command.
func NewCmdThesis(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "thesis <id>",
		Short: "Thesis detail + claims + market links",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"Missing required argument: id",
					"usage: midaz thesis <id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/theses/" + url.PathEscape(args[0]),
				Normalize: normalize,
			})
		},
	}
}

func normalize(body []byte) (any, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}

	viewURL := cmdutil.ExtractViewURL(rawMap)
	claimCount := cmdutil.CountArray(rawMap["claims"])
	marketLinkCount := cmdutil.CountArray(rawMap["market_links"])
	supporting := cmdutil.UnmarshalInt(rawMap["supporting_count"])
	contradicting := cmdutil.UnmarshalInt(rawMap["contradicting_count"])

	delete(rawMap, "has_market_link")
	delete(rawMap, "market_link_count")

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{
		"claim_count":         claimCount,
		"supporting_count":    supporting,
		"contradicting_count": contradicting,
		"market_link_count":   marketLinkCount,
	}
	if viewURL != "" {
		meta["view_url"] = viewURL
	}
	return data, meta, nil
}
