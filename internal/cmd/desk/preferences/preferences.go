// Package preferences exposes `midaz desk preferences` — read / write desk
// preferences such as preferred_language.
package preferences

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// SupportedLanguages mirrors Seer's DESK_PREFERRED_LANGUAGES
// (apps/api/src/routes/ws.ts). Keep in sync on rename.
var SupportedLanguages = []string{"en", "zh-CN", "ja", "ko", "es", "fr"}

// NewCmdPreferences builds the preferences subcommand tree.
func NewCmdPreferences(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "preferences",
		Short: "Read or update desk preferences (e.g. language)",
	}
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdSet(f))
	return cmd
}

func newCmdGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Show current preferences (from desk settings)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/desk/settings",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}
