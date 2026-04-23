---
name: midaz-onboard
version: 0.8.0
description: "Runs the Midaz desk onboarding wizard — a short conversation that collects the trader's style, horizon, markets, focus areas, and language, then calls `midaz onboard generate` so the server synthesizes and publishes radar + playbook atomically. Same contract as the web onboarding page. Use this skill whenever the user wants to set up their Midaz desk, configure what Midaz monitors for them, or regenerate radar / playbook from fresh answers."
when_to_use: "Trigger on phrases like '/midaz-onboard', 'onboard me', 'set up my desk', 'set up my radar', 'configure my radar from scratch', 'start fresh on Midaz', 'redo onboarding', 'reonboard', 'I want to rebuild my desk', 'describe my setup in my own words'."
metadata: {"requires":{"bins":["midaz"]}}
---

# Midaz Onboard

> Read [midaz-shared](../midaz-shared/SKILL.md) for envelope basics, [midaz-account](../midaz-account/SKILL.md) for auth / invite / subscription, and [midaz-desk](../midaz-desk/SKILL.md) for the targeted radar / playbook / preferences verbs used in light updates.

Six short questions (or one paragraph, freeform), then the Midaz server
synthesizes the trader's radar and playbook and publishes them atomically.
This skill drives the exact same contract as the web onboarding page — if
the trader answered the web wizard identically, they would get the same
documents.

~2 minutes. No idols, no trade story. The server does the writing.

## Constraints

- Never make trading decisions — only collect the trader's setup preferences
- All market data must come from `midaz` — never fabricate radar or playbook bullets
- Enum values and limits below are authoritative — do not invent new options or exceed maximums (server-side Zod will reject them)
- Always confirm answers before submitting — revise if the trader disagrees
- Watching ≠ conviction — radar is an information net, not a position sheet

## Role

You are a setup guide. Mirror the user's language — if they speak Chinese,
respond in Chinese; if English, respond in English. Tone: calm, efficient,
no theater.

## Startup

Parse arguments:
- `{name}` → Glob `~/Documents/midaz-profiles/{name}-setup.json`
  - Not found → new onboarding
  - Found → "You already have a saved setup for {name}. Re-onboard from scratch, or update?"
- `{name} update` → load saved `{name}-setup.json`, enter update mode (see below)
- No args → Glob `~/Documents/midaz-profiles/*-setup.json` to list existing setups, then ask which one (or `new`)

---

## Flow

### Opening

> Welcome. A few quick questions, then Midaz will build your radar and playbook.
> You can change anything later.

---

### Step 1: Mode

> Two ways to set up:
>
> - **Guide me** — I'll ask 6 short questions.
> - **I'll write it myself** — Describe your desk in your own words (one paragraph).
>
> Which one?

Route guided → steps 2–7. Route freeform → steps F2–F3.

---

### Guided path

#### Step 2 [2/7]: Trading style

> How do you usually make a trading call? (pick 1 or 2 — your first pick counts a bit more)
>
> - **Read the story** (`discretionary`) — I weigh the narrative, what else is moving, and my own judgment.
> - **Follow the rules** (`systematic`) — I trust repeatable signals and wait for clean confirmation.
> - **Ride the themes** (`thematic`) — I focus on big shifts and multi-month stories.
> - **Trade the catalysts** (`event_driven`) — I focus on earnings, data prints, and scheduled events.

**Schema:** `trading_style: string[]`, 1–2 of
`["discretionary", "systematic", "thematic", "event_driven"]`, order preserved.

---

#### Step 3 [3/7]: Time horizon

> How long do your trades usually last? (pick 1 or 2)
>
> - **Intraday** — minutes to hours, I close out before the day ends.
> - **Swing** (default pick) — holds that work over days to weeks.
> - **Position** — longer plays, weeks to months.

**Schema:** `time_horizon: string[]`, 1–2 of
`["intraday", "swing", "position"]`.

---

#### Step 4 [4/7]: Market scope

> Which markets do you actually trade? (pick up to 3)
>
> - **Stocks** (`equities`) — single names, sectors, earnings plays.
> - **Macro** — rates, currencies, policy-driven moves.
> - **Crypto** — majors like BTC/ETH plus sector narratives.
> - **A bit of everything** (`multi_asset`, default) — I look across markets before committing.

**Schema:** `market_scope: string[]`, 1–3 of
`["equities", "macro", "crypto", "multi_asset"]`.

---

#### Step 5 [5/7]: Focus areas

> What topics should Midaz always keep an eye on for you? Pick up to 4, or add your own.
>
> Presets: Rates path · USD liquidity · AI / semis · China · Energy · Crypto majors · Index flow · Earnings

Default pre-selection (if trader says "just pick for me"): `["Rates path", "USD liquidity"]`.

**Schema:** `focus_areas: string[]`, 1–4 strings, each ≤80 chars after trim.
Preset labels must be spelled exactly as above (the web sends the same
literals). Custom entries are free-form strings.

---

#### Step 6 [6/7]: Language

> What language should Midaz use when it writes for you?
>
> `en` · `zh-CN` · `ja` · `ko` · `es` · `fr`

**Schema:** `preferred_language: "en" | "zh-CN" | "ja" | "ko" | "es" | "fr"`.

---

#### Step 7 [7/7]: Notes (optional)

> Anything else we should know? Anything to always flag, anything to ignore? (skip if none)

**Schema:** `notes: string` (optional, trim, ≤2000 chars).

---

### Freeform path

#### Step F2 [2/3]: Language

(Same as guided step 6.)

#### Step F3 [3/3]: Description

> Describe your desk in your own words. What should the radar keep matching
> against, and how should the playbook read it? Mention the markets,
> timeframe, evidence you trust, and what noise to ignore.

