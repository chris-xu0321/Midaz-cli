package doctor

import (
	"context"
	"os"
	"time"

	authpkg "github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/client"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/config"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

type check struct {
	Name    string `json:"name"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewCmdDoctor(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnostic checks",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return runDoctor(f, opts)
		},
	}
}

func runDoctor(f *cmdutil.Factory, opts *cmdutil.RunOpts) error {
	var checks []check
	passed, failed, warned := 0, 0, 0

	configPath := config.ConfigPath()
	if _, err := os.Stat(configPath); err == nil {
		checks = append(checks, check{"config_source", "pass", "file at " + configPath})
		passed++
	} else if envAny("MIDAZ_API_URL", "SEER_API_URL", "MIDAZ_FRONTEND_URL", "SEER_FRONTEND_URL") {
		checks = append(checks, check{"config_source", "pass", "env vars"})
		passed++
	} else {
		checks = append(checks, check{"config_source", "pass", "defaults only"})
		passed++
	}

	cfg, err := f.Config()
	if err != nil {
		checks = append(checks, check{"api_url", "fail", "config error: " + err.Error()})
		failed++
	} else {
		checks = append(checks, check{"api_url", "pass", cfg.APIURL})
		passed++
	}

	if cfg != nil {
		c := client.New(cfg.APIURL)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, apiErr := c.Get(ctx, "/api/health", nil)
		if apiErr != nil {
			checks = append(checks, check{"api_reachable", "fail", apiErr.Error()})
			failed++
		} else {
			checks = append(checks, check{"api_reachable", "pass", "GET /api/health returned ok"})
			passed++
		}
	} else {
		checks = append(checks, check{"api_reachable", "fail", "skipped — no config"})
		failed++
	}

	if cfg != nil && cfg.FrontendURL != "" {
		checks = append(checks, check{"frontend_url", "pass", cfg.FrontendURL})
		passed++
	} else {
		checks = append(checks, check{"frontend_url", "warn", "not configured"})
		warned++
	}

	if _, err := os.Stat(configPath); err == nil {
		checks = append(checks, check{"config_file", "pass", configPath})
		passed++
	} else {
		checks = append(checks, check{"config_file", "warn", "not found at " + configPath})
		warned++
	}

	if creds, err := f.Auth(); err == nil && creds != nil {
		msg := loginDescription(creds)
		checks = append(checks, check{"auth", "pass", msg})
		passed++
	} else {
		checks = append(checks, check{"auth", "warn", "not logged in — run 'midaz auth login'"})
		warned++
	}

	data := map[string]any{"checks": checks}
	meta := map[string]any{
		"passed": passed,
		"failed": failed,
		"warned": warned,
	}
	return output.WriteSuccess(opts.Out, data, meta, opts.Format)
}

func envAny(names ...string) bool {
	for _, n := range names {
		if os.Getenv(n) != "" {
			return true
		}
	}
	return false
}

func loginDescription(c *authpkg.Creds) string {
	email := c.UserEmail
	if email == "" {
		email = c.UserID
	}
	if email == "" {
		email = "unknown user"
	}
	ws := c.WorkspaceSlug
	if ws == "" {
		ws = c.WorkspaceID
	}
	return email + " @ " + ws + " (" + authpkg.MaskKey(c.APIKey) + ")"
}
