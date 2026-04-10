// Package login implements the `seer-q login` command.
//
// Opens a browser to the Seer web login page, starts a localhost callback
// server, and exchanges the auth result for a Seer PAT (sk_...).
package login

import (
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdLogin creates the login command.
func NewCmdLogin(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: "Authenticate with Seer via browser",
		Long:  "Opens a browser to sign in with Google or email. Stores a Seer API key locally for subsequent CLI use.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runLogin(f, cmd)
		},
	}

	cmd.Flags().Bool("status", false, "Show current login status without logging in")

	return cmd
}

func runLogin(f *cmdutil.Factory, cmd *cobra.Command) error {
	opts := cmdutil.ResolveRunOpts(cmd, f)

	// --status: just show current state
	if showStatus, _ := cmd.Flags().GetBool("status"); showStatus {
		return showLoginStatus(opts)
	}

	cfg, err := f.Config()
	if err != nil {
		return output.ErrConfig("cannot load config: %s", err)
	}

	frontendURL := cfg.FrontendURL

	// Start localhost callback server
	ctx := cmd.Context()
	port, resultCh, cleanup := auth.StartCallbackServer(ctx)
	defer cleanup()

	if port == 0 {
		return output.ErrNetwork("failed to start local callback server")
	}

	// Open browser
	loginURL := fmt.Sprintf("%s/login?cli=true&port=%d", frontendURL, port)
	fmt.Fprintf(opts.ErrOut, "Opening browser: %s\n", loginURL)
	fmt.Fprintf(opts.ErrOut, "Waiting for authentication...\n")

	if err := openBrowser(loginURL); err != nil {
		fmt.Fprintf(opts.ErrOut, "Could not open browser. Please visit:\n  %s\n", loginURL)
	}

	// Wait for callback (60s timeout)
	select {
	case result := <-resultCh:
		if result.APIKey == "" {
			return output.ErrAPI("auth", "no API key received from server")
		}

		// Save credentials
		creds := &auth.Credentials{
			APIKey:        result.APIKey,
			WorkspaceID:   result.WorkspaceID,
			WorkspaceSlug: result.WorkspaceSlug,
			UserEmail:     result.UserEmail,
		}
		if err := auth.Save(creds); err != nil {
			return output.ErrConfig("failed to save credentials: %s", err)
		}

		loginInfo := map[string]any{
			"ok":    true,
			"email": result.UserEmail,
			"workspace": map[string]any{
				"id":   result.WorkspaceID,
				"slug": result.WorkspaceSlug,
			},
			"credentials_path": auth.CredentialsPath(),
		}

		return output.WriteSuccess(opts.Out, loginInfo, nil, opts.Format)

	case <-time.After(60 * time.Second):
		return output.ErrNetwork("authentication timed out after 60s")
	case <-ctx.Done():
		return output.ErrNetwork("authentication cancelled")
	}
}

func showLoginStatus(opts *cmdutil.RunOpts) error {
	creds, err := auth.Load()
	if err != nil {
		return output.ErrConfig("failed to read credentials: %s", err)
	}

	if creds == nil || creds.APIKey == "" {
		status := map[string]any{
			"ok":             true,
			"authenticated":  false,
			"hint":           "run: seer-q login",
		}
		return output.WriteSuccess(opts.Out, status, nil, opts.Format)
	}

	status := map[string]any{
		"ok":            true,
		"authenticated": true,
		"email":         creds.UserEmail,
		"workspace": map[string]any{
			"id":   creds.WorkspaceID,
			"slug": creds.WorkspaceSlug,
		},
		"key_prefix":       creds.APIKey[:min(11, len(creds.APIKey))],
		"credentials_path": auth.CredentialsPath(),
	}
	return output.WriteSuccess(opts.Out, status, nil, opts.Format)
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
	return cmd.Start()
}
