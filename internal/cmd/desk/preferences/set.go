package preferences

import (
	"slices"
	"strings"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdSet(f *cmdutil.Factory) *cobra.Command {
	var (
		language string
		yes      bool
	)
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Update desk preferences",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk preferences set requires --yes",
					"e.g. midaz desk preferences set --language en --yes")
			}
			if language == "" {
				return output.ErrValidation("--language is required")
			}
			if !isSupportedLanguage(language) {
				return output.ErrWithHint(output.ExitValidation, "validation",
					"unsupported language: "+language,
					"supported: "+strings.Join(SupportedLanguages, ", "))
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "PATCH",
				Path:      "/api/desk/preferences",
				Body:      map[string]any{"preferred_language": language},
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&language, "language", "",
		"Preferred language ("+strings.Join(SupportedLanguages, ", ")+")")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the update")
	return cmd
}

func isSupportedLanguage(lang string) bool {
	return slices.Contains(SupportedLanguages, lang)
}
