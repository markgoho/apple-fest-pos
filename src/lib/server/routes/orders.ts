import type { PlaceOrderRequest } from "#lib/types/api";
import { placeOrder } from "../services/order-service";

export async function handlePlaceOrder(request: Request): Promise<Response> {
  let body: PlaceOrderRequest;

  try {
    body = await request.json();
  } catch {
    return Response.json({ error: "Invalid JSON" }, { status: 400 });
  }

  const result = await placeOrder(body);
  if (!result.ok) {
    return Response.json({ error: result.error }, { status: result.status });
  }

  return Response.json(result.response, { status: 201 });
}
