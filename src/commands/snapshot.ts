import { apiFetch, buildQs, output } from "../api.js";

export async function snapshot(
  flags: { history?: boolean; limit?: string },
  pretty: boolean
): Promise<void> {
  if (flags.history) {
    const qs = buildQs({ limit: flags.limit });
    output(await apiFetch(`/api/global/snapshots${qs}`), pretty);
  } else {
    output(await apiFetch("/api/global/snapshot"), pretty);
  }
}
