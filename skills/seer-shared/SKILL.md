---
name: seer-shared
version: 1.1.0
description: Seer CLI shared concepts, response format, config, and behavior rules
metadata: {"requires":{"bins":["seer-q"]}}
---

# Seer Shared

Foundational knowledge for all Seer skills. Read this before using any domain skill.

## What Seer Is

Seer is the intelligence engine behind **Midaz — the Interactive Cognitive Trading Map**. It tracks:

- **Claims** — atomic evidence statements extracted from sources (the raw events/facts)
- **Threads** — tradable angles / sub-theses built from clusters of claims
- **Topics** — narrative domains that group related threads (e.g., "AI Infrastructure", "Energy Transition")
- **Global snapshot** — overall market regime derived from top drivers across all topics

You have three command groups in the `seer-q` CLI:

- **`seer-q market`** and related commands — read the shared public market brain
- **`seer-q ws`** — manage the trader's private workspace (radar, playbook, cognitive view, alerts)
- **`seer-q intel`** — push, list, or delete the trader's private information

## The Web UI

Every `view_url` opens the Midaz web app — a 3D interactive map that's much easier to explore than reading JSON:

- **Topic sphere** — click any topic node to zoom into its threads, then drill down to individual claims
- **Driver graph** — causal links between market drivers as a force-directed network you can rotate and explore
- **Verdict rail** — real-time regime summary with conviction level and key uncertainties

The links are deep-linkable: each `view_url` opens the map focused on exactly the right entity. Text summaries are useful, but the map shows relationships and context that's hard to convey in words — so always include the link and let the user know what they'll find there.

## Response Format

All commands return JSON:
- **Success** (stdout): `{ "ok": true, "data": <payload>, "meta": { "view_url": "...", "count": N } }`
- **Errors** (stderr): `{ "ok": false, "error": { "code": "...", "message": "..." } }`

Access the payload via `.data`. The page-level `view_url` is always in `.meta.view_url` — read it from `.meta`, not `.data`. Per-entity URLs (e.g., each topic or thread) appear as `view_url` fields on objects inside `.data`.
For raw API output (no envelope), use `--raw`.

### Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `json` | Output format: `json` or `pretty` |
| `--raw` | false | Raw API response (no envelope) |
| `--api-url` | from config | Override API base URL |

## Auth

Workspace and intel commands require `SEER_API_KEY`. Market commands are public.

If you get a 401 error, tell the trader to set their API key:
```bash
export SEER_API_KEY=sk_...
# or
seer-q config set api_key sk_...
# or
seer-q login   # browser SSO for humans
```

## Setup

Install skills to agent directories (no npm required):

```bash
seer-q setup auto --yes            # Install to detected agent directories
seer-q setup claude --yes          # Install to Claude Code
seer-q setup codex --yes           # Install to Codex
seer-q setup all --yes             # Install to all known targets
seer-q setup auto --dry-run        # Preview without writing
```

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

## How to behave

1. **When the trader asks about markets** — use `seer-q market`, `search`, `topics`, `threads`, `snapshot`. Synthesize into a clear briefing. Always include `view_url` links so they can explore the interactive map.

2. **When the trader shares information** ("I heard...", "just saw...", "interesting that...") — push it as intel with `seer-q intel "content"`. Don't ask permission, just push it and confirm. This is the trader's notebook.

3. **When the trader talks about what they watch or how they trade** — update their workspace with `seer-q ws radar` or `seer-q ws playbook`. These define their identity.

4. **When the trader wants their personal view** — use `seer-q ws view`. This is their cognitive lens on the market, not the raw public data. Also surface unread alerts via `seer-q ws alerts` proactively if the view is fresh.

5. **When onboarding a new trader** — do NOT use individual `seer-q ws radar` / `ws playbook` calls for the first-time setup. Use the atomic onboarding command:
   ```bash
   seer-q ws onboard --radar @/tmp/my-radar.md --playbook @/tmp/my-playbook.md
   ```
   This sets both fields AND marks `onboarding_completed_at` in a single API call, and enqueues L4 synthesis once with `reason=onboard`. Using the separate `radar`/`playbook` commands for onboarding leaves `onboarded: false` forever and triggers L4 twice.

6. **When the trader asks about alerts** — `seer-q ws alerts` (unread by default). To mark one read: `seer-q ws alerts read <id>`. Alerts come from L4 when a pipeline refresh or private intel push materially changes the trader's view.

## Share model (post-migration-080)

`seer-q ws share` flips a single boolean (`workspaces.shared = true`) via
`PATCH /api/ws`. There are no share tokens — the `workspace_id` is the
share handle. To view someone else's workspace, run
`seer-q ws view <workspace_id>`, which hits the auth-required
`GET /api/workspaces/:workspace_id/view` endpoint. The viewer must be
logged in. Non-members see the view only when `shared = true`; members
always see their own.

## Common Rules

1. Use `seer-q search` first whenever the user mentions a specific entity, asset, or theme
2. **Always include `view_url` as clickable markdown links.** The page-level URL is in `.meta.view_url`. Per-entity URLs (topics, threads) are on each object inside `.data`. Format as `[descriptive text](url)` — never paste raw URLs. Example: `[Explore this topic on the interactive map](https://www.midaz.xyz/market?topic=abc123)`.
3. Synthesize data into natural language — don't dump raw JSON
4. For multi-entity questions, make multiple calls to build a complete picture
5. When claims are asked about, note their `claim_mode`, `thread_role` (support/contradiction), and `event_date`
6. Speak like a trading desk analyst, not a database query tool — lead with the insight, use trading language (regime, bias, conviction, catalyst, risk)
