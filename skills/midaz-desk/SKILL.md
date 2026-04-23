---
name: midaz-desk
version: 0.7.3
description: "Manages the user's personal Midaz desk — radar (watchlist of assets/theses), playbook (custom trading rules), preferences, sharing, Telegram alerts, private intel notes, and asset tracking — via the midaz CLI. Use this skill whenever the user references their personal desk, radar, watchlist, playbook, or alerts — even if they don't mention Midaz by name."
when_to_use: "Trigger on phrases like 'my desk', 'my radar', 'my watchlist', 'my playbook', 'update my radar', 'add X to radar', 'remove X from my watchlist', 'my preferences', 'share my desk', 'Telegram alerts', 'mute alerts', 'notify me when', 'push a note', 'private intel', 'track [asset]', 'stop tracking X', 'regenerate my radar', 'reonboard', 'what am I watching', 'what's on my radar'. Also trigger on any first-person possessive followed by desk/radar/playbook terminology."
metadata: {"requires":{"bins":["midaz"]}}
---

# Midaz Desk

> Read [midaz-shared](../midaz-shared/SKILL.md) for auth/envelope basics and [midaz-account](../midaz-account/SKILL.md) for signin/onboarding.

Everything a signed-in user can do inside their desk at `/desk/view` and `/desk/settings`, exposed as CLI commands. All write commands require `--yes`.

## Desk Read Commands

```
midaz desk get          # Summary (name, shared flag, subscription, has_invite_access, onboarded)
midaz desk settings     # Owner-only: radar, playbook, telegram status  (GET /api/desk/settings)
midaz desk view         # Personal market read — subscription-gated     (GET /api/desk/view)
```

`desk get` is the cheapest way to inspect state at the start of a session.

## Radar (watchlist)

The radar is a short list of domains / assets / events the user wants Midaz to focus on. Rules:

- ≤12 items
- Each item ≤160 chars

```
midaz desk radar get
midaz desk radar set --items "Fed policy, AI capex, Oil, China CNY" --yes
midaz desk radar set --from-file radar.md --yes
midaz desk radar add --thesis <id> --yes
midaz desk radar add --driver <id> --yes
midaz desk radar add --asset NVDA  --yes
midaz desk radar add --url https://... --title "..." --yes
midaz desk radar add --text "free-form note" --yes
midaz desk radar remove --index 3 --yes
```

Updates enqueue an L4 refresh; `l4_enqueued: true` in the response means the personal market read will recompute soon.

### Entity pins (radar pin / unpin / pins)

Distinct from free-form radar lines: pins attach a specific **entity** (thesis, driver, asset) to the radar with provenance tracking, so the web market view can render a filled pin button and L4 can treat the entity as a first-class watch target.

```
midaz desk radar pin --kind Thesis --source-type thread --source-id <id> --label "Short text"  --yes
midaz desk radar pin --kind Driver --source-type driver --source-id <id> --label "Short text"  --yes
midaz desk radar pin --kind Asset  --source-type asset  --source-id <id> --label "Short text"  --yes

midaz desk radar unpin --source-type thread --source-id <id> --yes

midaz desk radar pins                 # List pins (provenance rows) — includes origin = pin | adopted
```

Rules enforced server-side:
- `kind` is the display label (`Thesis | Driver | Asset`); `source-type` is the DB key (`thread | driver | asset` — the DB schema keeps the legacy `thread` term internally).
- `label` ≤ 160 chars after whitespace collapse.
- 409 `already_pinned` if the (source-type, source-id) is already pinned.
- 409 `radar_full` if the radar already has 12 lines *and* no freeform line can be adopted.
- If a freeform radar line matches the `label`, the pin adopts it (`origin: adopted`); otherwise a new line is appended (`origin: pin`).
- Unpin is origin-aware: `origin=pin` strips the line; `origin=adopted` only removes provenance (the freeform line survives); unknown pairs no-op.

Pin/unpin enqueue an L4 refresh (`l4_enqueued: true` when work was queued).

## Playbook (trading rules)

Markdown, ≤20 000 chars. Describes how the user wants Midaz to interpret the market.

```
midaz desk playbook get
midaz desk playbook set --from-file playbook.md --yes
```

## Preferences

Per-desk preferences stored in `desk_profiles.soul_json.preferences`. Currently one setting: preferred output language for generated narratives.

```
midaz desk preferences get
midaz desk preferences set --language zh-CN --yes
```

Supported languages: `en`, `zh-CN`, `ja`, `ko`, `es`, `fr`. The preference affects LLM-generated narratives (desk view, bias cards); structural fields are unchanged.

## Sharing

Flip the `shared` boolean to expose a read-only desk page at `https://www.midaz.xyz/d/<desk_id>`:

```
midaz desk share --on --yes
midaz desk share --off --yes
```

After enabling, `midaz desk get` returns `desk.shared: true`. The public URL is computed client-side from the desk id.

## Refresh / regenerate (manual rebuilds)

Three owner-triggered verbs enqueue a rebuild. Pick based on scope:

