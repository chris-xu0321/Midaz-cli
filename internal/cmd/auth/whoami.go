package auth

import (
	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdWhoami(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show cached credentials (no network call)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			creds, err := f.Auth()
			if err != nil {
				return output.ErrAuth(err.Error(), "")
			}
			if creds == nil {
				return output.ErrAuth("not logged in", "")
			}
			data := map[string]any{
				"profile":     creds.Profile,
				"user_email":  creds.UserEmail,
				"user_id":     creds.UserID,
				"desk_id":     creds.DeskID,
				"desk_slug":   creds.DeskSlug,
				"key_prefix":  auth.MaskKey(creds.APIKey),
				"verified_at": creds.VerifiedAt,
				"label":       creds.Label,
			}
			return output.WriteSuccess(opts.Out, data, nil, opts.Format)
		},
	}
}
