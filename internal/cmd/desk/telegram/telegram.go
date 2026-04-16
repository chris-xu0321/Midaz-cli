// Package telegram exposes `midaz desk telegram` (status/connect/disconnect).
package telegram

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/SparkssL/Midaz-cli/internal/auth"
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
		Short: "Open the Telegram deep link to connect the bot",
		Args:  cobra.NoArgs,
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
			resp, err := c.Get(cmd.Context(), "/api/desk/settings", nil)
			if err != nil {
				return err
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
				return output.ErrAPI("api", "failed to parse settings: %s", err)
			}
			if settings.Telegram.BotUsername == "" {
				return output.ErrAPI("api", "TELEGRAM_BOT_USERNAME not configured on the API")
			}
			deskID := settings.Desk.ID
			if deskID == "" {
				deskID = creds.DeskID
			}
			deepLink := fmt.Sprintf("https://t.me/%s?start=%s",
				url.PathEscape(settings.Telegram.BotUsername), url.QueryEscape(deskID))

			_ = auth.OpenBrowser(deepLink)

			data := map[string]any{
				"url":          deepLink,
				"bot_username": settings.Telegram.BotUsername,
				"connected":    settings.Telegram.Connected,
			}
			meta := map[string]any{
				"hint":     "Tap 'Start' in Telegram, then run: midaz desk telegram status",
				"view_url": deepLink,
			}
			return output.WriteSuccess(opts.Out, data, meta, opts.Format)
		},
	}
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
