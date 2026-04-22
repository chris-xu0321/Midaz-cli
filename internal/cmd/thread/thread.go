package thread

import (
	"fmt"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdThread is a deprecated alias of `midaz thesis`.
// Kept for one release to avoid breaking existing agent scripts.
func NewCmdThread(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:        "thread <id>",
		Short:      "Deprecated alias of `thesis`",
		Deprecated: "use `midaz thesis`",
		Hidden:     true,
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
			fmt.Fprintln(opts.ErrOut, "note: `thread` is deprecated — use `midaz thesis` instead.")
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/theses/" + url.PathEscape(args[0]),
				Normalize: threadNormalize,
			})
		},
	}
}

func threadNormalize(body []byte) (any, map[string]any, error) {
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
