import { apiFetch, buildQs, output } from "../api.js";

export async function decisions(
  flags: {
    stage?: string;
    "entity-type"?: string;
    "entity-id"?: string;
    run?: string;
    limit?: string;
  },
  pretty: boolean
): Promise<void> {
  // --run with no other filters → dedicated sub-route
  if (
    flags.run &&
    !flags.stage &&
    !flags["entity-type"] &&
    !flags["entity-id"]
  ) {
    output(
      await apiFetch(
        `/api/decisions/run/${encodeURIComponent(flags.run)}`
      ),
      pretty
    );
    return;
  }

  const qs = buildQs({
    stage: flags.stage,
    entity_type: flags["entity-type"],
    entity_id: flags["entity-id"],
    pipeline_run_id: flags.run,
    limit: flags.limit,
  });
  output(await apiFetch(`/api/decisions${qs}`), pretty);
}
