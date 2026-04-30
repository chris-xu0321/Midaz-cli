package desk

import (
	"encoding/json"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdView(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "view",
		Short: "Personal market read (subscription-gated)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			creds, err := cmdutil.RequireAuth(f)
			if err != nil {
				return err
			}
			slug := auth.NonEmpty(creds.DeskSlug, creds.DeskID)
			if slug == "" {
				c, err := f.Client()
				if err != nil {
					return err
				}
				resp, err := c.Get(cmd.Context(), "/api/desk", nil)
				if err != nil {
					return err
				}
				var d struct {
					Desk struct {
						ID   string `json:"id"`
						Slug string `json:"slug"`
					} `json:"desk"`
				}
				if err := json.Unmarshal(resp.Body, &d); err != nil {
					return output.ErrAPI("api", "failed to parse desk: %s", err)
				}
				slug = auth.NonEmpty(d.Desk.Slug, d.Desk.ID)
				if slug == "" {
					return output.ErrAPI("api", "no desk slug or id found for current credentials")
				}
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/desks/" + url.PathEscape(slug) + "/read",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}
