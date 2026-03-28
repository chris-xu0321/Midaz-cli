import { apiFetch, output } from "../api.js";

export async function market(pretty: boolean): Promise<void> {
  output(await apiFetch("/api/market"), pretty);
}
