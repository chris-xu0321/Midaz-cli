const API_BASE = process.env.SEER_API_URL ?? "http://localhost:4000";

export async function apiFetch(path: string): Promise<unknown> {
  const url = `${API_BASE}${path}`;
  let res: Response;
  try {
    res = await fetch(url);
  } catch (err: unknown) {
    const msg = err instanceof Error ? err.message : String(err);
    console.error(`Error: Cannot connect to Seer API at ${API_BASE}. ${msg}`);
    process.exit(1);
  }
  if (!res.ok) {
    const body = await res.text().catch(() => "");
    console.error(`Error: ${res.status} ${res.statusText} — ${path}`);
    if (body) console.error(body);
    process.exit(1);
  }
  return res.json();
}

export function buildQs(params: Record<string, string | undefined>): string {
  const entries = Object.entries(params).filter(
    (pair): pair is [string, string] => pair[1] !== undefined
  );
  if (entries.length === 0) return "";
  return "?" + new URLSearchParams(entries).toString();
}

export function output(data: unknown, pretty: boolean): void {
  console.log(JSON.stringify(data, null, pretty ? 2 : undefined));
}
