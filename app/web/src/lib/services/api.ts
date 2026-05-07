import type {
  AdminSalesResponse,
  HealthResponse,
  MenuItem,
  PlaceOrderRequest,
  PlaceOrderResponse
} from "$lib/types/api";

export class ApiError extends Error {
  constructor(
    message: string,
    readonly status: number
  ) {
    super(message);
  }
}

const apiBaseUrl = (import.meta.env.VITE_API_BASE_URL || "").replace(/\/$/, "");

export async function getHealth(): Promise<HealthResponse> {
  return fetchJson<HealthResponse>(apiUrl("/api/health"));
}

export async function getMenu(): Promise<MenuItem[]> {
  return fetchJson<MenuItem[]>(apiUrl("/api/menu"));
}

export async function placeOrder(order: PlaceOrderRequest): Promise<PlaceOrderResponse> {
  return fetchJson<PlaceOrderResponse>(apiUrl("/api/orders"), {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(order)
  });
}

export async function getAdminSales(date?: string): Promise<AdminSalesResponse> {
  const path = date ? `/api/admin/sales?date=${encodeURIComponent(date)}` : "/api/admin/sales";
  return fetchJson<AdminSalesResponse>(apiUrl(path));
}

function apiUrl(path: string): string {
  return `${apiBaseUrl}${path}`;
}

async function fetchJson<T>(input: RequestInfo | URL, init?: RequestInit): Promise<T> {
  let response: Response;

  try {
    response = await fetch(input, init);
  } catch {
    throw new TypeError("Network request failed");
  }

  if (!response.ok) {
    const body = await response.json().catch(() => null);
    const message = typeof body?.error === "string" ? body.error : `Request failed with ${response.status}`;
    throw new ApiError(message, response.status);
  }

  return response.json() as Promise<T>;
}
