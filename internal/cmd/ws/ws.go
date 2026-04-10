// Package ws implements the `seer-q ws` command group.
// The workspace is a trader's private desk: identity, radar, playbook, view, share.
package ws

import (
	"encoding/json"
	"io"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// uuidRE validates the workspace_id argument on `ws view <workspace_id>`
// so a typo doesn't get quietly appended to the URL and 400'd server-side
// with an opaque message.
var uuidRE = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// NewCmdWs creates the ws command with subcommands.
func NewCmdWs(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ws",
		Short: "Your workspace — identity, radar, playbook, view",
		Long:  "Your private working desk. Without subcommand, shows workspace overview.",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/ws",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}

	cmd.AddCommand(newCmdRadar(f))
	cmd.AddCommand(newCmdPlaybook(f))
	cmd.AddCommand(newCmdOnboard(f))
	cmd.AddCommand(newCmdView(f))
	cmd.AddCommand(newCmdShare(f))
	cmd.AddCommand(newCmdUnshare(f))
	cmd.AddCommand(newCmdAlerts(f))

	return cmd
}

func newCmdRadar(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "radar <text | @file>",
		Short: "Set what you're watching",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			text, err := resolveTextArg(args[0])
			if err != nil {
				return output.Errorf(output.ExitValidation, "validation", "%s", err)
			}
			body, _ := json.Marshal(map[string]string{"radar": text})
			return cmdutil.RunMutatingAPICommand(f, opts, &cmdutil.MutatingAPISpec{
				Method: "PATCH", Path: "/api/ws/radar", Body: body,
			})
		},
	}
}

func newCmdPlaybook(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "playbook <text | @file>",
		Short: "Set how you trade",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			text, err := resolveTextArg(args[0])
			if err != nil {
				return output.Errorf(output.ExitValidation, "validation", "%s", err)
			}
			body, _ := json.Marshal(map[string]string{"playbook": text})
			return cmdutil.RunMutatingAPICommand(f, opts, &cmdutil.MutatingAPISpec{
				Method: "PATCH", Path: "/api/ws/playbook", Body: body,
			})
		},
	}
}

// newCmdOnboard writes radar + playbook + onboarding_completed_at in a
// single API call. `seer-q ws radar` / `ws playbook` update one field at
// a time and do NOT set onboarding_completed_at; until this command runs
// `seer-q ws` will keep reporting `onboarded: false`.
func newCmdOnboard(f *cmdutil.Factory) *cobra.Command {
	var radarArg, playbookArg string
	cmd := &cobra.Command{
		Use:   "onboard --radar <text|@file> --playbook <text|@file>",
		Short: "Atomically set radar + playbook and mark workspace onboarded",
		Long: `Initial onboarding call. Sets radar, playbook, and onboarding_completed_at
in a single API request, and triggers L4 synthesis once with reason=onboard.

Use 'seer-q ws radar' / 'seer-q ws playbook' for updates AFTER onboarding.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if radarArg == "" || playbookArg == "" {
				return output.Errorf(output.ExitValidation, "validation",
					"both --radar and --playbook are required")
			}
			radar, err := resolveTextArg(radarArg)
			if err != nil {
				return output.Errorf(output.ExitValidation, "validation", "%s", err)
			}
			playbook, err := resolveTextArg(playbookArg)
			if err != nil {
				return output.Errorf(output.ExitValidation, "validation", "%s", err)
			}
			opts := cmdutil.ResolveRunOpts(cmd, f)
			body, _ := json.Marshal(map[string]string{
				"radar": radar, "playbook": playbook,
			})
			return cmdutil.RunMutatingAPICommand(f, opts, &cmdutil.MutatingAPISpec{
				Method: "POST", Path: "/api/ws/onboard", Body: body,
			})
		},
	}
	cmd.Flags().StringVar(&radarArg, "radar", "", "Radar text or @file")
	cmd.Flags().StringVar(&playbookArg, "playbook", "", "Playbook text or @file")
	return cmd
}

// newCmdView fetches the caller's personal view by default. With a
// workspace_id argument it fetches another workspace's current view via
// the auth-required share endpoint (`GET /api/workspaces/:id/view`).
// Workspace members can always read their own view; non-members need
// `shared = true` on the target workspace.
func newCmdView(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "view [workspace_id]",
		Short: "Your personal market cognitive view",
		Long: `Without arguments: your own personal view (market baseline + your radar + your intel).
With a workspace_id: someone else's view via the auth-required share endpoint
(GET /api/workspaces/:workspace_id/view). The viewer must be logged in, and
the target workspace must have shared = true (unless the viewer is a member).`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if len(args) == 1 {
				if !uuidRE.MatchString(args[0]) {
					return output.Errorf(output.ExitValidation, "validation",
						"workspace_id must be a UUID")
				}
				return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
					Path:      "/api/workspaces/" + url.PathEscape(args[0]) + "/view",
					Normalize: cmdutil.NormalizePassthrough,
				})
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/ws/view",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}

// newCmdShare flips workspaces.shared = true via PATCH /api/ws.
// The workspace_id is the share handle; any logged-in user who knows it
// can then read the current view via 'seer-q ws view <workspace_id>'.
func newCmdShare(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "share",
		Short: "Mark your workspace as shared (readable by any logged-in user)",
		Long: `Flip workspaces.shared = true. Any logged-in user who knows your
workspace_id can then read your current view via
'seer-q ws view <workspace_id>'.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			body, _ := json.Marshal(map[string]bool{"shared": true})
			return cmdutil.RunMutatingAPICommand(f, opts, &cmdutil.MutatingAPISpec{
				Method: "PATCH", Path: "/api/ws", Body: body,
			})
		},
	}
}

