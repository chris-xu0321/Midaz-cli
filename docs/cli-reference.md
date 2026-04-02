# CLI Reference

Last updated: 2026-03-30

## Installation

### Step 1: Install CLI

```bash
npm install -g @midaz/cli
```

Requires Node.js >= 16. Supported platforms: Windows, macOS, Linux (x64, arm64).

### Step 2: Install Skills

Install all skills:

```bash
npx skills add SparkssL/Midaz-cli -y -g
```

Skills are installed directly from this repo via `npx skills add SparkssL/Midaz-cli -y -g`.

### Release (maintainers)

```bash
bash npm/publish.sh              # goreleaser + npm publish (single package)
bash npm/publish.sh --dry-run    # test without publishing
```

Skills are installed directly from this repo via `npx skills add SparkssL/Midaz-cli -y -g`. No separate skills publish step is required.

---

## Query CLI (`seer-q`)

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

### Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Internal error |
| 2 | Validation error (bad args/flags) |
| 3 | Config error |
| 4 | Network/timeout error |
| 5 | API error (4xx/5xx) |

### Entity Lookup

```bash
seer-q search "QUERY"           # Fuzzy search across topics, threads, assets
seer-q topic <ID>               # Topic detail: thesis, bias, all threads
seer-q thread <ID>              # Thread detail: snapshot, claims, market links
```

### List / Browse

```bash
seer-q market                   # Global regime + all topics with thread counts
seer-q topics                   # List all topics with thread counts
seer-q threads                  # List all threads (--topic ID, --status S)
seer-q claims                   # Latest 100 claims (--thread ID, --source ID, --status S, --mode M)
seer-q sources                  # Latest 100 sources (--decision D, --tier N)
```

### Snapshots

```bash
seer-q snapshot                 # Latest global regime snapshot
seer-q snapshot --history       # Regime snapshot history (--limit N, default 10)
```

### Usage & Audit

```bash
seer-q usage                    # Token usage summary (--since P, default 24h)
seer-q decisions                # Decision log (--stage S, --run ID, --entity-type T, --entity-id I, --limit N)
seer-q health                   # API health check
```

### Diagnostics

```bash
seer-q version                  # CLI version, Go version, OS/arch
seer-q doctor                   # Check API connectivity, config, health
seer-q schema                   # List all command contracts
seer-q schema <command>         # Describe one command's input/output contract
```

### Configuration

```bash
seer-q config get <key>         # Get config value
seer-q config set <key> <value> # Set config value (creates file if needed)
seer-q config list              # List all config (token masked)
seer-q config path              # Show config file path
```

Config precedence: CLI flags > env vars > config file > defaults.

Config file: `%APPDATA%\seer\config.json` (Windows), `~/.config/seer/config.json` (Linux), `~/Library/Application Support/seer/config.json` (macOS).

### Full Contract

See `testdata/golden/` for contract examples (golden JSON files for each command).
