// Package drivers hosts `midaz drivers`, `midaz driver <id>`, and
// `midaz driver-links` — the canonical world-layer that replaced topics.
package drivers

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// NewCmdDrivers returns `midaz drivers` — list active drivers.
func NewCmdDrivers(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "drivers",
		Short: "List active drivers (world-layer objects behind the market view)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			base := frontendBase(f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path: "/api/drivers",
				Normalize: cmdutil.InjectItemViewURL(func(id string) string {
					if base == "" {
						return ""
					}
					return base + "/market-read?driver=" + id
				}),
			})
		},
	}
}

// frontendBase returns the configured frontend URL ("" if unavailable) so
// list normalizers can synthesize per-item `view_url`s.
func frontendBase(f *cmdutil.Factory) string {
	cfg, err := f.Config()
	if err != nil || cfg == nil {
		return ""
	}
	return cfg.FrontendURL
}
