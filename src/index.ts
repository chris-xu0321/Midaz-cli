import { health } from "./commands/health.js";
import { search } from "./commands/search.js";
import { market } from "./commands/market.js";
import { topics, topic } from "./commands/topics.js";
import { threads, thread } from "./commands/threads.js";
import { claims } from "./commands/claims.js";
import { sources } from "./commands/sources.js";
import { snapshot } from "./commands/snapshot.js";
import { usage } from "./commands/usage.js";
import { decisions } from "./commands/decisions.js";

const HELP = `seer-q — Seer market intelligence CLI

Usage: seer-q <command> [args] [flags]

Commands:
  search <query>       Fuzzy search topics, threads, assets
  market               Global regime + all topics
  topics               List all topics
  topic <id>           Topic detail + threads
  threads              List threads (--topic ID, --status S)
  thread <id>          Thread detail + claims
  claims               List claims (--thread ID, --source ID, --status S, --mode M)
  sources              List sources (--decision D, --tier N)
  snapshot             Latest global snapshot (--history, --limit N)
  usage                Token usage (--since P)
  decisions            Decision log (--stage S, --run ID, --entity-type T, --entity-id I, --limit N)
  health               Health check

Flags:
  --pretty             Pretty-print JSON output

Environment:
  SEER_API_URL         API base URL (default: http://localhost:4000)
`;

function parseArgs(argv: string[]): {
  command: string | undefined;
  positional: string[];
  flags: Record<string, string | boolean>;
} {
  const args = argv.slice(2);
  const command = args[0] && !args[0].startsWith("--") ? args[0] : undefined;
  const rest = command ? args.slice(1) : args;
  const positional: string[] = [];
  const flags: Record<string, string | boolean> = {};

  for (let i = 0; i < rest.length; i++) {
    const arg = rest[i];
    if (arg.startsWith("--")) {
      const key = arg.slice(2);
      const next = rest[i + 1];
      if (next && !next.startsWith("--")) {
        flags[key] = next;
        i++;
      } else {
        flags[key] = true;
      }
    } else {
      positional.push(arg);
    }
  }
  return { command, positional, flags };
}

function flag(flags: Record<string, string | boolean>, key: string): string | undefined {
  const v = flags[key];
  return typeof v === "string" ? v : undefined;
}

const { command, positional, flags } = parseArgs(process.argv);
const pretty = flags.pretty === true;

switch (command) {
  case "search": {
    const query = positional[0];
    if (!query) {
      console.error("Usage: seer-q search <query>");
      process.exit(1);
    }
    await search(query, pretty);
    break;
  }
  case "market":
    await market(pretty);
    break;
  case "topics":
    await topics(pretty);
    break;
  case "topic": {
    const id = positional[0];
    if (!id) {
      console.error("Usage: seer-q topic <id>");
      process.exit(1);
    }
    await topic(id, pretty);
    break;
  }
  case "threads":
    await threads({ topic: flag(flags, "topic"), status: flag(flags, "status") }, pretty);
    break;
  case "thread": {
    const id = positional[0];
    if (!id) {
      console.error("Usage: seer-q thread <id>");
      process.exit(1);
    }
    await thread(id, pretty);
    break;
  }
  case "claims":
    await claims(
      {
        thread: flag(flags, "thread"),
        source: flag(flags, "source"),
        status: flag(flags, "status"),
        mode: flag(flags, "mode"),
      },
      pretty
    );
    break;
  case "sources":
    await sources(
      { decision: flag(flags, "decision"), tier: flag(flags, "tier") },
      pretty
    );
    break;
  case "snapshot":
    await snapshot(
      { history: flags.history === true, limit: flag(flags, "limit") },
      pretty
    );
    break;
  case "usage":
    await usage({ since: flag(flags, "since") }, pretty);
    break;
  case "decisions":
    await decisions(
      {
        stage: flag(flags, "stage"),
        "entity-type": flag(flags, "entity-type"),
        "entity-id": flag(flags, "entity-id"),
        run: flag(flags, "run"),
        limit: flag(flags, "limit"),
      },
      pretty
    );
    break;
  case "health":
    await health(pretty);
    break;
  default:
    console.error(HELP);
    process.exit(command ? 1 : 0);
}
