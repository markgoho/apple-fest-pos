import { getDatabase } from "../db/sqlite";
import { menuItemsById } from "./menu-items";
import { printOrder } from "./printer-service";
import type { OrderStatus, PlaceOrderRequest, PlaceOrderResponse } from "#lib/types/api";

type TransactionRow = {
  id: string;
  client_order_id: string;
  device_id: string;
  order_number: number;
  business_date: string;
  status: PlaceOrderResponse["order"]["status"];
  subtotal_cents: number;
  tax_cents: number;
  total_cents: number;
  payment_method: string;
  request_json: string;
  customer_print_status: PlaceOrderResponse["print"]["customer"];
  kitchen_print_status: PlaceOrderResponse["print"]["kitchen"];
  created_at: string;
};

const startingOrderNumber = Number(process.env.ORDER_NUMBER_START ?? 100);

type PlaceOrderResult =
  | { ok: true; response: PlaceOrderResponse }
  | { ok: false; error: string; status: number };

export async function placeOrder(body: PlaceOrderRequest): Promise<PlaceOrderResult> {
  const validationError = validateOrder(body);
  if (validationError) {
    return { ok: false, error: validationError, status: 400 };
  }

  const database = getDatabase();
  const existing = database
    .query<TransactionRow, [string]>("SELECT * FROM transactions WHERE client_order_id = ?1")
    .get(body.clientOrderId);

  if (existing) {
    return { ok: true, response: mapRow(existing) };
  }

  const createdAt = new Date().toISOString();
  const businessDate = createdAt.slice(0, 10);
  const subtotalCents = body.items.reduce((total, line) => {
    const menuItem = menuItemsById.get(line.menuItemId)!;
    return total + menuItem.priceCents * line.quantity;
  }, 0);
  const taxCents = 0;
  const totalCents = subtotalCents + taxCents;

  const insertOrder = database.transaction(() => {
    const storedBusinessDate = getMetadataValue("business_date") ?? businessDate;
    let nextOrderNumber = Number(getMetadataValue("next_order_number") ?? startingOrderNumber);

    if (storedBusinessDate !== businessDate) {
      nextOrderNumber = startingOrderNumber;
    }

    const orderNumber = nextOrderNumber;
    setMetadataValue("business_date", businessDate);
    setMetadataValue("next_order_number", String(orderNumber + 1));

    const row: TransactionRow = {
      id: crypto.randomUUID(),
      client_order_id: body.clientOrderId,
      device_id: body.deviceId,
      order_number: orderNumber,
      business_date: businessDate,
      status: "accepted",
      subtotal_cents: subtotalCents,
      tax_cents: taxCents,
      total_cents: totalCents,
      payment_method: body.payment.method,
      request_json: JSON.stringify(body),
      customer_print_status: "queued",
      kitchen_print_status: "queued",
      created_at: createdAt
    };

    database
      .query(
        `INSERT INTO transactions (
          id, client_order_id, device_id, order_number, business_date, status,
          subtotal_cents, tax_cents, total_cents, payment_method, request_json,
          customer_print_status, kitchen_print_status, created_at
        ) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7, ?8, ?9, ?10, ?11, ?12, ?13, ?14)`
      )
      .run(
        row.id,
        row.client_order_id,
        row.device_id,
        row.order_number,
        row.business_date,
        row.status,
        row.subtotal_cents,
        row.tax_cents,
        row.total_cents,
        row.payment_method,
        row.request_json,
        row.customer_print_status,
        row.kitchen_print_status,
        row.created_at
      );

    return row;
  });

  const row = insertOrder();
  const printStatus = await printOrder({
    orderId: row.id,
    orderNumber: row.order_number,
    createdAt: row.created_at,
    subtotalCents: row.subtotal_cents,
    totalCents: row.total_cents,
    items: body.items
  });

  const status = getOrderStatus(printStatus);

  database
    .query(
      `UPDATE transactions
       SET status = ?2, customer_print_status = ?3, kitchen_print_status = ?4
       WHERE id = ?1`
    )
    .run(row.id, status, printStatus.customer, printStatus.kitchen);

  return {
    ok: true,
    response: {
      order: {
        id: row.id,
        orderNumber: row.order_number,
        status,
        subtotalCents: row.subtotal_cents,
        taxCents: row.tax_cents,
        totalCents: row.total_cents,
        createdAt: row.created_at
      },
      print: printStatus
    }
  };

  function getMetadataValue(key: string): string | null {
    const metadata = database
      .query<{ value: string }, [string]>("SELECT value FROM metadata WHERE key = ?1")
      .get(key);

    return metadata?.value ?? null;
  }

  function setMetadataValue(key: string, value: string): void {
    database
      .query("INSERT INTO metadata (key, value) VALUES (?1, ?2) ON CONFLICT(key) DO UPDATE SET value = excluded.value")
      .run(key, value);
  }
}

export function validateOrder(body: PlaceOrderRequest): string | null {
  if (!body.clientOrderId || !body.deviceId) {
    return "Missing order identifiers";
  }

  if (body.payment?.method !== "cash") {
    return "Only cash payments are supported";
  }

  if (!Array.isArray(body.items) || body.items.length === 0) {
    return "Order must contain at least one item";
  }

  if (body.notes && body.notes.length > 500) {
    return "Order notes are too long";
  }

  for (const line of body.items) {
    if (!Number.isInteger(line.quantity) || line.quantity < 1) {
      return "Item quantity must be a positive integer";
    }

    if (line.notes && line.notes.length > 200) {
      return "Item notes are too long";
    }

    if (!menuItemsById.has(line.menuItemId)) {
      return `Unknown menu item: ${line.menuItemId}`;
    }
  }

  return null;
}

function getOrderStatus(printStatus: PlaceOrderResponse["print"]): OrderStatus {
  if (printStatus.customer === "failed" || printStatus.kitchen === "failed") {
    return "print_failed";
  }

  if (printStatus.customer === "printed" && printStatus.kitchen === "printed") {
    return "printed";
  }

  return "accepted";
}

function mapRow(row: TransactionRow): PlaceOrderResponse {
  return {
    order: {
      id: row.id,
      orderNumber: row.order_number,
      status: row.status,
      subtotalCents: row.subtotal_cents,
      taxCents: row.tax_cents,
      totalCents: row.total_cents,
      createdAt: row.created_at
    },
    print: {
      customer: row.customer_print_status,
      kitchen: row.kitchen_print_status
    }
  };
}
