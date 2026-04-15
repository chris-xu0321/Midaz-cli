// Package playbook exposes `midaz workspace playbook` (get/set).
package playbook

import (
	"os"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdPlaybook builds the playbook subcommand tree.
func NewCmdPlaybook(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "playbook",
		Short: "Read or update workspace playbook (trading rules)",
	}
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdSet(f))
	return cmd
}

func newCmdGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Show the current playbook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/ws/settings",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}

func newCmdSet(f *cmdutil.Factory) *cobra.Command {
	var (
		fromFile string
		yes      bool
	)
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Replace the playbook",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"workspace playbook set requires --yes", "run: midaz workspace playbook set --from-file playbook.md --yes")
			}
			if fromFile == "" {
				return output.ErrValidation("--from-file is required")
			}
			raw, err := os.ReadFile(fromFile)
			if err != nil {
				return output.ErrConfig("cannot read %s: %s", fromFile, err)
			}
			if len(raw) > 20_000 {
				return output.ErrValidation("playbook exceeds 20000 chars")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "PATCH",
				Path:      "/api/ws/playbook",
				Body:      map[string]any{"playbook": string(raw)},
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Path to the playbook Markdown file (required)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the update")
	return cmd
}
