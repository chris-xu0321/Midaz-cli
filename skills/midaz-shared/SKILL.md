---
name: midaz-shared
version: 0.7.2
description: Midaz CLI shared concepts — auth model, response format, global flags, and safety rules that apply to every midaz skill
metadata: {"requires":{"bins":["midaz"]}}
---

# Midaz Shared

Foundational knowledge for every `midaz-*` skill. Read this before using the other skills.

## What Midaz Is

Midaz is an **Interactive Cognitive Trading Map** plus the intelligence pipeline behind it. From the CLI an agent can mirror everything a user can do in the web app at https://www.midaz.xyz:

- **Account** — sign in, redeem invite codes, onboard, subscribe, manage API keys.
- **Desk** — edit the radar (watchlist) and playbook (trading rules), set language preference, toggle public sharing, connect Telegram for alerts, push private intel.
- **Market** — browse drivers, theses, claims, assets, klines, global snapshots, deltas.

Everything is JSON-in / JSON-out.

## Command Surface At A Glance

```
midaz auth {login,logout,status,whoami,keys}
midaz onboard {status,generate,complete}
midaz invite redeem <CODE>
midaz subscription {status,start,portal}
midaz desk {get,settings,view,share,regenerate,reonboard,refresh,radar,playbook,preferences,telegram}
midaz intel {list,push,rm}
midaz assets {list,get,timeline}
midaz klines [asset_id]
midaz delta
midaz search | market | drivers | driver | driver-links | theses | thesis | claims | sources | snapshot | usage | decisions | health
midaz usage by-run <run_id>
midaz doctor | version | config | schema | skills
```

If the user's intent isn't obviously covered, run `midaz schema` to list every registered command, then `midaz schema <command>` for the contract.

## Response Format

All commands return JSON:

- **Success** (stdout): `{ "ok": true, "data": <payload>, "meta": { "view_url": "...", "count": N, ... } }`
- **Errors** (stderr): `{ "ok": false, "error": { "code": "...", "message": "...", "hint": "..." } }`

Access the payload via `.data`. Page-level `view_url` lives in `.meta.view_url`. Per-entity URLs appear as `view_url` fields inside `.data` objects.

Pass `--raw` to bypass the envelope — useful when exploring unfamiliar endpoints.

### Exit codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | Internal CLI error |
| 2 | Validation (missing flag/arg, bad input) |
| 3 | Config error |
| 4 | Network / timeout |
| 5 | API 4xx/5xx (non-auth) |
| 6 | Auth required or invalid — run `midaz auth login` |
| 7 | Subscription required — run `midaz subscription start --yes` |

## Global Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--format` | `json` | `json` or `pretty` |
| `--raw` | `false` | Raw API response (no envelope) |
| `--api-url` | from config | Override API base URL |
| `--profile` | current | Auth profile name (set by last `auth login`) |

## Auth Model

Midaz uses **personal access tokens (PATs)** — strings prefixed `sk_`. Flow:

1. `midaz auth login` opens the browser; user signs in at midaz.xyz.
2. The web app calls `POST /api/app/auth/exchange` and POSTs the PAT to a local loopback server the CLI started.
3. The CLI stores the PAT at `~/.config/midaz/auth.json` (mode `0600`).

For CI / headless contexts:

- `midaz auth login --paste` prompts for a PAT created on the website.
- `midaz auth login --token sk_…` stores a PAT inline.
- `MIDAZ_TOKEN=sk_…` env var overrides the stored PAT.

If a command returns exit code 6, the PAT is missing or expired — run `midaz auth login` again.

## Side-Effect Gating

Every command that mutates server state **requires `--yes`**. Prevents an agent from silently sending messages, spending money, or changing the user's desk. Examples:

```
midaz desk radar set --items "Fed,AI,Oil" --yes
midaz subscription start --yes
midaz invite redeem ABC-123 --yes
midaz intel push --from-file note.md --yes
midaz desk telegram disconnect --yes
```

Read-only commands never require `--yes`.

## Browser Handoffs

Commands that need the user to complete a flow in the browser (Stripe checkout, Stripe billing portal, Telegram bot deep-link, login) print the URL in the response envelope and best-effort auto-open it. In headless / SSH / CI contexts the CLI skips auto-open; surface the URL to the user as a clickable markdown link.

Env var `MIDAZ_NO_BROWSER=1` disables auto-open globally.

## Config & Diagnostics

```
midaz version                    # CLI version
midaz doctor                     # Connectivity + auth check
midaz config get <key>           # Show config value
midaz config set <key> <value>   # Update config
midaz config list                # Dump config
midaz config path                # Show config file path
midaz schema                     # List every registered command
midaz schema <command>           # Describe a command's contract
midaz health                     # API health check
midaz setup all --yes            # (Re)install skills to agent directories
```

Config path: `~/.config/midaz/config.json` (Linux/macOS) or `%APPDATA%\midaz\config.json` (Windows).

## Common Rules

1. **Run `midaz auth status` at the top of a session** when you're about to call authenticated commands — it returns the user's id, desk id, subscription status, and onboarding state in one call, so you can decide what to do.
2. **Never drop `--yes` on write commands.** If a command exits 2 with `confirmation_required`, add the flag — don't retry blindly.
3. **Always surface `view_url` as clickable markdown links.** Page-level in `.meta.view_url`; per-entity in each object under `.data`. Format as `[descriptive text](url)` — never paste raw URLs.
4. **Synthesize, don't dump.** Convert JSON into natural language before replying.
5. **Respect exit code 7.** On subscription-required, ask the user before running `subscription start` unless they already asked to subscribe.
6. **Thread → thesis.** The product now says "thesis". Use `midaz theses` / `midaz thesis <id>`. `threads` / `thread` still work but are deprecated.
7. **Topic → driver.** The old topic layer was replaced by the driver ontology. Use `midaz drivers` / `midaz driver <id>` / `midaz driver-links`. The `topics` / `topic` commands are gone.
