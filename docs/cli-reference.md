# CLI Reference

Last updated: 2026-04-22

## Installation

### One-line install (recommended)

**macOS / Linux:**
```bash
curl -fsSL https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.sh | sh
```

**Windows (PowerShell):**
```powershell
irm https://raw.githubusercontent.com/SparkssL/Midaz-cli/main/install.ps1 | iex
```

This installs the `midaz` binary. Run `midaz skills install --yes` afterwards to register agent skills.

### Release (maintainers)

```bash
git tag vX.Y.Z
git push --tags    # CI runs goreleaser and publishes the GitHub Release
```

Skills are embedded in the binary.

---

## Query CLI (`midaz`)

### Response Format

All commands return JSON envelopes:

**Success** (stdout):
```json
{ "ok": true, "data": <payload>, "meta": { "view_url": "...", "count": N } }
```

**Error** (stderr):
```json
{ "ok": false, "error": { "code": "not_found", "message": "...", "hint": "..." } }
```

Use `--raw` to bypass the envelope and get raw API JSON.

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `json` | Output format: `json` or `pretty` |
| `--raw` | false | Raw API response (no envelope) |
| `--api-url` | from config | Override API base URL |
| `--profile` | current | Auth profile name |

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Internal error |
| 2 | Validation error (bad args/flags) |
| 3 | Config error |
| 4 | Network/timeout error |
| 5 | API error (4xx/5xx) |
| 6 | Auth required — run `midaz auth login` |
| 7 | Subscription required — run `midaz subscription start --yes` |

### Entity Lookup

```bash
midaz search "QUERY"            # Fuzzy search across drivers, theses, assets
midaz driver <ID>               # Driver detail: thread members + asset contributions
midaz thesis <ID>               # Thesis detail: snapshot, claims, market links
```

`thread` / `threads` are deprecated hidden aliases of `thesis` / `theses`. The
`topics` / `topic` commands were removed in 0.7.2 — use `drivers` / `driver`.

### List / Browse

```bash
midaz market                    # Global regime + drivers + thesis memberships (composite)
midaz drivers                   # List active drivers (world-layer objects)
midaz driver-links              # Causal edges between drivers (sphere graph)
midaz theses [--status S]       # List all theses (active/weakening/divided/resolved)
midaz claims                    # Latest 100 claims (--source ID, --status S, --mode M)
midaz sources                   # Latest 100 sources (--decision D, --tier N)
midaz delta                     # Recent claims + theses + drivers (--hours N, default 12)
```

### Snapshots

```bash
midaz snapshot                  # Latest global regime snapshot
```

### Desk

```bash
midaz desk get                  # Desk summary (name, shared, subscription, onboarded)
midaz desk settings             # Owner-only: radar, playbook, telegram (GET /api/desk/settings)
midaz desk view                 # Personal market read (GET /api/desks/<own-slug>/read, subscription-gated)
midaz desk share --on --yes     # Toggle public sharing
midaz desk regenerate --yes     # Rebuild personal desk only (fast)
midaz desk reonboard  --yes     # Resubmit current radar + playbook to trigger a rebuild
midaz desk refresh    --yes     # Full pipeline refresh (market + desk; slower)
midaz desk radar {get,set,add,remove,pin,unpin,pins}
midaz desk playbook {get,set}
midaz desk preferences {get,set}   # e.g. set --language zh-CN --yes
midaz desk telegram {status,connect,disconnect}
```

### Account & Subscription

```bash
midaz auth {login,logout,status,whoami,keys}
midaz onboard {status,generate,complete}
midaz invite redeem <CODE> --yes
midaz subscription {status,start,portal}
midaz intel {list,push,rm}
```

### Assets

```bash
midaz assets list [--tier N] [--bias B]
midaz assets get <ID>                      # Bias direction + driver contributions
midaz assets timeline <ID> [--limit N]     # Per-asset event timeline
midaz klines                               # Assets with kline coverage
midaz klines <ID>                          # Candlestick history + latest for one asset
```

### Usage & Audit

```bash
midaz usage                     # Token usage summary (--since P, default 24h)
midaz usage by-run <RUN_ID>     # Per-pipeline-run token-cost breakdown
midaz decisions                 # Decision log (--stage S, --run ID, --entity-type T, --entity-id I, --limit N)
midaz health                    # API health check
```

### Skills

```bash
midaz skills install --yes                       # Install skills to detected agent directories
midaz skills install --target claude --yes       # Install to Claude Code
midaz skills install --target codex --yes        # Install to Codex
midaz skills install --target all --yes          # Install to all known targets
midaz skills install --dry-run                   # Preview without writing
midaz skills install --yes --force               # Overwrite existing skill files
midaz skills install --yes --skill-dir /path     # Custom directory
```

### Diagnostics

```bash
midaz version                   # CLI version, Go version, OS/arch
midaz doctor                    # Check API connectivity, config, health
midaz schema                    # List all command contracts
midaz schema <command>          # Describe one command's input/output contract
```

### Configuration

```bash
midaz config get <key>          # Get config value
midaz config set <key> <value>  # Set config value (creates file if needed)
midaz config list               # List all config (token masked)
midaz config path               # Show config file path
```

Config precedence: CLI flags > env vars > config file > defaults.

Config file: `%APPDATA%\midaz\config.json` (Windows), `~/.config/midaz/config.json` (Linux), `~/Library/Application Support/midaz/config.json` (macOS).

### Full Contract

See `testdata/golden/` for contract examples (golden JSON files for each command).
