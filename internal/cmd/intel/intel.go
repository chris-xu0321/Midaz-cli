// Package intel implements the `seer-q intel` command group.
// Intel is any information a trader pushes into their workspace:
// ideas, notes, observations, articles, research.
package intel

import (
	"encoding/json"
	"io"
	"net/url"
	"os"

	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// NewCmdIntel creates the intel command group.
func NewCmdIntel(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "intel <content | @file | -> | list | rm",
		Short: "Push intel into your workspace",
		Long: `Feed any information into your workspace — ideas, notes, articles, observations.

Without subcommand, pushes content:
  seer-q intel "Fed raised rates by 25bp"
  seer-q intel "OPEC+ considering cut" --title "OPEC Intel"
  seer-q intel @research.md
  cat article.txt | seer-q intel -

Subcommands:
  seer-q intel list           List your intel
  seer-q intel rm <id>        Delete intel`,
		Args:                  cobra.ArbitraryArgs,
		DisableFlagParsing:    false,
		DisableFlagsInUseLine: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cmd.Help()
			}
			// Default action: push
			return pushIntel(f, cmd, args)
		},
	}

	var title, sourceURL string
	cmd.Flags().StringVarP(&title, "title", "t", "", "Title (auto-generated if omitted)")
	cmd.Flags().StringVarP(&sourceURL, "url", "u", "", "Source URL if applicable")

	cmd.AddCommand(newCmdList(f))
	cmd.AddCommand(newCmdRm(f))

	return cmd
}

func pushIntel(f *cmdutil.Factory, cmd *cobra.Command, args []string) error {
	opts := cmdutil.ResolveRunOpts(cmd, f)

	content, err := resolveInput(args[0])
	if err != nil {
		return output.Errorf(output.ExitValidation, "validation", "%s", err)
	}
	if content == "" {
		return output.Errorf(output.ExitValidation, "validation", "content is empty")
	}

	payload := map[string]interface{}{"content": content}

	title, _ := cmd.Flags().GetString("title")
	if title != "" {
		payload["title"] = title
	}
	sourceURL, _ := cmd.Flags().GetString("url")
	if sourceURL != "" {
		payload["url"] = sourceURL
	}

	body, _ := json.Marshal(payload)

	return cmdutil.RunMutatingAPICommand(f, opts, &cmdutil.MutatingAPISpec{
		Method: "POST", Path: "/api/intel", Body: body,
	})
}

func newCmdList(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List your intel",
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/intel",
				Normalize: cmdutil.NormalizePassthrough,
			})
		},
	}
}

func newCmdRm(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "rm <id>",
		Short: "Delete intel by ID",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			return cmdutil.RunMutatingAPICommand(f, opts, &cmdutil.MutatingAPISpec{
				Method: "DELETE",
				Path:   "/api/intel/" + url.PathEscape(args[0]),
			})
		},
	}
}

func resolveInput(arg string) (string, error) {
	if arg == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	if len(arg) > 1 && arg[0] == '@' {
		data, err := os.ReadFile(arg[1:])
		if err != nil {
			return "", err
		}
		return string(data), nil
	}
	return arg, nil
}
