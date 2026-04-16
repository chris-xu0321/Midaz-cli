package radar

import (
	"strings"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// newCmdAdd builds `midaz desk radar add`.
//
// Exactly one of --thesis / --topic / --url / --asset / --text is required.
// The command fetches the current radar items, appends the rendered line, and
// writes the full list back via PATCH /api/desk/radar (the backend endpoint
// only supports full-list replacement).
func newCmdAdd(f *cmdutil.Factory) *cobra.Command {
	var (
		thesisID string
		topicID  string
		urlStr   string
		asset    string
		text     string
		title    string
		yes      bool
	)
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add an item to the radar (thesis, topic, url, asset, or free-text note)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if !yes {
				return output.ErrWithHint(output.ExitValidation, "confirmation_required",
					"desk radar add requires --yes",
					"e.g. --asset AAPL --yes")
			}

			kind, err := pickOne(map[string]string{
				"--thesis": thesisID,
				"--topic":  topicID,
				"--url":    urlStr,
				"--asset":  asset,
				"--text":   text,
			})
			if err != nil {
				return err
			}

			c, _, err := cmdutil.AuthedClient(f)
			if err != nil {
				return err
			}

			items, err := fetchCurrentItems(opts.Ctx, c)
			if err != nil {
				return err
			}
			if len(items) >= maxRadarItems {
				return output.ErrWithHint(output.ExitValidation, "radar_full",
					"radar already has the maximum of 12 items",
					"remove one with 'midaz desk radar remove --index N --yes' first")
			}

			var line string
			switch kind {
			case "--thesis":
				label := strings.TrimSpace(title)
				if label == "" {
					label, err = resolveTitle(opts.Ctx, c, "/api/theses/", thesisID, "title")
					if err != nil {
						return err
					}
				}
				line = renderThesisItem(thesisID, label)
			case "--topic":
				label := strings.TrimSpace(title)
				if label == "" {
					label, err = resolveTitle(opts.Ctx, c, "/api/topics/", topicID, "name")
					if err != nil {
						return err
					}
				}
				line = renderTopicItem(topicID, label)
			case "--url":
				if strings.TrimSpace(title) == "" {
					return output.ErrValidation("--title is required with --url")
				}
				line = renderURLItem(urlStr, title)
			case "--asset":
				line = renderAssetItem(asset)
			case "--text":
				line = strings.TrimSpace(text)
				if line == "" {
					return output.ErrValidation("--text is empty after trimming whitespace")
				}
			}

			clipped, truncated := clipTo160(line)
			line = clipped

			if findIndex(items, line) >= 0 {
				return output.ErrWithHint(output.ExitValidation, "already_on_radar",
					"item is already on the radar: "+line,
					"run 'midaz desk radar get' to view current items")
			}

			newItems := append(items, line)
			resp, err := pushItems(opts.Ctx, c, newItems)
			if err != nil {
				return err
			}

			meta := map[string]any{
				"added": line,
				"count": len(newItems),
			}
			if truncated {
				meta["truncated"] = true
			}
			return output.WriteSuccess(opts.Out, resp, meta, opts.Format)
		},
	}
	cmd.Flags().StringVar(&thesisID, "thesis", "", "Thesis ID to pin on the radar")
	cmd.Flags().StringVar(&topicID, "topic", "", "Topic ID to pin on the radar")
	cmd.Flags().StringVar(&urlStr, "url", "", "External URL to add to the radar (requires --title)")
	cmd.Flags().StringVar(&asset, "asset", "", "Ticker symbol to watch (uppercased)")
	cmd.Flags().StringVar(&text, "text", "", "Free-form note (max 160 chars)")
	cmd.Flags().StringVar(&title, "title", "", "Override or provide a human label (used with --thesis/--topic/--url)")
	cmd.Flags().BoolVar(&yes, "yes", false, "Confirm the update")
	return cmd
}

// pickOne validates that exactly one of the given flag/value pairs is non-empty
// and returns its flag name.
func pickOne(flags map[string]string) (string, error) {
	var chosen string
	count := 0
	for k, v := range flags {
		if strings.TrimSpace(v) != "" {
			chosen = k
			count++
		}
	}
	if count == 0 {
		return "", output.ErrWithHint(output.ExitValidation, "validation",
			"specify one of --thesis, --topic, --url, --asset, --text",
			"e.g. --thesis <id> --yes")
	}
	if count > 1 {
		return "", output.ErrWithHint(output.ExitValidation, "validation",
			"specify exactly one of --thesis, --topic, --url, --asset, --text",
			"")
	}
	return chosen, nil
}
