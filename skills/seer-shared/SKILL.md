---
name: seer-shared
version: 0.2.0
description: Seer CLI shared concepts, response format, config, and common rules
metadata: {"requires":{"bins":["seer-q"]}}
---

# Seer Shared

Foundational knowledge for all Seer skills. Read this before using any domain skill.

## What Seer Is

Seer is a market thesis intelligence system. It tracks:

- **Claims** — atomic evidence statements extracted from sources (the raw events/facts)
- **Threads** — tradable angles / sub-theses built from clusters of claims
- **Topics** — narrative domains that group related threads (e.g., "AI Infrastructure", "Energy Transition")
- **Global snapshot** — overall market regime derived from top drivers across all topics

All queries use the `seer-q` CLI.

## Response Format

All commands return JSON:
- **Success** (stdout): `{ "ok": true, "data": <payload>, "meta": { "view_url": "...", "count": N } }`
- **Errors** (stderr): `{ "ok": false, "error": { "code": "...", "message": "..." } }`

Access the payload via `.data`. The `meta` field contains `view_url` and count hints.
For raw API output (no envelope), use `--raw`.

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `json` | Output format: `json` or `pretty` |
| `--raw` | false | Raw API response (no envelope) |
| `--api-url` | from config | Override API base URL |

## Config & Diagnostics

```bash
seer-q version                  # CLI version, Go version, OS/arch
seer-q doctor                   # Verify API connectivity and config
seer-q config get <key>         # Get config value
seer-q config set <key> <value> # Set config value
seer-q config list              # Show active configuration
seer-q config path              # Show config file path
seer-q schema                   # List all command contracts
seer-q schema <command>         # Describe a command's input/output contract
seer-q health                   # API health check
```

## Common Rules

1. Use `seer-q search` first whenever the user mentions a specific entity, asset, or theme
2. ALWAYS include the `view_url` from `.data` or `.meta` as a link for the user
3. Synthesize data into natural language — don't dump raw JSON
4. For multi-entity questions, make multiple calls to build a complete picture
5. When claims are asked about, note their `claim_mode`, `thread_role` (support/contradiction), and `event_date`