func newCmdUnshare(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "unshare",
		Short: "Revoke shared access to your workspace",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			body, _ := json.Marshal(map[string]bool{"shared": false})
			return cmdutil.RunMutatingAPICommand(f, opts, &cmdutil.MutatingAPISpec{
				Method: "PATCH", Path: "/api/ws", Body: body,
			})
		},
	}
}

// newCmdAlerts exposes the workspace_alerts feed written by L4.
//   seer-q ws alerts                # list unread
//   seer-q ws alerts --all           # include already-read
//   seer-q ws alerts --limit N       # cap result count
//   seer-q ws alerts read <id>       # mark a single alert as read
func newCmdAlerts(f *cmdutil.Factory) *cobra.Command {
	var includeRead bool
	var limit int

	cmd := &cobra.Command{
		Use:   "alerts",
		Short: "Your L4 alert feed",
		Long: `Alerts are generated by L4 when a pipeline refresh or user edit
materially changes your personal view. By default only unread alerts
are returned — pass --all to include previously-read alerts too.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			query := url.Values{}
			if includeRead {
				query.Set("include_read", "true")
			}
			if limit > 0 {
				query.Set("limit", strconv.Itoa(limit))
			}
			path := "/api/ws/alerts"
			if q := query.Encode(); q != "" {
				path += "?" + q
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      path,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().BoolVar(&includeRead, "all", false, "Include already-read alerts")
	cmd.Flags().IntVar(&limit, "limit", 0, "Max alerts to return (default: server default)")

	cmd.AddCommand(&cobra.Command{
		Use:   "read <alert_id>",
		Short: "Mark an alert as read",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunMutatingAPICommand(f, opts, &cmdutil.MutatingAPISpec{
				Method: "PATCH",
				Path:   "/api/ws/alerts/" + url.PathEscape(args[0]) + "/read",
			})
		},
	})

	return cmd
}

func resolveTextArg(arg string) (string, error) {
	if len(arg) > 1 && arg[0] == '@' {
		f, err := os.Open(arg[1:])
		if err != nil {
			return "", err
		}
		defer f.Close()
		data, err := io.ReadAll(f)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return arg, nil
}
