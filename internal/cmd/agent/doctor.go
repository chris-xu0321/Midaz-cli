package agent

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	agentres "github.com/SparkssL/seer-cli/agent"
	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/SparkssL/seer-cli/internal/output"
	"github.com/spf13/cobra"
)

type fileCheck struct {
	File    string `json:"file"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

func newCmdDoctor(f *cmdutil.Factory) *cobra.Command {
	var workspace string

	cmd := &cobra.Command{
		Use:   "doctor",
		Short: "Check agent file status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)

			if workspace == "" {
				var err error
				workspace, err = os.Getwd()
				if err != nil {
					return output.ErrConfig("failed to get working directory: %s", err)
				}
			}

			var checks []fileCheck

			// Check installed files
			for _, t := range targets {
				dest := filepath.Join(workspace, ".claude", t.RelPath)
				relPath := ".claude/" + t.RelPath

				data, err := os.ReadFile(dest)
				if err != nil {
					checks = append(checks, fileCheck{relPath, "fail", "not found"})
					continue
				}
				if bytes.Equal(data, *t.Content) {
					msg := "up to date"
					if v := extractVersion(data); v != "" {
						msg = "v" + v + ", up to date"
					}
					checks = append(checks, fileCheck{relPath, "pass", msg})
				} else {
					msg := "exists but differs from source"
					if v := extractVersion(data); v != "" {
						msg = "v" + v + ", differs from source"
					}
					checks = append(checks, fileCheck{relPath, "warn", msg})
				}
			}

			// Check seer-q binary on PATH
			if path, err := exec.LookPath("seer-q"); err == nil {
				checks = append(checks, fileCheck{"seer-q", "pass", path})
			} else {
				checks = append(checks, fileCheck{"seer-q", "warn", "not found on PATH"})
			}

			result := map[string]any{"checks": checks}
			meta := map[string]any{"workspace": workspace}
			return output.WriteSuccess(opts.Out, result, meta, opts.Format)
		},
	}

	cmd.Flags().StringVar(&workspace, "workspace", "", "Target workspace (default: current directory)")
	return cmd
}

// extractVersion pulls version from YAML frontmatter: "version: X.Y.Z"
var versionRe = regexp.MustCompile(`(?m)^version:\s*(.+)$`)

func extractVersion(data []byte) string {
	m := versionRe.FindSubmatch(data)
	if m != nil {
		return string(bytes.TrimSpace(m[1]))
	}
	return ""
}

// ensure agentres is imported
var _ = agentres.SeerCmdMD
