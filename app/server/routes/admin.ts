import { getDatabase } from "../db/sqlite";
import { menuItemsById } from "../services/menu-items";
import { menuItems } from "./menu";
import type {
  AdminSalesResponse,
  AdminSalesOrder,
  AdminSalesItemLine,
  CartLine,
  OrderStatus,
  PrintStatus
} from "../types/api";

type SalesRow = {
  id: string;
  order_number: number;
  business_date: string;
  status: OrderStatus;
  subtotal_cents: number;
  tax_cents: number;
  total_cents: number;
  payment_method: string;
  request_json: string;
  customer_print_status: PrintStatus;
  kitchen_print_status: PrintStatus;
  created_at: string;
};

export function handleAdminSales(request: Request): Response {
  const url = new URL(request.url);
  const businessDate = url.searchParams.get("date") ?? currentBusinessDate();

  const rows = getDatabase()
    .query<SalesRow, [string]>(
      `SELECT id, order_number, business_date, status, subtotal_cents, tax_cents,
              total_cents, payment_method, request_json,
              customer_print_status, kitchen_print_status, created_at
       FROM transactions
       WHERE business_date = ?1
       ORDER BY order_number DESC`
    )
    .all(businessDate);

  const orders: AdminSalesOrder[] = rows.map(toOrder);
  const items = aggregateItems(rows);
  const totalCents = orders.reduce((sum, order) => sum + order.totalCents, 0);
  const printFailures = orders.filter((order) => order.status === "print_failed").length;

  const response: AdminSalesResponse = {
    businessDate,
    serverTime: new Date().toISOString(),
    summary: {
      orderCount: orders.length,
      totalCents,
      printFailures
    },
    items,
    orders
  };

  return Response.json(response);
}

function currentBusinessDate(): string {
  return new Date().toISOString().slice(0, 10);
}

function toOrder(row: SalesRow): AdminSalesOrder {
  const request = safeParseRequest(row.request_json);
  return {
    id: row.id,
    orderNumber: row.order_number,
    status: row.status,
    subtotalCents: row.subtotal_cents,
    taxCents: row.tax_cents,
    totalCents: row.total_cents,
    paymentMethod: row.payment_method,
    customerPrintStatus: row.customer_print_status,
    kitchenPrintStatus: row.kitchen_print_status,
    createdAt: row.created_at,
    items: request.items.map((line) => ({
      menuItemId: line.menuItemId,
      name: menuItemsById.get(line.menuItemId)?.name ?? line.menuItemId,
      quantity: line.quantity
    }))
  };
}

function aggregateItems(rows: SalesRow[]): AdminSalesItemLine[] {
  const totals = new Map<string, number>();

  for (const row of rows) {
    const request = safeParseRequest(row.request_json);
    for (const line of request.items) {
      totals.set(line.menuItemId, (totals.get(line.menuItemId) ?? 0) + line.quantity);
    }
  }

  return menuItems
    .map((item) => ({
      menuItemId: item.id,
      name: item.name,
      quantity: totals.get(item.id) ?? 0,
      revenueCents: (totals.get(item.id) ?? 0) * item.priceCents
    }))
    .filter((line) => line.quantity > 0)
    .sort((a, b) => b.quantity - a.quantity);
}

function safeParseRequest(json: string): { items: CartLine[] } {
  try {
    const parsed = JSON.parse(json);
    if (parsed && Array.isArray(parsed.items)) {
      return { items: parsed.items as CartLine[] };
    }
  } catch {}
  return { items: [] };
}
