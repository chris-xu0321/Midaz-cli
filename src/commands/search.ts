import { apiFetch, output } from "../api.js";

export async function search(query: string, pretty: boolean): Promise<void> {
  output(await apiFetch(`/api/search?q=${encodeURIComponent(query)}`), pretty);
}