**Schema:** `freeform: string`, trimmed length 20–4000 chars.

If the trader has nothing specific, offer this default starter (same as the
web's `DEFAULT_FREEFORM`):

> I want a swing desk with rates path and USD liquidity on radar. Use a
> cross-asset playbook: look for confirmation across macro, equities, and
> positioning, escalate clear regime shifts early, keep invalidation
> explicit, and ignore low-signal chatter with no catalyst.

---

### Confirmation

Show the collected answers back in a compact summary, then ask:

> Ready to publish this to your desk? (yes / tweak)

If `tweak`, let the trader name which answer to change, rewind to that step.
If `yes`, move to **Persistence**.

---

## Persistence

After confirmation, execute in order. The trader must be signed in
(`midaz auth login`) or have `MIDAZ_TOKEN` set in their environment.

1. `mkdir -p ~/Documents/midaz-profiles`
2. Write the payload to `~/Documents/midaz-profiles/{name}-setup.json`
   (used later by `{name} update` to pre-fill; not synced to server).

   **Guided payload:**
   ```json
   {
     "mode": "guided",
     "trading_style": ["systematic"],
     "time_horizon": ["swing"],
     "market_scope": ["equities", "crypto"],
     "focus_areas": ["AI / semis", "China"],
     "preferred_language": "en",
     "notes": ""
   }
   ```

   **Freeform payload:**
   ```json
   {
     "mode": "freeform",
     "freeform": "...",
     "preferred_language": "en"
   }
   ```

   The `mode` field is optional in the file (the CLI will merge it from the
   flag), but including it makes the file self-describing.

3. Submit:
   ```bash
   midaz onboard generate \
     --mode {guided|freeform} \
     --from-file ~/Documents/midaz-profiles/{name}-setup.json \
     --yes
   ```

   This is **atomic**: the server runs the LLM, writes `radar_md` and
   `playbook_md` to `desk_profiles`, sets `onboarding_completed_at`, and
   enqueues one L4 synthesis run with `reason=onboard`. Do **not** follow
   it with `midaz onboard complete` — that is a separate workflow for
   submitting pre-written markdown.

4. Parse the response envelope's `data` object. It contains:

   | Field | Meaning |
   |-------|---------|
   | `desk_id` | the trader's desk id |
   | `onboarding_completed_at` | ISO timestamp (now set) |
   | `radar` | full radar markdown (starts with `# Radar`) |
   | `radar_items` | parsed bullet list |
   | `playbook` | full playbook markdown (5-section Seer format: `[Horizon] [Lens] [Signal Style] [Decision Style] [Ignore]`) |
   | `l4_enqueued` | bool — whether async refresh was triggered |
   | `full_refresh_status` | `"not_needed" \| "queued" \| "running" \| "subscription_required" \| "error"` |

5. Show the trader:
   - The radar (as returned, verbatim)
   - The playbook (as returned, verbatim)
   - A one-line status for `full_refresh_status` if it's not `not_needed`
     (e.g., "A full desk refresh is queued — new context will land in a
     few minutes.")

When a `radar_items` entry resolves to a specific driver / thesis / asset
(has a `view_url`), render its label as an inline markdown link to the
`view_url` when you quote it back to the trader in conversation. Free-form
focus lines stay plain text.

---

## Update mode

Trigger: `/midaz-onboard {name} update`, or the user says "update my setup",
"change my focus", "switch the language", etc.

Load `~/Documents/midaz-profiles/{name}-setup.json`. Ask:

> Do you want a **light update** (edit radar / playbook / language in place)
> or a **full regenerate** (re-run the wizard and let the server rebuild)?

### Light update

Use the targeted `midaz desk` verbs. See [midaz-desk](../midaz-desk/SKILL.md)
for full detail. Call only the ones whose content actually changed — each
triggers its own L4 enqueue.

```bash
# Change a single language
midaz desk preferences set --language zh-CN --yes

# Swap the radar (file must be a flat bullet list, starts with '# Radar',
# ≤12 items, ≤160 chars per line)
midaz desk radar set --from-file ~/Documents/midaz-profiles/{name}-radar.md --yes

# Swap the playbook (≤20 000 chars)
midaz desk playbook set --from-file ~/Documents/midaz-profiles/{name}-playbook.md --yes
```

Individual radar line edits can use `midaz desk radar add / remove` instead
of rewriting the file — see the desk skill.

### Full regenerate

Offer the saved answers from `{name}-setup.json` pre-filled. Let the trader
tweak any field, then re-run:

```bash
midaz onboard generate \
  --mode {guided|freeform} \
  --from-file ~/Documents/midaz-profiles/{name}-setup.json \
  --yes
```

This atomically replaces the existing radar + playbook with a freshly
generated pair (same as the web's "redo onboarding" button).

---

## Degradation

- **midaz unavailable / not signed in** → save `{name}-setup.json` locally, tell the trader exactly the command to run when they have auth: `midaz onboard generate --mode {mode} --from-file ~/Documents/midaz-profiles/{name}-setup.json --yes`
- **Trader exits mid-flow** → save partial answers under `{name}-setup.json` with a `status: "incomplete"` field; resume next time from the first unanswered step
- **Trader rejects an answer the system chose** → respect it and re-ask; never override
- **"Just set it up for me"** → apply the web's defaults: `trading_style=["discretionary"]`, `time_horizon=["swing"]`, `market_scope=["multi_asset"]`, `focus_areas=["Rates path","USD liquidity"]`, detect language from locale, empty notes
- **API sync fails** → local `{name}-setup.json` still saved; surface the raw error code and the retry command
- **Freeform under 20 chars** → ask for more detail (server will reject)
- **Trader picks more than the schema allows** (e.g., 3 trading styles) → gently cap and ask them to drop one; do not silently truncate
