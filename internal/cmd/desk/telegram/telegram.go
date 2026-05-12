// Package telegram exposes `midaz desk telegram` (status/connect/disconnect).
package telegram

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/client"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdTelegram builds the telegram subcommand tree.
func NewCmdTelegram(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "telegram",
		Short: "Connect/disconnect the desk Telegram bot",
	}
	cmd.AddCommand(newCmdStatus(f))
	cmd.AddCommand(newCmdConnect(f))
	cmd.AddCommand(newCmdDisconnect(f))
	return cmd
}

func newCmdStatus(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show Telegram connection status",
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

func newCmdConnect(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "connect",
		Short: "Generate a one-time Telegram deep link to link the bot",
		Long: `Generate a token-backed Telegram deep link via
POST /api/desk/telegram/link-token. The server returns a 10-minute
single-use start payload; the CLI prints the URL and auto-opens it.

If the API rejects the new route (older Seer), the command falls back
to the legacy bot-username deep link derived from /api/desk/settings.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			creds, err := cmdutil.RequireAuth(f)
			if err != nil {
				return err
			}
			c, err := f.Client()
			if err != nil {
				return err
			}

			data, meta, err := tryLinkToken(cmd.Context(), c)
			if err != nil {
				if !isNotFound(err) {
					return err
				}
				// Older Seer: legacy bot-username flow.
				data, meta, err = legacyConnect(cmd.Context(), c, creds)
				if err != nil {
					return err
				}
			}
			if u, ok := data["url"].(string); ok && u != "" {
				_ = auth.OpenBrowser(u)
			}
			return output.WriteSuccess(opts.Out, data, meta, opts.Format)
		},
	}
}

// tryLinkToken posts to /api/desk/telegram/link-token and shapes the
// success payload into (data, meta).
func tryLinkToken(ctx context.Context, c *client.Client) (map[string]any, map[string]any, error) {
	resp, err := c.Post(ctx, "/api/desk/telegram/link-token", nil)
	if err != nil {
		return nil, nil, err
	}
	var parsed struct {
		URL             string `json:"url"`
		StartCommand    string `json:"start_command"`
		TokenExpiresAt  string `json:"token_expires_at"`
	}
	if err := json.Unmarshal(resp.Body, &parsed); err != nil {
		return nil, nil, output.ErrAPI("api", "failed to parse link-token response: %s", err)
	}
	if parsed.URL == "" {
		return nil, nil, output.ErrAPI("api", "link-token response missing url")
	}
	data := map[string]any{
		"url":               parsed.URL,
		"start_command":     parsed.StartCommand,
		"token_expires_at":  parsed.TokenExpiresAt,
		"flow":              "link_token",
	}
	meta := map[string]any{
		"hint":     "Tap 'Start' in Telegram within 10 minutes, then run: midaz desk telegram status",
		"view_url": parsed.URL,
	}
	return data, meta, nil
}

// legacyConnect mirrors the pre-0.8.0 bot-username flow for compatibility
// with older Seer that hasn't shipped /api/desk/telegram/link-token.
func legacyConnect(ctx context.Context, c *client.Client, creds *auth.Creds) (map[string]any, map[string]any, error) {
	resp, err := c.Get(ctx, "/api/desk/settings", nil)
	if err != nil {
		return nil, nil, err
	}
	var settings struct {
		Telegram struct {
			Connected   bool   `json:"connected"`
			BotUsername string `json:"bot_username"`
		} `json:"telegram"`
		Desk struct {
			ID string `json:"id"`
		} `json:"desk"`
	}
	if err := json.Unmarshal(resp.Body, &settings); err != nil {
		return nil, nil, output.ErrAPI("api", "failed to parse settings: %s", err)
	}
	if settings.Telegram.BotUsername == "" {
		return nil, nil, output.ErrAPI("api", "TELEGRAM_BOT_USERNAME not configured on the API")
	}
	deskID := settings.Desk.ID
	if deskID == "" {
		deskID = creds.DeskID
	}
	deepLink := fmt.Sprintf("https://t.me/%s?start=%s",
		url.PathEscape(settings.Telegram.BotUsername), url.QueryEscape(deskID))

	data := map[string]any{
		"url":          deepLink,
		"bot_username": settings.Telegram.BotUsername,
		"connected":    settings.Telegram.Connected,
		"flow":         "legacy_bot_username",
	}
	meta := map[string]any{
		"hint":     "Tap 'Start' in Telegram, then run: midaz desk telegram status",
		"view_url": deepLink,
	}
	return data, meta, nil
}

// isNotFound reports whether err comes from a 404 classification (the
// legacy fallback signal — older Seer doesn't mount link-token).
func isNotFound(err error) bool {
	var exitErr *output.ExitError
	if !errors.As(err, &exitErr) {
		return false
	}
	return exitErr.Detail != nil && exitErr.Detail.Code == "not_found"
}

func newCmdDisconnect(f *cmdutil.Factory) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "disconnect",
		Short: "Disconnect the Telegram bot",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk telegram disconnect requires --yes",
					"run: midaz desk telegram disconnect --yes")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "DELETE",
				Path:      "/api/desk/telegram",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm disconnection")
	return cmd
}