```
midaz desk regenerate --yes     # Rebuild only your personal desk (fast, reuses market state)
midaz desk reonboard  --yes     # Same scope as regenerate but resubmits current radar+playbook
midaz desk refresh    --yes     # Rebuild the whole market pipeline, then your desk (slow, shared)
```

Notes:

- `regenerate` sends no body and is the direct equivalent of the "Regenerate personal desk" button. Subscription-gated; returns `{status:"queued"}` or 409 with `refresh_id` if one is already in flight.
- `reonboard` reads your current `radar_items` + `playbook` via `GET /api/desk/settings`, then POSTs them back to `/api/desk/onboard` unchanged. Because the settings read is subscription-gated, reonboard is too. Response includes `l4_enqueued: true` on success.
- `refresh` triggers a full pipeline refresh (market rebuild + personal desk). Slower and shared across all desks. Returns `{status:"queued"}`.
- All three are fire-and-forget; the rebuild runs asynchronously. Poll `midaz desk view` to see the recomputed personal read once ready.
- For first-time onboarding (not a refresh), use `midaz onboard` instead.

## Telegram Alerts

Connect the desk to the Midaz Telegram bot so alerts are delivered to chat.

```
midaz desk telegram status       # Polls GET /api/desk/settings → telegram.{connected,bot_username}
midaz desk telegram connect      # Prints + opens https://t.me/<bot>?start=<desk_id>
midaz desk telegram disconnect --yes
```

Flow:
1. Run `telegram connect`. The envelope has a `view_url` pointing at the Telegram deep link; the CLI also attempts to auto-open it. Instruct the user to tap "Start" inside Telegram.
2. Poll `midaz desk telegram status`. `telegram.connected` flips to `true` once the bot webhook has stored the chat id.
3. The desk will then receive alerts as they're produced by the L4 pipeline.

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

Read-only; no auth required, but the richest context appears when paired with a desk.

```
midaz assets list [--tier 1|2] [--bias bullish|bearish|neutral|mixed]
midaz assets get <asset_id>
midaz assets timeline <asset_id> [--limit N]
midaz klines <asset_id>
```

Fields of interest:

- `bias.{direction,axis_state,resonance_n}` — aggregate stance
- `driver_contributions[]` — each has `driver_id`, `role`, axis votes (fundamental/macro/flows), `why`
- `view_url` — deep link into the map

## Agent Recipes

### "What's in my desk right now?"

1. `midaz desk get` — prints subscription, onboarded, has_invite_access.
2. If `has_invite_access: false` → point them to `midaz invite redeem`.
3. If `onboarded: false` → point them to `midaz onboard`.
4. If `subscription.allowed: false` → point them to `midaz subscription start`.
5. Otherwise: `midaz desk view` to show the personal market read.

### "Update my radar to focus on X"

1. `midaz desk radar get` — confirm current state.
2. Compose the new item list with the user.
3. `midaz desk radar set --items "item1, item2, …" --yes`
4. Report back the `l4_enqueued` status and mention that the personal market read will refresh shortly. When you name a thesis or driver that the new radar will match, format it as an inline markdown link to its `view_url` — e.g. `will match the **[<driver_name>](<view_url>)** driver`. Look up the URL via `midaz search` if you don't already have it.

### "Change my desk output language"

1. `midaz desk preferences get` — confirm current language.
2. `midaz desk preferences set --language <code> --yes`
3. `midaz desk regenerate --yes` to re-render narratives in the new language.

### "Hook up Telegram"

1. `midaz desk telegram status` — skip the rest if already connected.
2. `midaz desk telegram connect` — surface the deep link as a clickable markdown link (also auto-opens).
3. Ask the user to tap Start in Telegram.
4. `midaz desk telegram status` once or twice to confirm it flipped to `connected: true`.

### "Share my desk with a friend"

1. Confirm with the user — sharing exposes their personal read.
2. `midaz desk share --on --yes`
3. Read `desk.id` from `midaz desk get`.
4. Share `https://www.midaz.xyz/d/<id>` as a markdown link.

### "Dump everything I know about NVDA"

1. `midaz assets get NVDA` — bias direction, driver contributions (each has `view_url` that opens NVDA's own page with that contribution panel highlighted — not a separate driver page).
2. `midaz assets timeline NVDA --limit 20` — recent events.
3. `midaz klines NVDA` — price history.
4. For each contributing driver, `midaz driver <id>` to surface the thesis members (thesis list items also each carry a `view_url`).
5. When you name any contributing driver or thesis in the reply, make the name an inline markdown link to its `view_url` — e.g. `the **[<driver_name>](<view_url>)** driver contributes +0.5 on fundamentals`. The asset's own `view_url` goes at the end as `[View NVDA on the map](<url>)`.

### "Push my morning note to intel"

1. Read subscription: `midaz subscription status` (must be allowed).
2. Save the note to a file or pipe into one.
3. `midaz intel push --from-file note.md --title "Morning note 2026-04-15" --yes`.
4. Report the returned `id` so they can delete later if needed.
