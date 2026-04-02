package schema

import (
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
	"github.com/SparkssL/Midaz-cli/internal/output"
	"github.com/spf13/cobra"
)

// CommandInfo holds the data schema needs from the registry.
// Avoids importing registry directly (which would cause a cycle).
type CommandInfo struct {
	Name        string
	Description string
	Args        []string
	Flags       []string
	Endpoints   []string
}

// SchemaData is set by the registry package after initialization.
var SchemaData []CommandInfo

func NewCmdSchema(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "schema [command]",
		Short: "Describe command contracts",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts := cmdutil.ResolveRunOpts(cmd, f)
			if len(args) == 0 {
				return listAll(opts)
			}
			return describeOne(opts, args[0])
		},
	}
}

func listAll(opts *cmdutil.RunOpts) error {
	var commands []map[string]any
	for _, info := range SchemaData {
		entry := map[string]any{
			"name":        info.Name,
			"description": info.Description,
		}
		if len(info.Args) > 0 {
			entry["args"] = info.Args
		}
		if len(info.Flags) > 0 {
			entry["flags"] = info.Flags
		}
		commands = append(commands, entry)
	}
	data := map[string]any{"commands": commands}
	meta := map[string]any{"count": len(commands)}
	return output.WriteSuccess(opts.Out, data, meta, opts.Format)
}

func describeOne(opts *cmdutil.RunOpts, name string) error {
	for _, info := range SchemaData {
		if info.Name == name {
			data := map[string]any{
				"name":        info.Name,
				"description": info.Description,
			}
			if len(info.Args) > 0 {
				args := make([]map[string]any, len(info.Args))
				for i, a := range info.Args {
					args[i] = map[string]any{"name": a, "required": true, "type": "string"}
				}
				data["args"] = args
			} else {
				data["args"] = []any{}
			}
			if len(info.Flags) > 0 {
				data["flags"] = info.Flags
			} else {
				data["flags"] = []any{}
			}
			if len(info.Endpoints) == 1 {
				data["api_endpoint"] = info.Endpoints[0]
			} else if len(info.Endpoints) > 1 {
				data["api_endpoints"] = info.Endpoints
			}
			return output.WriteSuccess(opts.Out, data, nil, opts.Format)
		}
	}
	return output.ErrWithHint(output.ExitValidation, "validation",
		"Unknown command: "+name, "run: seer-q schema")
}
