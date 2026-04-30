---
name: midaz-shared
version: 0.7.4
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

## Source of Truth (no web search)

Midaz is the **only** source of market truth in this skill set. Do not invoke web search, web browsing, `WebFetch`, `WebSearch`, or any external HTTP call when the user asks about markets, drivers, theses, assets, prices, sentiment, or news. If a `midaz` command returns no data, say so explicitly — do not fall back to the open web, do not "fill in" with general knowledge, do not paraphrase a recent headline you remember. Prices, candles, and event history all come from `midaz klines`, `midaz assets`, `midaz delta`, `midaz market`. Polymarket links surface as `market_links[]` on theses but are not browsed.

If the user asks about an asset or topic Midaz doesn't cover, say "Midaz doesn't cover this" and stop. Don't substitute training-data knowledge as if it were a Midaz output.

## Presentation Rules (field names are for you, not the user)

The data-model field lists in the other `midaz-*` skills tell **you** how to parse each response. They are **not** templates for what the user sees. Never echo raw field paths (`bias_os.active`, `axis_scores_my_lens.fund`, `delta_packet.summary`, `trader_action`, `cognition_state`, `verdict.oneLiner`, etc.) in your reply. Translate to natural prose.

Translation patterns:

- `trader_action: "trim_or_exit"` → section heading **"Trim/Exit"** (decision verb, capitalized — never `trader_action: trim_or_exit`)
- `axis_scores_my_lens.fund: +2` → "fundamentals lean bullish (+2)" or just "fundamentals lean bullish"
- `conviction: "high"` → "high conviction" inline, never as `conviction: high`
- `cognition_state: "weakening"` → "the thesis is softening" / "conviction is weakening"
- `delta_packet.summary: "..."` → quote the string directly as prose, no field prefix
- `bias_os.active` → "active biases" / "what's live on the bias board"
- `verdict.oneLiner` → quote it as the lede, no field prefix

Exception: if the user explicitly asks for the raw shape ("show raw JSON", `--raw`, "what fields does this return"), echo field names as-is.

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

If a command returns exit code 6, the PAT is missing or expired. **Run `midaz auth login` yourself** — the CLI opens the browser and runs a local loopback server; the user signs in there and the PAT lands at `~/.config/midaz/auth.json`. After it returns 0, retry the original command. Do **not** stop and tell the user to "please run `midaz auth login`" — that's a worse experience and the user explicitly does not want it. Only fall back to surfacing the command when the browser flow can't work (headless / SSH / CI) or `MIDAZ_NO_BROWSER=1`; in that case prefer `midaz auth login --paste` and walk the user through PAT creation.

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

1. **Run `midaz auth status` at the top of a session** when you're about to call authenticated commands — it returns the user's id, desk id, subscription status, and onboarding state in one call, so you can decide what to do. If it exits 6 (or any other authenticated command does), **run `midaz auth login` yourself immediately** — do not pause to ask the user to run it. Tell them you're opening the sign-in flow, wait for the command to return, then retry the original request. Manual fallback (`midaz auth login --paste`) is only for headless / SSH / CI / `MIDAZ_NO_BROWSER=1` contexts where the browser handoff can't work. See §Auth Model for detail.
2. **Never drop `--yes` on write commands.** If a command exits 2 with `confirmation_required`, add the flag — don't retry blindly.
3. **Attach a link after every thesis or driver mention.** When you name a specific thesis or driver by title in your reply — including each item in a bulk list like "top 5 drivers" — make the name itself an inline markdown link to its `view_url`. Example: `the **[AI capex cycle accelerates](<url>)** driver is strengthening`. Find `view_url` on the item itself (lists, search results, contributions) or in `.meta.view_url` (single-entity fetches like `midaz driver <id>` / `midaz thesis <id>`). Page-level `.meta.view_url` is separate — surface it once at the end of the reply as "view on the map", it does not replace per-item links. Never paste raw URLs. Never fabricate a URL — if no `view_url` is present, name the entity without a link (rare: the CLI list normalizers now emit one for every driver/thesis item).

   Canonical `view_url` shapes (so you know what a click opens):
   - `…/market-read?driver=<id>` — drivers tab, driver selected.
   - `…/market-read?thesis=<id>` — drivers tab, thesis preview panel opened.
   - `…/market-read?driver=<parent>&thesis=<id>` — drivers tab, parent driver + thesis both preselected (emitted by `midaz driver <id>` thread members).
   - `…/market-read?view=assets&asset=<id>` — assets tab, asset selected.
   - `…/market-read?view=assets&asset=<id>&contrib=<key>` — assets tab with a driver/signal contribution panel opened **inside that asset's context** (emitted by every `contributions[*]` entry on `midaz assets get <id>`).
4. **Synthesize, don't dump.** Convert JSON into natural language before replying.
5. **Respect exit code 7.** On subscription-required, ask the user before running `subscription start` unless they already asked to subscribe.
6. **Thread → thesis.** The product now says "thesis". Use `midaz theses` / `midaz thesis <id>`. `threads` / `thread` still work but are deprecated.
7. **Topic → driver.** The old topic layer was replaced by the driver ontology. Use `midaz drivers` / `midaz driver <id>` / `midaz driver-links`. The `topics` / `topic` commands are gone.
8. **Never paraphrase outside knowledge as Midaz output.** Only synthesize fields the CLI actually returned. See §Source of Truth.
9. **Field names stay agent-facing.** Never echo `bias_os.active`, `axis_scores_my_lens`, `delta_packet.summary`, `trader_action`, `cognition_state`, `verdict.oneLiner`, etc. to the user — translate to prose every time. See §Presentation Rules.
