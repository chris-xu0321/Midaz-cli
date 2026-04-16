package radar

import (
	"strings"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// newCmdPin builds `midaz desk radar pin`.
//
// Pins an entity (thesis/topic/driver/asset) to the radar with provenance
// tracking. Distinct from `radar add` — pins participate in L4 refresh and
// render as filled buttons in the web market view.
func newCmdPin(f *cmdutil.Factory) *cobra.Command {
	var (
		kind       string
		sourceType string
		sourceID   string
		label      string
		yes        bool
	)
	cmd := &cobra.Command{
		Use:   "pin",
		Short: "Pin an entity to the radar (entity-level, with provenance)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk radar pin requires --yes",
					"e.g. --kind Thesis --source-type thread --source-id <id> --label \"...\" --yes")
			}
			if strings.TrimSpace(kind) == "" {
				return output.ErrValidation("--kind is required (Thesis | Topic | Driver | Asset)")
			}
			if strings.TrimSpace(sourceType) == "" {
				return output.ErrValidation("--source-type is required (thread | topic | driver | asset)")
			}
			if strings.TrimSpace(sourceID) == "" {
				return output.ErrValidation("--source-id is required")
			}
			if strings.TrimSpace(label) == "" {
				return output.ErrValidation("--label is required")
			}
			if len(label) > 160 {
				return output.ErrValidation("--label exceeds 160 character limit")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method: "POST",
				Path:   "/api/desk/radar/pin",
				Body: map[string]any{
					"kind":       kind,
					"sourceType": sourceType,
					"sourceId":   sourceID,
					"label":      label,
				},
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&kind, "kind", "", "Entity kind: Thesis | Topic | Driver | Asset")
	cmd.Flags().StringVar(&sourceType, "source-type", "", "Source type: thread | topic | driver | asset")
	cmd.Flags().StringVar(&sourceID, "source-id", "", "Entity id or slug")
	cmd.Flags().StringVar(&label, "label", "", "Radar line text (≤ 160 chars)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the pin")
	return cmd
}

// newCmdUnpin builds `midaz desk radar unpin`.
func newCmdUnpin(f *cmdutil.Factory) *cobra.Command {
	var (
		sourceType string
		sourceID   string
		yes        bool
	)
	cmd := &cobra.Command{
		Use:   "unpin",
		Short: "Unpin an entity from the radar",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk radar unpin requires --yes",
					"e.g. --source-type thread --source-id <id> --yes")
			}
			if strings.TrimSpace(sourceType) == "" {
				return output.ErrValidation("--source-type is required (thread | topic | driver | asset)")
			}
			if strings.TrimSpace(sourceID) == "" {
				return output.ErrValidation("--source-id is required")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method: "DELETE",
				Path:   "/api/desk/radar/pin",
				Body: map[string]any{
					"sourceType": sourceType,
					"sourceId":   sourceID,
				},
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&sourceType, "source-type", "", "Source type: thread | topic | driver | asset")
	cmd.Flags().StringVar(&sourceID, "source-id", "", "Entity id or slug")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the unpin")
	return cmd
}

// newCmdPins builds `midaz desk radar pins`.
func newCmdPins(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "pins",
		Short: "List entity pins on the radar (with provenance)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/desk/radar/pins",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}
