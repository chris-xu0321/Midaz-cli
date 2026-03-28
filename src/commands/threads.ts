import { apiFetch, buildQs, output } from "../api.js";

export async function threads(
  flags: { topic?: string; status?: string },
  pretty: boolean
): Promise<void> {
  const qs = buildQs({ topic_id: flags.topic, status: flags.status });
  output(await apiFetch(`/api/threads${qs}`), pretty);
}

export async function thread(id: string, pretty: boolean): Promise<void> {
  output(await apiFetch(`/api/threads/${encodeURIComponent(id)}`), pretty);
}
