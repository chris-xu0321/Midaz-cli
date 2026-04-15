// Package invite hosts `midaz invite redeem`.
package invite

import (
	"strings"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdInvite builds the invite command tree.
func NewCmdInvite(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invite",
		Short: "Redeem invitation codes",
	}
	cmd.AddCommand(newCmdRedeem(f))
	return cmd
}

func newCmdRedeem(f *cmdutil.Factory) *cobra.Command {
	var yes bool
	cmd := &cobra.Command{
		Use:   "redeem <code>",
		Short: "Redeem an invitation code to unlock the workspace",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return output.ErrValidation("missing required argument: code")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"invite redeem requires --yes",
					"run: midaz invite redeem "+args[0]+" --yes")
			}
			code := strings.ToUpper(strings.TrimSpace(args[0]))
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "POST",
				Path:      "/api/invite/redeem",
				Body:      map[string]any{"code": code},
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm redemption")
	return cmd
}
