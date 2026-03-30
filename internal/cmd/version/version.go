package version

import (
	"github.com/chris-xu0321/Midaz-cli/internal/build"
	"github.com/chris-xu0321/Midaz-cli/internal/cmdutil"
	"github.com/chris-xu0321/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func NewCmdVersion(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print CLI version info",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			data := map[string]string{
				"version":    build.Version,
				"go_version": build.GoVersion(),
				"os":         build.OS(),
				"arch":       build.Arch(),
			}
			return output.WriteSuccess(opts.Out, data, nil, opts.Format)
		},
	}
}
