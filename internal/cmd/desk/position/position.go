// Package position hosts `midaz desk position` (open/update/close).
//
// Positions are DB-owned trader stances on a desk; the API stores
// (asset, bias_direction, entry_thesis) and L4 fills health/prose around
// them. Opening or updating a position enqueues an L4 rebuild so the
// downstream desk view picks the change up on the next refresh.
package position

import (
	"context"
	"encoding/json"

	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdPosition builds the `desk position` subcommand tree.
func NewCmdPosition(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "position",
		Short: "Open, update, and close trader positions on your desk",
	}
	cmd.AddCommand(newCmdOpen(f))
	cmd.AddCommand(newCmdUpdate(f))
	cmd.AddCommand(newCmdClose(f))
	return cmd
}

// resolveDeskSlug returns the slug (or id, whichever is set) for the
// authenticated user's desk. Mirrors the inline pattern in
// internal/cmd/desk/view.go so /api/desks/{slug}/positions works whether
// or not the slug was cached in auth.json.
func resolveDeskSlug(ctx context.Context, f *cmdutil.Factory, creds *auth.Creds) (string, error) {
	slug := auth.NonEmpty(creds.DeskSlug, creds.DeskID)
	if slug != "" {
		return slug, nil
	}
	c, err := f.Client()
	if err != nil {
		return "", err
	}
	resp, err := c.Get(ctx, "/api/desk", nil)
	if err != nil {
		return "", err
	}
	var d struct {
		Desk struct {
			ID   string `json:"id"`
			Slug string `json:"slug"`
		} `json:"desk"`
	}
	if err := json.Unmarshal(resp.Body, &d); err != nil {
		return "", output.ErrAPI("api", "failed to parse desk: %s", err)
	}
	slug = auth.NonEmpty(d.Desk.Slug, d.Desk.ID)
	if slug == "" {
		return "", output.ErrAPI("api", "no desk slug or id found for current credentials")
	}
	return slug, nil
}
