const API = process.env.SEER_API_URL || "http://localhost:4000";

async function apiFetch(path: string): Promise<any> {
  const res = await fetch(`${API}${path}`);
  if (!res.ok) {
    const body = await res.text();
    console.error(`API error ${res.status}: ${body}`);
    process.exit(1);
  }
  return res.json();
}

function output(data: any, pretty: boolean): void {
  console.log(pretty ? JSON.stringify(data, null, 2) : JSON.stringify(data));
}

function parseFlags(args: string[]): Record<string, string> {
  const flags: Record<string, string> = {};
  for (let i = 0; i < args.length; i++) {
    if (args[i].startsWith("--")) {
      const key = args[i].slice(2);
      if (i + 1 < args.length && !args[i + 1].startsWith("--")) {
        flags[key] = args[i + 1];
        i++;
      } else {
        flags[key] = "";  // boolean flag
      }
    }
  }
  return flags;
}

function buildQs(params: Record<string, string | undefined>): string {
  const qs = new URLSearchParams();
  for (const [k, v] of Object.entries(params)) {
    if (v !== undefined) qs.set(k, v);
  }
  const s = qs.toString();
  return s ? `?${s}` : "";
}

// --- Commands ---

async function search(query: string, pretty: boolean): Promise<void> {
  const data = await apiFetch(`/api/search?q=${encodeURIComponent(query)}`);
  output(data, pretty);
}

async function market(pretty: boolean): Promise<void> {
  const data = await apiFetch("/api/market");
  output(data, pretty);
}

async function topic(id: string, pretty: boolean): Promise<void> {
  const data = await apiFetch(`/api/topics/${id}`);
  output(data, pretty);
}

async function thread(id: string, pretty: boolean): Promise<void> {
  const data = await apiFetch(`/api/threads/${id}`);
  output(data, pretty);
}

async function topics(pretty: boolean): Promise<void> {
  const data = await apiFetch("/api/topics");
  output(data, pretty);
}

async function threads(flags: Record<string, string>, pretty: boolean): Promise<void> {
  const qs = buildQs({ topic_id: flags.topic, status: flags.status });
  const data = await apiFetch(`/api/threads${qs}`);
  output(data, pretty);
}

async function claims(flags: Record<string, string>, pretty: boolean): Promise<void> {
  const qs = buildQs({
    thread_id: flags.thread,
    source_id: flags.source,
    status: flags.status,
    claim_mode: flags.mode,
  });
  const data = await apiFetch(`/api/claims${qs}`);
  output(data, pretty);
}

async function sources(flags: Record<string, string>, pretty: boolean): Promise<void> {
  const qs = buildQs({ decision: flags.decision, tier: flags.tier });
  const data = await apiFetch(`/api/sources${qs}`);
  output(data, pretty);
}

async function snapshot(flags: Record<string, string>, pretty: boolean): Promise<void> {
  if ("history" in flags || flags.limit) {
    const qs = buildQs({ limit: flags.limit });
    const data = await apiFetch(`/api/global/snapshots${qs}`);
    output(data, pretty);
  } else {
    const data = await apiFetch("/api/global/snapshot");
    output(data, pretty);
  }
}

// --- Arg parsing ---

const args = process.argv.slice(2);
const pretty = args.includes("--pretty");
const filtered = args.filter((a) => a !== "--pretty");
const command = filtered[0];
const rest = filtered.slice(1);
const flags = parseFlags(rest);
const positional = rest.filter((a, i) => !a.startsWith("--") && (i === 0 || !rest[i - 1].startsWith("--")));
const arg = positional.join(" ");

switch (command) {
  case "search":
    if (!arg) { console.error("Usage: seer-q search <query>"); process.exit(1); }
    await search(arg, pretty);
    break;
  case "market":
    await market(pretty);
    break;
  case "topic":
    if (!arg) { console.error("Usage: seer-q topic <id>"); process.exit(1); }
    await topic(arg, pretty);
    break;
  case "thread":
    if (!arg) { console.error("Usage: seer-q thread <id>"); process.exit(1); }
    await thread(arg, pretty);
    break;
  case "topics":
    await topics(pretty);
    break;
  case "threads":
    await threads(flags, pretty);
    break;
  case "claims":
    await claims(flags, pretty);
    break;
  case "sources":
    await sources(flags, pretty);
    break;
  case "snapshot":
    await snapshot(flags, pretty);
    break;
  default:
    console.error(`Seer Query CLI

Commands:
  seer-q search <query>       Search topics, threads, assets
  seer-q market               Global market overview
  seer-q topic <id>           Topic detail with threads
  seer-q thread <id>          Thread detail with claims
  seer-q topics               List all topics with thread counts
  seer-q threads              List all threads
    --topic <id>              Filter by topic
    --status <status>         Filter by status (active/weakening/divided/resolved)
  seer-q claims               List recent claims
    --thread <id>             Filter by thread
    --source <id>             Filter by source
    --status <status>         Filter by status (pending/current/stale/discarded)
    --mode <mode>             Filter by claim mode (observed/interpreted/forecast/attributed)
  seer-q sources              List recent sources
    --decision <d>            Filter by gate decision (process/drop)
    --tier <n>                Filter by source tier (1/2/3)
  seer-q snapshot             Latest global regime snapshot
    --history                 Show snapshot history instead
    --limit <n>               Limit history results (default 10)

Options:
  --pretty                    Pretty-print JSON output

Environment:
  SEER_API_URL                API base URL (default: http://localhost:4000)`);
    process.exit(command ? 1 : 0);
}
