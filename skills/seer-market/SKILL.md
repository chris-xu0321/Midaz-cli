---
name: seer-market
version: 1.1.0
description: Help traders read markets, manage their workspace, and push intel — via seer-q CLI
metadata: {"requires":{"bins":["seer-q"]}}
---

# Seer Market Intelligence

> Read [seer-shared](../seer-shared/SKILL.md) first for auth and behavior rules.

## Commands

### Read the market (public, no auth)

| What | Command |
|------|---------|
| Full market overview | `seer-q market` |
| Search anything | `seer-q search "QUERY"` |
| All topics | `seer-q topics` |
| One topic + its threads | `seer-q topic ID` |
| List threads | `seer-q threads [--topic ID] [--status active]` |
| One thread + evidence | `seer-q thread ID` |
| Global regime snapshot | `seer-q snapshot [--history] [--limit N]` |

### My workspace (auth required)

| What | Command |
|------|---------|
| My desk | `seer-q ws` |
| First-time onboarding | `seer-q ws onboard --radar @file.md --playbook @file.md` |
| Update my watchlist | `seer-q ws radar "TEXT"` or `seer-q ws radar @file.md` |
| Update my trading rules | `seer-q ws playbook "TEXT"` or `seer-q ws playbook @file.md` |
| My personal market view | `seer-q ws view` |
| See someone's shared view | `seer-q ws view <workspace_id>` |
| My alerts (unread) | `seer-q ws alerts` |
| My alerts (incl. read) | `seer-q ws alerts --all` |
| Mark alert read | `seer-q ws alerts read <alert_id>` |
| Share my view publicly | `seer-q ws share` |
| Revoke public sharing | `seer-q ws unshare` |

> **Share model:** `ws share` flips `workspaces.shared = true` via
> `PATCH /api/ws`. There are no share tokens — your `workspace_id` IS
> the share handle. Give it to another trader; they must be logged in
> themselves and run `seer-q ws view <workspace_id>`, which hits the
> auth-required `/api/workspaces/:workspace_id/view` endpoint. Members
> can read their own workspace regardless of the `shared` flag.

### My intel (auth required)

| What | Command |
|------|---------|
| Push intel | `seer-q intel "CONTENT"` or `seer-q intel "CONTENT" -t "Title"` |
| Push from file | `seer-q intel @article.md` |
| List my intel | `seer-q intel list` |
| Delete intel | `seer-q intel rm ID` |

## What the trader says → what you do

| Trader says | You do |
|---|---|
| "how's the market" / "what's going on" | `seer-q market` → briefing with regime, top movers, key risks |
| "tell me about oil" / "what's happening with X" | `seer-q search "X"` → `seer-q topic ID` or `seer-q thread ID` |
| "what's the bear case for AI" | `seer-q search "AI"` → find bearish threads → focus on risk_case, contradictions |
| "analyze NVDA" | `seer-q search "NVDA"` → fetch relevant threads → synthesize bull/bear/catalysts |
| "what changed recently" | `seer-q snapshot --history` → compare regime shifts |
| "I heard OPEC might cut" | `seer-q intel "Hearing OPEC+ may cut 500k bpd"` → confirm pushed |
| "saw this article about China PMI" | `seer-q intel "China PMI came in at 49.2, below expectations" -t "China PMI Miss"` |
| "I'm watching oil, rates, and BTC" | `seer-q ws radar "Oil supply | Fed rate path | BTC ETF flows"` |
| "I swing trade, max 3% risk" | `seer-q ws playbook "Swing 2-10d. Max 3% position risk. Scale in on breakouts."` |
| "show me my view" / "my market" | `seer-q ws view` |
| "what's my setup" | `seer-q ws` |
| "any new alerts?" | `seer-q ws alerts` |
| "share my view with Alex" | `seer-q ws share` → give them the `workspace_id` returned |
| "let me see Alex's view" | `seer-q ws view <workspace_id>` (requires Alex to have run `ws share`) |
| "what intel have I saved" | `seer-q intel list` |

## Key fields to highlight

**When briefing on market regime:**
- `regime_summary` — the one-liner
- `verdict.stance` + `verdict.confidence` — bullish/bearish + how sure
- `major_drivers` — what's driving the market
- `key_uncertainties` — what could change everything

**When analyzing a thread:**
- `thesis` + `bias` — what the thread says and which way it leans
- `snapshot.assessment` — full analysis
- `snapshot.risk_case` + `snapshot.what_breaks_it` — the counterargument
- `snapshot.assets_exposed` — what trades are affected
- `snapshot.top_contradiction` — the strongest evidence against

**When showing workspace view:**
- `profile.radar` — what they told you they watch
- `global_snapshot` — the current regime
- `topics` — filtered/ranked by relevance to their radar (future)
- `view` — L4-synthesized personal market cognitive view (nullable if not yet computed)
- `has_view` — presence flag; `false` does NOT mean "running", just "no row for current refresh"
- `alerts_unread` — how many unread alerts are queued

**When reading alerts:**
- `level` — urgency tier
- `headline` — short summary
- `detail` — full context
- `asset` — related ticker/topic
- `source_thread_id` — thread the alert was generated from
- `read_at` — null if unread
