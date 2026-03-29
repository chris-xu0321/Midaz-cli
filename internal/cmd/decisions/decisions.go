package decisions

import (
	"net/url"
	"strconv"

	"github.com/SparkssL/seer-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdDecisions(f *cmdutil.Factory) *cobra.Command {
	var stage, runID, entityType, entityID string
	var limit int

	cmd := &cobra.Command{
		Use:   "decisions",
		Short: "Decision audit log",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)

			// Routing: if --run set and no other filters → dedicated sub-route
			if runID != "" && stage == "" && entityType == "" && entityID == "" {
				return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
					Path:      "/api/decisions/run/" + url.PathEscape(runID),
					Normalize: cmdutil.NormalizeBareArray,
				})
			}

			params := url.Values{}
			if stage != "" {
				params.Set("stage", stage)
			}
			if entityType != "" {
				params.Set("entity_type", entityType)
			}
			if entityID != "" {
				params.Set("entity_id", entityID)
			}
			if runID != "" {
				params.Set("pipeline_run_id", runID)
			}
			if limit > 0 {
				params.Set("limit", strconv.Itoa(limit))
			}

			return cmdutil.RunAPICommand(f, opts, &cmdutil.APISpec{
				Path:      "/api/decisions",
				Params:    params,
				Normalize: cmdutil.NormalizeBareArray,
			})
		},
	}

	cmd.Flags().StringVar(&stage, "stage", "", "Filter by stage")
	cmd.Flags().StringVar(&runID, "run", "", "Filter by pipeline run ID")
	cmd.Flags().StringVar(&entityType, "entity-type", "", "Filter by entity type")
	cmd.Flags().StringVar(&entityID, "entity-id", "", "Filter by entity ID")
	cmd.Flags().IntVar(&limit, "limit", 0, "Max results (default 50)")
	return cmd
}
