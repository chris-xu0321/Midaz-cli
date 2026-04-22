package drivers

import (
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdDriver returns `midaz driver <id>` — driver detail with thread members.
func NewCmdDriver(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "driver <id>",
		Short: "Driver detail with thread members and asset contributions",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"Missing required argument: id",
					"usage: midaz driver <id>")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/drivers/" + url.PathEscape(args[0]),
				Normalize: normalizeDriver,
			})
		},
	}
}

func normalizeDriver(body []byte) (any, map[string]any, error) {
	rawMap, err := cmdutil.ParseMap(body)
	if err != nil {
		return nil, nil, err
	}

	viewURL := cmdutil.ExtractViewURL(rawMap)
	memberCount := cmdutil.CountArray(rawMap["thread_members"])

	data, err := cmdutil.RebuildMap(rawMap)
	if err != nil {
		return nil, nil, err
	}

	meta := map[string]any{"thread_member_count": memberCount}
	if viewURL != "" {
		meta["view_url"] = viewURL
	}
	return data, meta, nil
}
