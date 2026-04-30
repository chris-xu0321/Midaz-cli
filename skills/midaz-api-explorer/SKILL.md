---
name: midaz-api-explorer
version: 0.7.2
description: "Discovers and introspects Midaz CLI commands and schemas — lists available commands, dumps JSON schemas for any verb, and surfaces flags/arguments for commands not covered by the dedicated skills. Use this skill whenever the user asks what Midaz can do, requests command-level help, or needs a capability the other midaz-* skills don't cover."
when_to_use: "Trigger on phrases like 'what can midaz do', 'list midaz commands', 'midaz help', 'show me the schema for X', 'what flags does Y take', 'is there a midaz command for Z', 'explore the API', 'midaz --help', or when the user describes a need the other midaz skills clearly don't handle."
metadata: {"requires":{"bins":["midaz"]}}
---

# Midaz API Explorer

> Read [midaz-shared](../midaz-shared/SKILL.md) for response format and common rules.

Fallback skill. Use when the user's need isn't obviously covered by `midaz-market`, `midaz-account`, or `midaz-desk`.

## When to Use

- User asks about a command you don't recognize.
- User needs data that the domain skills don't document.
- You want to discover what commands are available.
- You need the exact input/output contract for a command.

## Discovery Flow

### 1. Check domain skills first

Before running `schema`, ask: does `midaz-market`, `midaz-account`, or `midaz-desk` already cover this? If yes, use that skill.

### 2. List every command

```
midaz schema
```

Returns every registered command with description, positional args, flags, and the endpoints it hits. Use this index to pick the right command.

### 3. Inspect one command

```
midaz schema <command>
```

Returns the full contract for that command: arguments, flags, and endpoint mapping.

### 4. Call with --raw

```
midaz <command> [args] [flags] --raw
```

Bypasses the envelope and prints the raw API response. Handy for:

- Inspecting full response structure when the envelope normalizer hides fields.
- Debugging unexpected output.
- Exploring fields not documented in skills.

### 5. Interpret & synthesize

1. Parse the JSON (`.data` payload when not using `--raw`).
2. Identify the relevant fields for the user's question.
3. Synthesize into natural language.
4. When mentioning a driver or thesis by name, make the name itself an inline markdown link to its `view_url` — per-item `view_url` on list/search results and contributions, `.meta.view_url` on single-entity fetches. Never fabricate a URL. Surface page-level `.meta.view_url` separately as `[View on the map](<url>)`.

## Auth & Subscription Gotchas

- Exit code 6 → the command needs auth. Run `midaz auth status`; if it also fails, **run `midaz auth login` yourself** (don't ask the user to run it manually) — the CLI opens the browser, the user signs in, and the PAT is exchanged via the local loopback server. Retry the original command after it returns 0. Headless / SSH / CI / `MIDAZ_NO_BROWSER=1` is the only case where you should fall back to telling the user to use `midaz auth login --paste`.
- Exit code 7 → the command needs an allowed subscription. Run `midaz subscription status`; surface the trial/active state to the user before running `midaz subscription start --yes`.
- Exit code 2 with `confirmation_required` → the command is a write; add `--yes`.

## Example

User: "what stages does the pipeline use for processing?"

1. Not covered by the market skill → use API explorer.
2. `midaz schema decisions` → see `--stage` flag.
3. `midaz decisions --limit 5` → real stage names in output.
4. Summarize the processing stages found.

User: "is there a way to list all my PATs?"

1. `midaz schema` → find `auth` command.
2. `midaz schema auth` → surface `keys {list,create,revoke}`.
3. `midaz auth keys list` → summarize active keys with labels and last-used dates.
