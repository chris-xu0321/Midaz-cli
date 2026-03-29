Query the Seer market intelligence system to answer: $ARGUMENTS

## Instructions

1. Parse the user's question to determine intent and pick the right command(s):

   | Intent | Commands |
   |---|---|
   | Overall market | `seer-q market` |
   | All topics / hottest topic | `seer-q topics` → analyze by thread_count |
   | Specific sector/asset | `seer-q search "KEYWORDS"` → `seer-q topic ID` |
   | Specific thread/angle | `seer-q search "KEYWORDS"` → `seer-q thread ID` |
   | Analyze an asset | `seer-q search "ASSET"` → multiple topic + thread calls |
   | Latest events/claims | `seer-q claims` |
   | Recent sources | `seer-q sources` |
   | Bear/bull case | `seer-q search "X"` → `seer-q thread ID` → focus on risk_case/contradictions |
   | Market regime history | `seer-q snapshot --history` |
   | Token usage / costs | `seer-q usage` |
   | Decision audit trail | `seer-q decisions` |

2. Synthesize the response in natural language. Focus on:
   - Current bias and thesis
   - Key risks or catalysts
   - Evidence balance (supporting vs contradicting claims)

3. ALWAYS include the `view_url` from `.data` or `.meta` as a clickable link.

4. Output is JSON wrapped in `{ "ok": true, "data": ..., "meta": { ... } }`.
   Access the payload via `.data`. Check `.ok` before processing.
   Use `--raw` for raw API output (no envelope).

## Commands

```bash
# Entity lookup
seer-q search "QUERY"           # Fuzzy search topics, threads, assets
seer-q topic ID                 # Topic detail + threads
seer-q thread ID                # Thread detail + claims + snapshot

# List / browse
seer-q market                   # Global regime + all topics
seer-q topics                   # All topics with thread counts
seer-q threads                  # All threads (--topic ID, --status S)
seer-q claims                   # Latest 100 claims (--thread ID, --source ID, --status S, --mode M)
seer-q sources                  # Latest 100 sources (--decision D, --tier N)

# Snapshots & usage
seer-q snapshot                 # Latest global regime
seer-q snapshot --history       # Regime history (--limit N)
seer-q usage                    # Token usage summary (--since P, default 24h)
seer-q decisions                # Decision log (--stage S, --run ID, --entity-type T, --limit N)
seer-q health                   # API health check

# Configuration & diagnostics
seer-q version                  # CLI version info
seer-q doctor                   # Check API connectivity and config
seer-q config list              # Show active configuration
seer-q schema <command>         # Describe a command's contract
seer-q agent install claude     # Install agent files to workspace
```

All commands output JSON. Use `--format pretty` for indented output.
