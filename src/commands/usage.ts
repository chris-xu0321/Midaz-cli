import { apiFetch, buildQs, output } from "../api.js";

export async function usage(
  flags: { since?: string },
  pretty: boolean
): Promise<void> {
  const qs = buildQs({ since: flags.since });
  output(await apiFetch(`/api/usage${qs}`), pretty);
}
