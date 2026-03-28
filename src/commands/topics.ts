import { apiFetch, output } from "../api.js";

export async function topics(pretty: boolean): Promise<void> {
  output(await apiFetch("/api/topics"), pretty);
}

export async function topic(id: string, pretty: boolean): Promise<void> {
  output(await apiFetch(`/api/topics/${encodeURIComponent(id)}`), pretty);
}
