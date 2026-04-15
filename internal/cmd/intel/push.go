package intel

import (
	"os"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

func newCmdPush(f *cmdutil.Factory) *cobra.Command {
	var (
		fromFile    string
		title       string
		sourceURL   string
		publishedAt string
		yes         bool
	)
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Upload a piece of private intel",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"intel push requires --yes",
					"run: midaz intel push --from-file note.md --title \"…\" --yes")
			}
			if fromFile == "" {
				return output.ErrValidation("--from-file is required")
			}
			raw, err := os.ReadFile(fromFile)
			if err != nil {
				return output.ErrConfig("cannot read %s: %s", fromFile, err)
			}
			if len(raw) > 100_000 {
				return output.ErrValidation("content exceeds 100000 chars")
			}
			if _, err := cmdutil.RequireAuth(f); err != nil {
				return err
			}
			body := map[string]any{"content": string(raw)}
			if title != "" {
				body["title"] = title
			}
			if sourceURL != "" {
				body["url"] = sourceURL
			}
			if publishedAt != "" {
				body["published_at"] = publishedAt
			}
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Method:    "POST",
				Path:      "/api/intel",
				Body:      body,
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
	cmd.Flags().StringVar(&fromFile, "from-file", "", "Path to the content file (required)")
	cmd.Flags().StringVar(&title, "title", "", "Optional title")
	cmd.Flags().StringVar(&sourceURL, "url", "", "Optional source URL")
	cmd.Flags().StringVar(&publishedAt, "published-at", "", "Optional ISO-8601 timestamp")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm upload")
	return cmd
}
