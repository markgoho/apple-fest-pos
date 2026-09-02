import { getDatabase } from "../db/sqlite";
import { menuItems, menuItemsById } from "./menu-items";
import type {
  AdminSalesResponse,
  AdminSalesOrder,
  AdminSalesItemLine,
  CartLine,
  OrderStatus,
  PrintStatus
} from "#lib/types/api";

export type SalesRow = {
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

export function currentBusinessDate(): string {
  return new Date().toISOString().slice(0, 10);
}

export function readSalesRows(businessDate: string): SalesRow[] {
  return getDatabase()
    .query<SalesRow, [string]>(
      `SELECT id, order_number, business_date, status, subtotal_cents, tax_cents,
              total_cents, payment_method, request_json,
              customer_print_status, kitchen_print_status, created_at
       FROM transactions
       WHERE business_date = ?1
       ORDER BY order_number DESC`
    )
    .all(businessDate);
}

export function getAdminSales(date?: string | null): AdminSalesResponse {
  const businessDate = date || currentBusinessDate();
  const rows = readSalesRows(businessDate);
  const orders: AdminSalesOrder[] = rows.map(toOrder);
  const totalCents = orders.reduce((sum, order) => sum + order.totalCents, 0);
  const printFailures = orders.filter((order) => order.status === "print_failed").length;

  return {
    businessDate,
    serverTime: new Date().toISOString(),
    summary: {
      orderCount: orders.length,
      totalCents,
      printFailures
    },
    items: aggregateItems(rows),
    orders
  };
}

export function parseOrderRequest(json: string): { items: CartLine[]; notes?: string } {
  try {
    const parsed = JSON.parse(json);
    if (parsed && Array.isArray(parsed.items)) {
      return { items: parsed.items as CartLine[], notes: parsed.notes };
    }
  } catch {}
  return { items: [] };
}

function toOrder(row: SalesRow): AdminSalesOrder {
  const request = parseOrderRequest(row.request_json);
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
    for (const line of parseOrderRequest(row.request_json).items) {
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
