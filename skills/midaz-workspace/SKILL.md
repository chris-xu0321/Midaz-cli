---
name: midaz-workspace
version: 0.6.1
description: Manage the Midaz workspace — radar, playbook, sharing, Telegram alerts, private intel, and asset tracking via the CLI
metadata: {"requires":{"bins":["midaz"]}}
---

# Midaz Workspace

> Read [midaz-shared](../midaz-shared/SKILL.md) for auth/envelope basics and [midaz-account](../midaz-account/SKILL.md) for signin/onboarding.

Everything a signed-in user can do inside their workspace at `/workspace/view` and `/workspace/settings`, exposed as CLI commands. All write commands require `--yes`.

## Workspace Read Commands

```
midaz workspace get          # Summary (name, shared flag, subscription, has_invite_access, onboarded)
midaz workspace settings     # Owner-only: radar, playbook, telegram status  (GET /api/ws/settings)
midaz workspace view         # Personal market view — subscription-gated    (GET /api/ws/view)
```

`workspace get` is the cheapest way to inspect state at the start of a session.

## Radar (watchlist)

The radar is a short list of domains / assets / events the user wants Midaz to focus on. Rules:

- ≤5 items
- Each item ≤200 chars

```
midaz workspace radar get
midaz workspace radar set --items "Fed policy, AI capex, Oil, China CNY" --yes
midaz workspace radar set --from-file radar.md --yes
```

Updates enqueue an L4 refresh; `l4_enqueued: true` in the response means the personal view will recompute soon.

### Entity pins (radar pin / unpin / pins)

Distinct from free-form radar lines: pins attach a specific **entity** (thesis, topic, driver, asset) to the radar with provenance tracking, so the web market view can render a filled pin button and L4 can treat the entity as a first-class watch target.

```
midaz workspace radar pin --kind Thesis --source-type thread --source-id <id> --label "Short text"  --yes
midaz workspace radar pin --kind Topic  --source-type topic  --source-id <id> --label "Short text"  --yes
midaz workspace radar pin --kind Driver --source-type driver --source-id <id> --label "Short text"  --yes
midaz workspace radar pin --kind Asset  --source-type asset  --source-id <id> --label "Short text"  --yes

midaz workspace radar unpin --source-type thread --source-id <id> --yes

midaz workspace radar pins                 # List pins (provenance rows) — includes origin = pin | adopted
```

Rules enforced server-side:
- `kind` is the display label (`Thesis | Topic | Driver | Asset`); `source-type` is the DB key (`thread | topic | driver | asset`).
- `label` ≤ 160 chars after whitespace collapse.
- 409 `already_pinned` if the (source-type, source-id) is already pinned.
- 409 `radar_full` if the radar already has 12 lines *and* no freeform line can be adopted.
- If a freeform radar line matches the `label`, the pin adopts it (`origin: adopted`); otherwise a new line is appended (`origin: pin`).
- Unpin is origin-aware: `origin=pin` strips the line; `origin=adopted` only removes provenance (the freeform line survives); unknown pairs no-op.

Pin/unpin enqueue an L4 refresh (`l4_enqueued: true` when work was queued).

## Playbook (trading rules)

Markdown, ≤20 000 chars. Describes how the user wants Midaz to interpret the market.

```
midaz workspace playbook get
midaz workspace playbook set --from-file playbook.md --yes
```

## Sharing

Flip the `shared` boolean to expose a read-only workspace page at `https://www.midaz.xyz/w/<workspace_id>`:

```
midaz workspace share --on --yes
midaz workspace share --off --yes
```

After enabling, `midaz workspace get` returns `workspace.shared: true`. The public URL is computed client-side from the workspace id.

## Telegram Alerts

Connect the workspace to the Midaz Telegram bot so alerts are delivered to chat.

```
midaz workspace telegram status       # Polls GET /api/ws/settings → telegram.{connected,bot_username}
midaz workspace telegram connect      # Prints + opens https://t.me/<bot>?start=<workspace_id>
midaz workspace telegram disconnect --yes
```

Flow:
1. Run `telegram connect`. The envelope has a `view_url` pointing at the Telegram deep link; the CLI also attempts to auto-open it. Instruct the user to tap "Start" inside Telegram.
2. Poll `midaz workspace telegram status`. `telegram.connected` flips to `true` once the bot webhook has stored the chat id.
3. The workspace will then receive alerts as they're produced by the L4 pipeline.

## Intel (private notes / sources)

Push private research into the Midaz intel store. Subscription-gated (exit 7 if no trial/active subscription).

```
midaz intel list [--limit 50]
midaz intel push --from-file note.md --title "My note" --url https://… --published-at 2026-04-15T00:00:00Z --yes
midaz intel rm <id> --yes
```

Limits:
- `content` ≤ 100 000 chars
- `--url` and `--published-at` are optional
- `--title` defaults to the first line of the file

Intel items feed into L4 like any other source; a push or delete enqueues a refresh.

## Assets

Read-only; no auth required, but the richest context appears when paired with a workspace.

```
midaz assets list [--tier 1|2] [--bias bullish|bearish|neutral|mixed]
midaz assets get <asset_id>
midaz assets thesis <asset_id> <thesis_id>
```

Fields of interest:

- `bias`, `bias_score` — aggregate stance across linked theses
- `thesis_count`, `bull_count`, `bear_count`, `mixed_count`
- `links[]` — each has `thesis_id`, `stance`, `weight`, `rationale`, and `evidence[]` with claim snippets
- `view_url` — deep link into the map

## Agent Recipes

### "What's in my workspace right now?"

1. `midaz workspace get` — prints subscription, onboarded, has_invite_access.
2. If `has_invite_access: false` → point them to `midaz invite redeem`.
3. If `onboarded: false` → point them to `midaz onboard`.
4. If `subscription.allowed: false` → point them to `midaz subscription start`.
5. Otherwise: `midaz workspace view` to show the personal market view.

### "Update my radar to focus on X"

1. `midaz workspace radar get` — confirm current state.
2. Compose the new item list with the user.
3. `midaz workspace radar set --items "item1, item2, …" --yes`
4. Report back the `l4_enqueued` status and mention that the personal view will refresh shortly.

### "Hook up Telegram"

1. `midaz workspace telegram status` — skip the rest if already connected.
2. `midaz workspace telegram connect` — surface the deep link as a clickable markdown link (also auto-opens).
3. Ask the user to tap Start in Telegram.
4. `midaz workspace telegram status` once or twice to confirm it flipped to `connected: true`.

### "Share my workspace with a friend"

1. Confirm with the user — sharing exposes their personal read.
2. `midaz workspace share --on --yes`
3. Read `workspace.id` from `midaz workspace get`.
4. Share `https://www.midaz.xyz/w/<id>` as a markdown link.

### "Dump everything I know about NVDA"

1. `midaz assets get NVDA`
2. For each thesis in `links[]`, summarize `stance`, `weight`, `rationale`, top `evidence[]` snippets.
3. Drill into the top-weighted link with `midaz assets thesis NVDA <thesis_id>` for full evidence.
4. Share the asset `view_url` + each linked thesis `view_url`.

### "Push my morning note to intel"

1. Read subscription: `midaz subscription status` (must be allowed).
2. Save the note to a file or pipe into one.
3. `midaz intel push --from-file note.md --title "Morning note 2026-04-15" --yes`.
4. Report the returned `id` so they can delete later if needed.
