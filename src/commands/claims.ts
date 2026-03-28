import { apiFetch, buildQs, output } from "../api.js";

export async function claims(
  flags: { thread?: string; source?: string; status?: string; mode?: string },
  pretty: boolean
): Promise<void> {
  const qs = buildQs({
    thread_id: flags.thread,
    source_id: flags.source,
    status: flags.status,
    claim_mode: flags.mode,
  });
  output(await apiFetch(`/api/claims${qs}`), pretty);
}
