// Package registry is the single source of truth for all seer-q commands.
// Used by root.go for registration, schema for introspection, and tests.
package registry

import (
	cmdconfig "github.com/chris-xu0321/Midaz-cli/internal/cmd/config"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/claims"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/decisions"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/doctor"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/health"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/market"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/schema"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/search"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/snapshot"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/sources"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/thread"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/threads"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/topic"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/topics"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/usage"
	"github.com/chris-xu0321/Midaz-cli/internal/cmd/version"
	"github.com/chris-xu0321/Midaz-cli/internal/cmdutil"
	"github.com/spf13/cobra"
)

// ArgDef describes a positional argument.
type ArgDef struct {
	Name     string
	Required bool
}

// FlagDef describes a command flag.
type FlagDef struct {
	Name string
}

// CommandDef describes one CLI command.
type CommandDef struct {
	Name        string
	Description string
	Args        []ArgDef
	Flags       []FlagDef
	Endpoints   []string // informational — for schema display
	NewCmd      func(*cmdutil.Factory) *cobra.Command
}

// Commands is the canonical list of all seer-q commands.
var Commands = []CommandDef{
	{
		Name:        "search",
		Description: "Fuzzy search across topics, threads, assets",
		Args:        []ArgDef{{Name: "query", Required: true}},
		Endpoints:   []string{"GET /api/search?q={query}"},
		NewCmd:      search.NewCmdSearch,
	},
	{
		Name:        "market",
		Description: "Global regime + all topics with thread counts",
		Endpoints:   []string{"GET /api/market"},
		NewCmd:      market.NewCmdMarket,
	},
	{
		Name:        "topics",
		Description: "List all topics with thread counts",
		Endpoints:   []string{"GET /api/topics"},
		NewCmd:      topics.NewCmdTopics,
	},
	{
		Name:        "topic",
		Description: "Topic detail + threads",
		Args:        []ArgDef{{Name: "id", Required: true}},
		Endpoints:   []string{"GET /api/topics/{id}"},
		NewCmd:      topic.NewCmdTopic,
	},
	{
		Name:        "threads",
		Description: "List threads",
		Flags:       []FlagDef{{Name: "topic"}, {Name: "status"}},
		Endpoints:   []string{"GET /api/threads"},
		NewCmd:      threads.NewCmdThreads,
	},
	{
		Name:        "thread",
		Description: "Thread detail + claims + market links",
		Args:        []ArgDef{{Name: "id", Required: true}},
		Endpoints:   []string{"GET /api/threads/{id}"},
		NewCmd:      thread.NewCmdThread,
	},
	{
		Name:        "claims",
		Description: "List claims",
		Flags:       []FlagDef{{Name: "thread"}, {Name: "source"}, {Name: "status"}, {Name: "mode"}},
		Endpoints:   []string{"GET /api/claims"},
		NewCmd:      claims.NewCmdClaims,
	},
	{
		Name:        "sources",
		Description: "List ingested sources",
		Flags:       []FlagDef{{Name: "decision"}, {Name: "tier"}},
		Endpoints:   []string{"GET /api/sources"},
		NewCmd:      sources.NewCmdSources,
	},
	{
		Name:        "snapshot",
		Description: "Global regime snapshot",
		Flags:       []FlagDef{{Name: "history"}, {Name: "limit"}},
		Endpoints:   []string{"GET /api/global/snapshot", "GET /api/global/snapshots"},
		NewCmd:      snapshot.NewCmdSnapshot,
	},
	{
		Name:        "usage",
		Description: "Token usage and cost summary",
		Flags:       []FlagDef{{Name: "since"}},
		Endpoints:   []string{"GET /api/usage"},
		NewCmd:      usage.NewCmdUsage,
	},
	{
		Name:        "decisions",
		Description: "Decision audit log",
		Flags:       []FlagDef{{Name: "stage"}, {Name: "run"}, {Name: "entity-type"}, {Name: "entity-id"}, {Name: "limit"}},
		Endpoints:   []string{"GET /api/decisions", "GET /api/decisions/run/{id}"},
		NewCmd:      decisions.NewCmdDecisions,
	},
	{
		Name:        "health",
		Description: "API health check",
		Endpoints:   []string{"GET /api/health"},
		NewCmd:      health.NewCmdHealth,
	},
	{
		Name:        "version",
		Description: "CLI version info",
		NewCmd:      version.NewCmdVersion,
	},
	{
		Name:        "doctor",
		Description: "Diagnostic checks",
		NewCmd:      doctor.NewCmdDoctor,
	},
	{
		Name:        "config",
		Description: "Configuration management",
		NewCmd:      cmdconfig.NewCmdConfig,
	},
	{
		Name:        "schema",
		Description: "Command contract introspection",
		NewCmd:      schema.NewCmdSchema,
	},
}
