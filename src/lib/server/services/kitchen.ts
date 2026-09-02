import { menuItemsById } from "./menu-items";
import { currentBusinessDate, parseOrderRequest, readSalesRows } from "./sales";
import type { KitchenBoard, KitchenTicket } from "#lib/types/api";

const maxTickets = 20;

export function getKitchenBoard(): KitchenBoard {
  const rows = readSalesRows(currentBusinessDate()).slice(0, maxTickets);

  const tickets: KitchenTicket[] = rows.map((row) => {
    const request = parseOrderRequest(row.request_json);
    return {
      orderNumber: row.order_number,
      createdAt: row.created_at,
      notes: request.notes,
      lines: request.items
        .filter((line) => menuItemsById.get(line.menuItemId)?.printGroup !== "customer")
        .map((line) => ({
          menuItemId: line.menuItemId,
          name: menuItemsById.get(line.menuItemId)?.name ?? line.menuItemId,
          quantity: line.quantity,
          notes: line.notes
        }))
    };
  });

  return {
    serverTime: new Date().toISOString(),
    tickets: tickets.filter((ticket) => ticket.lines.length > 0)
  };
}
