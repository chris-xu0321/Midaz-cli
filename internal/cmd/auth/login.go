package auth

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/client"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/config"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdLogin(f *cmdutil.Factory) *cobra.Command {
	var (
		paste     bool
		token     string
		label     string
		timeoutSe int
	)
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Sign in to Midaz",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			cfg, err := f.Config()
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(cmd.Context(), time.Duration(timeoutSe)*time.Second)
			defer cancel()

			switch {
			case token != "":
				return completeLogin(opts, cfg, opts.Profile, token, "MIDAZ_TOKEN flag", ctx)
			case paste:
				return pasteLogin(opts, cfg, opts.Profile, ctx)
			default:
				return browserLogin(opts, cfg, opts.Profile, label, ctx)
			}
		},
	}
	cmd.Flags().BoolVar(&paste, "paste", false, "Paste a PAT you created on the website")
	cmd.Flags().StringVar(&token, "token", "", "Provide a PAT inline (skips browser)")
	cmd.Flags().StringVar(&label, "label", "", "PAT label (default: cli:<hostname>:<date>)")
	cmd.Flags().IntVar(&timeoutSe, "timeout", 180, "Login timeout in seconds")
	return cmd
}

func browserLogin(opts *cmdutil.RunOpts, cfg *config.Config, profile, label string, ctx context.Context) error {
	listener, err := auth.NewListener()
	if err != nil {
		return output.ErrWithHint(output.ExitInternal, "internal", err.Error(),
			"fallback: midaz auth login --paste")
	}
	defer listener.Close()

	listener.Origin = frontendOrigin(cfg)

	if label == "" {
		label = defaultLabel()
	}
	// The /cli-auth page rejects labels >80 chars; truncate gracefully
	// rather than fail the whole flow.
	if len(label) > 80 {
		label = label[:80]
	}

	relayURL := buildRelayURL(cfg.FrontendURL, listener.Port, listener.Nonce, label)
	fmt.Fprintln(opts.ErrOut, "Opening browser for sign-in at "+cfg.FrontendURL+"/cli-auth …")
	fmt.Fprintln(opts.ErrOut, "  ", relayURL)
	fmt.Fprintln(opts.ErrOut, "Waiting for the sign-in to complete (timeout: 3m). Press Ctrl+C to abort.")
	if err := auth.OpenBrowser(relayURL); err != nil {
		fmt.Fprintln(opts.ErrOut, "(could not auto-open browser — copy the URL above into your browser)")
	}

	res, err := listener.Serve(ctx, 3*time.Minute)
	if err != nil {
		return output.ErrWithHint(output.ExitAuth, "auth", err.Error(),
			"try again, or run: midaz auth login --paste")
	}

	creds := &auth.Credentials{
		APIKey:        res.APIKey,
		WorkspaceID:   res.WorkspaceID,
		WorkspaceSlug: res.WorkspaceSlug,
		UserEmail:     res.UserEmail,
		UserID:        res.UserID,
		Label:         label,
	}
	return storeAndReport(opts, profile, creds)
}

func pasteLogin(opts *cmdutil.RunOpts, cfg *config.Config, profile string, ctx context.Context) error {
	fmt.Fprintln(opts.ErrOut, "1. Visit", cfg.FrontendURL+"/workspace/settings")
	fmt.Fprintln(opts.ErrOut, "2. Open the API keys section and create a new key")
	fmt.Fprintln(opts.ErrOut, "3. Paste the key below and press Enter")
	fmt.Fprint(opts.ErrOut, "PAT> ")

	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return output.ErrAuth("failed to read token: "+err.Error(), "")
	}
	token := strings.TrimSpace(line)
	if !strings.HasPrefix(token, "sk_") {
		return output.ErrAuth("invalid token (expected sk_… prefix)", "")
	}
	return completeLogin(opts, cfg, profile, token, "paste", ctx)
}

// completeLogin validates the token against /api/app/me and persists credentials.
func completeLogin(opts *cmdutil.RunOpts, cfg *config.Config, profile, token, source string, ctx context.Context) error {
	me, err := fetchMe(ctx, cfg.APIURL, token)
	if err != nil {
		return err
	}
	creds := &auth.Credentials{
		APIKey:        token,
		WorkspaceID:   me.WorkspaceID,
		WorkspaceSlug: me.WorkspaceSlug,
		UserEmail:     me.Email,
		UserID:        me.UserID,
		Label:         source,
	}
	return storeAndReport(opts, profile, creds)
}

func storeAndReport(opts *cmdutil.RunOpts, profile string, creds *auth.Credentials) error {
	if err := auth.SetCurrent(profile, creds); err != nil {
		return output.ErrWithHint(output.ExitInternal, "internal", err.Error(), "")
	}
	data := map[string]any{
		"user_email":     creds.UserEmail,
		"user_id":        creds.UserID,
		"workspace_id":   creds.WorkspaceID,
		"workspace_slug": creds.WorkspaceSlug,
		"profile":        auth.NonEmpty(profile, auth.DefaultProfile),
		"key_prefix":     keyPrefix(creds.APIKey),
		"auth_file":      auth.AuthPath(),
	}
	meta := map[string]any{
		"message": "Signed in. Credentials stored at " + auth.AuthPath(),
	}
	return output.WriteSuccess(opts.Out, data, meta, opts.Format)
}

type meResponse struct {
	UserID        string
	Email         string
	WorkspaceID   string
	WorkspaceSlug string
}

func fetchMe(ctx context.Context, apiURL, token string) (*meResponse, error) {
	c := client.New(apiURL).WithToken(token)
	resp, err := c.Get(ctx, "/api/app/me", nil)
	if err != nil {
		return nil, err
	}
	var parsed struct {
		User struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
		Workspace *struct {
			ID   string `json:"id"`
			Slug string `json:"slug"`
		} `json:"workspace"`
	}
	if err := json.Unmarshal(resp.Body, &parsed); err != nil {
		return nil, output.ErrAPI("api", "failed to parse /api/app/me response: %s", err)
	}
	out := &meResponse{
		UserID: parsed.User.ID,
		Email:  parsed.User.Email,
	}
	if parsed.Workspace != nil {
		out.WorkspaceID = parsed.Workspace.ID
		out.WorkspaceSlug = parsed.Workspace.Slug
	}
	return out, nil
}

func buildRelayURL(frontendURL string, port int, nonce, label string) string {
	base := strings.TrimRight(frontendURL, "/")
	q := url.Values{}
	q.Set("port", fmt.Sprintf("%d", port))
	q.Set("nonce", nonce)
	if label != "" {
		q.Set("label", label)
	}
	return base + "/cli-auth?" + q.Encode()
}

func frontendOrigin(cfg *config.Config) string {
	u, err := url.Parse(cfg.FrontendURL)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return "https://www.midaz.xyz"
	}
	return u.Scheme + "://" + u.Host
}

func defaultLabel() string {
	host, _ := os.Hostname()
	if host == "" {
		host = "cli"
	}
	return fmt.Sprintf("cli:%s:%s", host, time.Now().UTC().Format("2006-01-02"))
}

func keyPrefix(k string) string {
	if len(k) < 11 {
		return ""
	}
	return k[:11]
}

