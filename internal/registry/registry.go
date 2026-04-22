// Package registry is the single source of truth for all midaz commands.
// Used by root.go for registration, schema for introspection, and tests.
package registry

import (
	"github.com/SparkssL/Midaz-cli/internal/cmd/assets"
	"github.com/SparkssL/Midaz-cli/internal/cmd/auth"
	"github.com/SparkssL/Midaz-cli/internal/cmd/claims"
	cmdconfig "github.com/SparkssL/Midaz-cli/internal/cmd/config"
	"github.com/SparkssL/Midaz-cli/internal/cmd/decisions"
	"github.com/SparkssL/Midaz-cli/internal/cmd/delta"
	"github.com/SparkssL/Midaz-cli/internal/cmd/desk"
	"github.com/SparkssL/Midaz-cli/internal/cmd/doctor"
	"github.com/SparkssL/Midaz-cli/internal/cmd/drivers"
	"github.com/SparkssL/Midaz-cli/internal/cmd/health"
	"github.com/SparkssL/Midaz-cli/internal/cmd/intel"
	"github.com/SparkssL/Midaz-cli/internal/cmd/invite"
	"github.com/SparkssL/Midaz-cli/internal/cmd/klines"
	"github.com/SparkssL/Midaz-cli/internal/cmd/market"
	"github.com/SparkssL/Midaz-cli/internal/cmd/onboard"
	"github.com/SparkssL/Midaz-cli/internal/cmd/schema"
	"github.com/SparkssL/Midaz-cli/internal/cmd/search"
	"github.com/SparkssL/Midaz-cli/internal/cmd/skills"
	"github.com/SparkssL/Midaz-cli/internal/cmd/snapshot"
	"github.com/SparkssL/Midaz-cli/internal/cmd/sources"
	"github.com/SparkssL/Midaz-cli/internal/cmd/subscription"
	"github.com/SparkssL/Midaz-cli/internal/cmd/theses"
	"github.com/SparkssL/Midaz-cli/internal/cmd/thesis"
	"github.com/SparkssL/Midaz-cli/internal/cmd/thread"
	"github.com/SparkssL/Midaz-cli/internal/cmd/threads"
	"github.com/SparkssL/Midaz-cli/internal/cmd/usage"
	"github.com/SparkssL/Midaz-cli/internal/cmd/version"
	"github.com/SparkssL/Midaz-cli/internal/cmdutil"
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

// Commands is the canonical list of all midaz commands.
var Commands = []CommandDef{
	// --- Account & auth -------------------------------------------------------
	{
		Name:        "auth",
		Description: "Authenticate and manage API credentials",
		Endpoints:   []string{"POST /api/app/auth/exchange", "GET /api/app/me", "GET/POST/DELETE /api/app/api-keys"},
		NewCmd:      auth.NewCmdAuth,
	},
	{
		Name:        "onboard",
		Description: "Complete desk onboarding (radar + playbook)",
		Endpoints:   []string{"POST /api/desk/onboard", "POST /api/desk/onboard/generate"},
		NewCmd:      onboard.NewCmdOnboard,
	},
	{
		Name:        "invite",
		Description: "Redeem invitation codes",
		Endpoints:   []string{"POST /api/invite/redeem"},
		NewCmd:      invite.NewCmdInvite,
	},
	{
		Name:        "subscription",
		Description: "Subscribe, manage billing, check subscription status",
		Endpoints:   []string{"POST /api/stripe/checkout", "POST /api/stripe/portal", "GET /api/desk"},
		NewCmd:      subscription.NewCmdSubscription,
	},
	// --- Desk -----------------------------------------------------------------
	{
		Name:        "desk",
		Description: "Manage your desk: radar, playbook, preferences, sharing, Telegram",
		Endpoints: []string{
			"GET /api/desk", "GET /api/desk/settings", "GET /api/desk/view",
			"PATCH /api/desk*", "PATCH /api/desk/preferences",
			"DELETE /api/desk/telegram",
			"POST /api/desk/radar/pin", "DELETE /api/desk/radar/pin", "GET /api/desk/radar/pins",
			"POST /api/desk/personal-desk/regenerate", "POST /api/desk/refresh",
			"POST /api/desk/onboard",
		},
		NewCmd: desk.NewCmdDesk,
	},
	{
		Name:        "intel",
		Description: "Push, list, and delete private intel (notes / sources)",
		Endpoints:   []string{"GET /api/intel", "POST /api/intel", "DELETE /api/intel/{id}"},
		NewCmd:      intel.NewCmdIntel,
	},
	// --- Market read ----------------------------------------------------------
	{
		Name:        "search",
		Description: "Fuzzy search across drivers, theses, assets",
		Args:        []ArgDef{{Name: "query", Required: true}},
		Endpoints:   []string{"GET /api/search?q={query}"},
		NewCmd:      search.NewCmdSearch,
	},
	{
		Name:        "market",
		Description: "Global regime + drivers + thesis memberships (composite)",
		Endpoints:   []string{"GET /api/market"},
		NewCmd:      market.NewCmdMarket,
	},
	{
		Name:        "drivers",
		Description: "List active drivers (world-layer objects behind the market view)",
		Endpoints:   []string{"GET /api/drivers"},
		NewCmd:      drivers.NewCmdDrivers,
	},
	{
		Name:        "driver",
		Description: "Driver detail with thread members and asset contributions",
		Args:        []ArgDef{{Name: "id", Required: true}},
		Endpoints:   []string{"GET /api/drivers/{id}"},
		NewCmd:      drivers.NewCmdDriver,
	},
	{
		Name:        "driver-links",
		Description: "List causal edges between drivers (sphere/radar graph)",
		Endpoints:   []string{"GET /api/driver-links"},
		NewCmd:      drivers.NewCmdDriverLinks,
	},
	{
		Name:        "theses",
		Description: "List theses (market arguments)",
		Flags:       []FlagDef{{Name: "status"}},
		Endpoints:   []string{"GET /api/theses"},
		NewCmd:      theses.NewCmdTheses,
	},
	{
		Name:        "thesis",
		Description: "Thesis detail + claims + market links",
		Args:        []ArgDef{{Name: "id", Required: true}},
		Endpoints:   []string{"GET /api/theses/{id}"},
		NewCmd:      thesis.NewCmdThesis,
	},
	{
		Name:        "threads",
		Description: "Deprecated alias of `theses`",
		Flags:       []FlagDef{{Name: "status"}},
		Endpoints:   []string{"GET /api/theses"},
		NewCmd:      threads.NewCmdThreads,
	},
	{
		Name:        "thread",
		Description: "Deprecated alias of `thesis`",
		Args:        []ArgDef{{Name: "id", Required: true}},
		Endpoints:   []string{"GET /api/theses/{id}"},
		NewCmd:      thread.NewCmdThread,
	},
	{
		Name:        "assets",
		Description: "List and inspect assets (bias, contributions, timeline)",
		Endpoints:   []string{"GET /api/assets", "GET /api/assets/{id}", "GET /api/assets/{id}/timeline"},
		NewCmd:      assets.NewCmdAssets,
	},
	{
		Name:        "klines",
		Description: "Asset kline (candlestick) data",
		Args:        []ArgDef{{Name: "asset_id", Required: false}},
		Endpoints:   []string{"GET /api/klines", "GET /api/klines/{asset_id}"},
		NewCmd:      klines.NewCmdKlines,
	},
	{
		Name:        "delta",
		Description: "Recent claims + theses + drivers from the last N hours",
		Flags:       []FlagDef{{Name: "hours"}},
		Endpoints:   []string{"GET /api/delta"},
		NewCmd:      delta.NewCmdDelta,
	},
	{
		Name:        "claims",
		Description: "List claims",
		Flags:       []FlagDef{{Name: "source"}, {Name: "status"}, {Name: "mode"}},
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
		Endpoints:   []string{"GET /api/global"},
		NewCmd:      snapshot.NewCmdSnapshot,
	},
	{
		Name:        "usage",
		Description: "Token usage and cost summary",
		Flags:       []FlagDef{{Name: "since"}},
		Endpoints:   []string{"GET /api/usage", "GET /api/usage/by-run/{run_id}"},
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
	{
		Name:        "skills",
		Description: "Manage embedded agent skills (install to Claude Code, Codex, etc.)",
		NewCmd:      skills.NewCmdSkills,
	},
}
