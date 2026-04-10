// Package apikey implements the `seer-q api-key` command group.
package apikey

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SparkssL/Midaz-cli/internal/auth"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdAPIKey creates the api-key command group.
func NewCmdAPIKey(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "api-key",
		Short: "Manage Seer API keys",
	}

	cmd.AddCommand(newCreateCmd(f))
	cmd.AddCommand(newListCmd(f))
	cmd.AddCommand(newRevokeCmd(f))

	return cmd
}

// --- create ---

func newCreateCmd(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new API key",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			label, _ := cmd.Flags().GetString("name")
			if label == "" {
				return output.ErrValidation("--name is required")
			}
			return runCreate(f, opts, label)
		},
	}
	cmd.Flags().String("name", "", "Label for the API key (required)")
	_ = cmd.MarkFlagRequired("name")
	return cmd
}

func runCreate(f *cmdutil.Factory, opts *cmdutil.RunOpts, label string) error {
	cfg, err := f.Config()
	if err != nil {
		return output.ErrConfig("cannot load config: %s", err)
	}

	token := auth.ResolveToken(cfg.APIKey)
	if token == "" {
		return output.ErrWithHint(output.ExitAPI, "unauthorized",
			"Not authenticated", "run: seer-q login")
	}

	body, _ := json.Marshal(map[string]string{"label": label})
	req, _ := http.NewRequestWithContext(opts.Ctx, http.MethodPost,
		cfg.APIURL+"/api/app/api-keys", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return output.ErrNetwork("request failed: %s", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 201 {
		return output.ErrAPI("api", "failed to create API key: %s", string(respBody))
	}

	var result map[string]any
	_ = json.Unmarshal(respBody, &result)
	return output.WriteSuccess(opts.Out, result, nil, opts.Format)
}

// --- list ---

func newListCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List your API keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return runList(f, opts)
		},
	}
}

func runList(f *cmdutil.Factory, opts *cmdutil.RunOpts) error {
	cfg, err := f.Config()
	if err != nil {
		return output.ErrConfig("cannot load config: %s", err)
	}

	token := auth.ResolveToken(cfg.APIKey)
	if token == "" {
		return output.ErrWithHint(output.ExitAPI, "unauthorized",
			"Not authenticated", "run: seer-q login")
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet,
		cfg.APIURL+"/api/app/api-keys", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return output.ErrNetwork("request failed: %s", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return output.ErrAPI("api", "failed to list API keys: %s", string(respBody))
	}

	var result []map[string]any
	_ = json.Unmarshal(respBody, &result)
	return output.WriteSuccess(opts.Out, result, nil, opts.Format)
}

// --- revoke ---

func newRevokeCmd(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "revoke <key-id>",
		Short: "Revoke an API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return runRevoke(f, opts, args[0])
		},
	}
}

func runRevoke(f *cmdutil.Factory, opts *cmdutil.RunOpts, keyID string) error {
	cfg, err := f.Config()
	if err != nil {
		return output.ErrConfig("cannot load config: %s", err)
	}

	token := auth.ResolveToken(cfg.APIKey)
	if token == "" {
		return output.ErrWithHint(output.ExitAPI, "unauthorized",
			"Not authenticated", "run: seer-q login")
	}

	ctx := opts.Ctx
	if ctx == nil {
		ctx = context.Background()
	}

	url := fmt.Sprintf("%s/api/app/api-keys/%s", cfg.APIURL, keyID)
	req, _ := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return output.ErrNetwork("request failed: %s", err)
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return output.ErrAPI("api", "failed to revoke API key: %s", string(respBody))
	}

	var result map[string]any
	_ = json.Unmarshal(respBody, &result)
	return output.WriteSuccess(opts.Out, result, nil, opts.Format)
}
