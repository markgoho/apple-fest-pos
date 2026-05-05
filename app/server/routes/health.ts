import type { HealthResponse } from "../types/api";

export function handleHealth(): Response {
  return Response.json({ ok: true, serverTime: new Date().toISOString() } satisfies HealthResponse);
}
