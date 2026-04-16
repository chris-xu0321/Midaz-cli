// Package radar exposes `midaz workspace radar` (get/set).
package radar

import (
	"os"
	"strings"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdRadar builds the radar subcommand tree.
func NewCmdRadar(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "radar",
		Short: "Read or update workspace radar (watchlist)",
	}
	cmd.AddCommand(newCmdGet(f))
	cmd.AddCommand(newCmdSet(f))
	cmd.AddCommand(newCmdAdd(f))
	cmd.AddCommand(newCmdRemove(f))
	cmd.AddCommand(newCmdPin(f))
	cmd.AddCommand(newCmdUnpin(f))
	cmd.AddCommand(newCmdPins(f))
	return cmd
}

func newCmdGet(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "get",
		Short: "Show the current radar",
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
		items    string
		yes      bool
	)
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Replace the radar",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"workspace radar set requires --yes", "e.g. --items \"Fed,AI,Oil\" --yes")
			}
			if fromFile == "" && items == "" {
				return output.ErrValidation("provide --from-file or --items")
			}
			body := map[string]any{}
			switch {
			case fromFile != "":
				raw, err := os.ReadFile(fromFile)
				if err != nil {
					return output.ErrConfig("cannot read %s: %s", fromFile, err)
				}
				body["radar"] = string(raw)
			default:
				parts := strings.Split(items, ",")
				trimmed := make([]string, 0, len(parts))
				for _, p := range parts {
					p = strings.TrimSpace(p)
					if p != "" {
						trimmed = append(trimmed, p)
					}
				}
				if len(trimmed) == 0 {
					return output.ErrValidation("no radar items parsed from --items")
				}
				body["radar_items"] = trimmed
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "PATCH",
				Path:      "/api/ws/radar",
				Body:      body,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Path to a Markdown file")
	cmd.Flags().StringVar(&items, "items", "", "Comma-separated radar items")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the update")
	return cmd
}
