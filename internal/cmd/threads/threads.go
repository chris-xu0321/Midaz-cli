package threads

import (
	"fmt"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdThreads is a deprecated alias of `midaz theses`.
// Kept for one release to avoid breaking existing agent scripts.
func NewCmdThreads(f *cmdutil.Factory) *cobra.Command {
	var status string
	cmd := &cobra.Command{
		Use:        "threads",
		Short:      "Deprecated alias of `theses`",
		Deprecated: "use `midaz theses`",
		Hidden:     true,
		Args:       cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			fmt.Fprintln(opts.ErrOut, "note: `threads` is deprecated — use `midaz theses` instead.")
			params := url.Values{}
			if status != "" {
				params.Set("status", status)
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/theses",
				Params:    params,
				Normalize: cmdutil.NormalizeBareArray,
			})
		},
	}
	cmd.Flags().StringVar(&status, "status", "", "Filter by status")
	return cmd
}
