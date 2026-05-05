import type { PlaceOrderRequest, PlaceOrderResponse } from "../types/api";
import { menuItems } from "./menu";

const ordersByClientId = new Map<string, PlaceOrderResponse>();
const menuItemsById = new Map(menuItems.map((item) => [item.id, item]));
let nextOrderNumber = 100;

export async function handlePlaceOrder(request: Request): Promise<Response> {
  let body: PlaceOrderRequest;

  try {
    body = await request.json();
  } catch {
    return Response.json({ error: "Invalid JSON" }, { status: 400 });
  }

  const validationError = validateOrder(body);
  if (validationError) {
    return Response.json({ error: validationError }, { status: 400 });
  }

  const existingOrder = ordersByClientId.get(body.clientOrderId);
  if (existingOrder) {
    return Response.json(existingOrder);
  }

  const subtotalCents = body.items.reduce((total, line) => {
    const menuItem = menuItemsById.get(line.menuItemId)!;
    return total + menuItem.priceCents * line.quantity;
  }, 0);
  const taxCents = 0;
  const totalCents = subtotalCents + taxCents;

  const response: PlaceOrderResponse = {
    order: {
      id: crypto.randomUUID(),
      orderNumber: nextOrderNumber++,
      status: "accepted",
      subtotalCents,
      taxCents,
      totalCents,
      createdAt: new Date().toISOString()
    },
    print: {
      customer: "queued",
      kitchen: "queued"
    }
  };

  ordersByClientId.set(body.clientOrderId, response);
  return Response.json(response, { status: 201 });
}

function validateOrder(body: PlaceOrderRequest): string | null {
  if (!body.clientOrderId || !body.deviceId) {
    return "Missing order identifiers";
  }

  if (body.payment?.method !== "cash") {
    return "Only cash payments are supported";
  }

  if (!Array.isArray(body.items) || body.items.length === 0) {
    return "Order must contain at least one item";
  }

  for (const line of body.items) {
    if (!Number.isInteger(line.quantity) || line.quantity < 1) {
      return "Item quantity must be a positive integer";
    }

    if (!menuItemsById.has(line.menuItemId)) {
      return `Unknown menu item: ${line.menuItemId}`;
    }
  }

  return null;
}
