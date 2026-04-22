---
name: midaz-market
version: 0.7.2
description: Search, browse, and analyze drivers, theses, claims, assets, klines, deltas, and market regime via the midaz CLI
metadata: {"requires":{"bins":["midaz"]}}
---

# Midaz Market Intelligence

> Read [midaz-shared](../midaz-shared/SKILL.md) for response format, auth model, and common rules.

Read-only commands for exploring the market map. No `--yes` required; most endpoints are public (no auth).

## Command Reference

### Entity lookup

```
midaz search "QUERY"        # Fuzzy search across drivers, theses, assets
midaz driver DRIVER_ID      # Driver detail: thread members + asset contributions
midaz thesis THESIS_ID      # Thesis detail: snapshot, claims, evidence counts, market links
```

Note: `thread` / `threads` are deprecated aliases of `thesis` / `theses`. Prefer the new names. The old `topics` / `topic` commands are gone — use `drivers` / `driver` instead (driver ontology replaced the topic layer).

### List / browse

```
midaz market                        # Global regime + drivers composite view
midaz drivers                       # All active drivers
midaz driver-links                  # Causal edges between drivers (sphere graph)
midaz theses                        # All theses, newest activity first
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
midaz assets list                    # All assets with bias + contribution counts
midaz assets list --tier 1           # Filter by tier
midaz assets list --bias bullish     # Filter by bias
midaz assets get ASSET_ID            # Asset detail + driver contributions
midaz assets timeline ASSET_ID       # Event timeline for one asset
midaz assets timeline ASSET_ID --limit 50
```

### Klines (price history)

```
midaz klines                         # List assets that have kline coverage
midaz klines ASSET_ID                # Candlestick history + latest for one asset
```

### Delta (what's new)

```
midaz delta                         # New claims + theses + drivers, last 12h
midaz delta --hours 24              # Last 24h (1-168 allowed)
```

### Snapshots

```
midaz snapshot                       # Latest global regime snapshot
```

### Usage & audit

```
midaz usage                          # Token usage (--since P, default 24h)
midaz usage by-run RUN_ID            # Per-run usage breakdown
midaz decisions                      # Decision log (--stage S, --run ID, --entity-type T, --entity-id I, --limit N)
midaz health                         # API health
```

## Query Strategy

| User intent | Commands |
|---|---|
| "how's the market" | `midaz market` |
| List all drivers | `midaz drivers` |
| Hottest driver | `midaz drivers` → highest `thread_count` / recent `driver_delta` |
| Sector deep-dive | `midaz search "KEYWORDS"` → `midaz driver ID` |
| Specific thesis | `midaz search "KEYWORDS"` → `midaz thesis ID` |
| Analyze an asset | `midaz assets get TICKER` + `midaz assets timeline TICKER` |
| Price history | `midaz klines TICKER` |
| Driver graph | `midaz driver-links` |
| Latest events | `midaz delta --hours 24` or `midaz claims` |
| Claims for a thesis | `midaz thesis THID` (returns embedded `claims[]`) |
| Recent sources | `midaz sources` |
| Bull/bear case for X | `midaz search "X"` → `midaz thesis ID` → look at `risk_case`, contradicting claims |
| Global regime detail | `midaz snapshot` |
| Per-run cost | `midaz usage by-run RUN_ID` |

## Key Response Fields

**Drivers**
- `id`, `name`, `summary`, `driver_kind`, `lifecycle`
- `driver_delta` — recent activity score
- `candidate_assets[]` — assets the driver projects onto
- `thread_count` — theses that reference this driver
- `view_url`

**Market (composite)**
- `regime_summary` — one-line market regime
- `global_snapshot` — current regime object
- `drivers[]` — active drivers with summary fields
- `driver_thread_members[]` — (driver_id, thread_id, role) edges
- `driver_links[]` — causal edges between drivers

**Theses**
- `title`, `thesis` — the argument
- `bias`, `status` — stance and lifecycle
- `snapshot.{assessment,conviction,catalysts,outcomes,risk_case,what_breaks_it,assets_exposed,top_contradiction}`
- `supporting_count`, `contradicting_count` — evidence balance
- `market_links[]` — Polymarket markets referenced
- `view_url`

**Claims**
- `statement`, `claim_mode` (observed/interpreted/forecast/attributed)
- `thesis_role` (support / contradiction)
- `event_date`, `status`, `asset_mentions`, `source_id`

**Assets**
- `asset_id`, `name`, `aliases[]`, `asset_class`, `tier`
- `bias.{direction,axis_state,resonance_n}`
- `driver_contributions[]` — per-driver role + fund/macro/flows axis votes + why
- `view_url`

**Asset timeline**
- List of events (claims, deltas, signal changes) for an asset, newest first.

**Klines**
- List: `{ assets: [{ asset_id, driver_count, net_score }] }`
- Detail: `{ history: [...], latest: {...} }`

**Delta**
- Returns canonical JSONB — new claims grouped by thesis and driver.

**Global snapshot**
- `regime_summary`, `snapshot.{verdict,major_drivers,causal_links,key_uncertainties,environment_summary}`

## Examples

User: "how's the market"
→ `midaz market`
→ Summarize regime_summary, verdict, top drivers by `driver_delta` / `thread_count`
→ Share the market view_url as a markdown link
→ Share each driver's view_url as you mention them

User: "latest events"
→ `midaz delta --hours 24`
→ Summarize grouped by thesis: what changed, which drivers were touched, which assets mentioned

User: "analyze NVDA"
→ `midaz assets get NVDA`
→ Summarize bias direction, axis_state, top driver contributions
→ `midaz assets timeline NVDA --limit 20` for recent events
→ Optionally `midaz klines NVDA` for price context
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

User: "which drivers cause which"
→ `midaz driver-links`
→ Summarize top causal edges (`from_driver_id` → `to_driver_id`, `strength`, `explanation`)
