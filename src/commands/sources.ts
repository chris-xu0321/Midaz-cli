import { apiFetch, buildQs, output } from "../api.js";

export async function sources(
  flags: { decision?: string; tier?: string },
  pretty: boolean
): Promise<void> {
  const qs = buildQs({ decision: flags.decision, tier: flags.tier });
  output(await apiFetch(`/api/sources${qs}`), pretty);
}
