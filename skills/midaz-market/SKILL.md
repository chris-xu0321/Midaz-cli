---
name: midaz-market
version: 0.7.0
description: Search, browse, and analyze topics, theses, claims, assets, deltas, and market regime via the midaz CLI
metadata: {"requires":{"bins":["midaz"]}}
---

# Midaz Market Intelligence

> Read [midaz-shared](../midaz-shared/SKILL.md) for response format, auth model, and common rules.

Read-only commands for exploring the market map. No `--yes` required; most endpoints are public (no auth).

## Command Reference

### Entity lookup

```
midaz search "QUERY"        # Fuzzy search across topics, theses, assets
midaz topic TOPIC_ID        # Topic detail: thesis, bias, theses within it
midaz thesis THESIS_ID      # Thesis detail: snapshot, claims, evidence counts, market links
```

Note: `thread` / `threads` are deprecated aliases of `thesis` / `theses`. Prefer the new names.

### List / browse

```
midaz market                        # Global regime + topics
midaz topics                        # All topics with thesis counts
midaz theses                        # All theses, newest activity first
midaz theses --topic TOPIC_ID       # Theses within a topic
midaz theses --status active        # Filter by status (active/weakening/divided/resolved)
midaz claims                        # Latest claims
midaz claims --status current       # Filter (pending/current/stale/discarded)
midaz claims --mode observed        # Filter (observed/interpreted/forecast/attributed)
midaz sources                       # Recently ingested sources
midaz sources --decision process    # Processed sources only
midaz sources --tier 1              # Tier-1 only
```

### Assets

```
midaz assets list                    # All assets with bias + thesis_count
midaz assets list --tier 1           # Filter by tier
midaz assets list --bias bullish     # Filter by bias
midaz assets get ASSET_ID            # Asset detail + linked theses + evidence
midaz assets thesis ASSET_ID THID    # Drill-down for one (asset, thesis) pair
```

### Delta (what's new)

```
midaz delta                         # New claims + theses + topics, last 12h
midaz delta --hours 24              # Last 24h (1-168 allowed)
```

### Snapshots

```
midaz snapshot                       # Latest global regime snapshot
```

### Usage & audit

```
midaz usage                          # Token usage (--since P, default 24h)
midaz decisions                      # Decision log (--stage S, --run ID, --entity-type T, --entity-id I, --limit N)
midaz health                         # API health
```

## Query Strategy

| User intent | Commands |
|---|---|
| "how's the market" | `midaz market` |
| List all topics | `midaz topics` |
| Hottest topic | `midaz topics` → highest `thesis_count` or most recent update |
| Sector deep-dive | `midaz search "KEYWORDS"` → `midaz topic ID` |
| Specific thesis | `midaz search "KEYWORDS"` → `midaz thesis ID` |
| Analyze an asset | `midaz assets get TICKER` (primary) + `midaz assets thesis TICKER THID` for detail |
| Latest events | `midaz delta --hours 24` (new in Midaz) or `midaz claims` |
| Claims for a thesis | `midaz thesis THID` (returns embedded `claims[]`) |
| Recent sources | `midaz sources` |
| Bull/bear case for X | `midaz search "X"` → `midaz thesis ID` → look at `risk_case`, contradicting claims |
| Global regime detail | `midaz snapshot` |
| Theses in a topic | `midaz theses --topic TOPIC_ID` |

## Key Response Fields

**Market / Topics**
- `regime_summary` — one-line market regime
- `standing_thesis` / `standing_digest` — topic narrative
- `bias` — bullish / bearish / neutral / mixed / unclear
- `thesis_count` — number of theses
- `view_url` — interactive map link (always surface as markdown)

**Theses**
- `title`, `thesis` — the argument
- `bias`, `status` — stance and lifecycle
- `snapshot.{assessment,conviction,catalysts,outcomes,risk_case,what_breaks_it,assets_exposed,top_contradiction}`
- `supporting_count`, `contradicting_count` — evidence balance
- `market_links[]` — Polymarket markets referenced
- `view_url`

**Claims**
- `statement`, `claim_mode` (observed/interpreted/forecast/attributed)
- `thesis_role` (support / contradiction) — server renamed from `thread_role`
- `event_date`, `status`, `asset_mentions`, `source_id`

**Assets**
- `asset_id`, `name`, `aliases[]`, `asset_class`, `tier`
- `bias`, `bias_score` — tanh-scaled aggregate over linked theses
- `thesis_count`, `bull_count`, `bear_count`, `mixed_count`
- `links[]` — linked theses with `stance`, `weight`, `rationale`, and evidence claims
- `view_url`

**Delta**
- Returns canonical JSONB from `get_recent_delta(p_hours)` — new claims grouped by thesis and topic.

**Global snapshot**
- `regime_summary`, `snapshot.{verdict,major_drivers,key_uncertainties}`

## Examples

User: "how's the market"
→ `midaz market`
→ Summarize regime_summary, verdict, top topics by thesis_count
→ Share the market view_url as a markdown link
→ Share each topic's view_url as you mention them

User: "latest events"
→ `midaz delta --hours 24`
→ Summarize grouped by thesis: what changed, which topics were touched, which assets mentioned

User: "analyze NVDA"
→ `midaz assets get NVDA`
→ Summarize bias, linked theses (bull vs bear), top evidence snippets per thesis
→ Drill into `midaz assets thesis NVDA <thesis_id>` for specific arguments
→ Share `view_url` for the asset page

User: "bear case for AI"
→ `midaz search "AI"`
→ Filter theses where `bias` is bearish/weakening
→ `midaz thesis <id>` for each — focus on `risk_case`, contradicting claims
→ Share thesis view_urls

User: "what claims support thesis X"
→ `midaz thesis X` — returns embedded `claims[]` for the thesis
→ Filter claims where `thesis_role` is "support"
→ Summarize key supporting evidence
