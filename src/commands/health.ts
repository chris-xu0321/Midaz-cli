import { apiFetch, output } from "../api.js";

export async function health(pretty: boolean): Promise<void> {
  output(await apiFetch("/api/health"), pretty);
}
