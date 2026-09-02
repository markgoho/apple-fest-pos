import type { ReceiptOrder } from "#lib/types/api";
import { menuItemsById } from "./menu-items";

const encoder = new TextEncoder();

export function buildCustomerReceipt(order: ReceiptOrder): Uint8Array {
  const lines = [
    "Apple Fest POS",
    `Order #${order.orderNumber}`,
    formatTimestamp(order.createdAt),
    "",
    ...order.items.flatMap((line) => {
      const menuItem = menuItemsById.get(line.menuItemId);
      const total = menuItem ? formatCurrency(menuItem.priceCents * line.quantity) : "$0.00";
      return [`${line.quantity} x ${menuItem?.name ?? line.menuItemId}`, `  ${total}`];
    }),
    "",
    `Total ${formatCurrency(order.totalCents)}`,
    "Thank you!"
  ];

  return concatEscPos([initializePrinter(), encodeLines(lines), cutPaper()]);
}

export function buildKitchenTicket(order: ReceiptOrder): Uint8Array {
  const lines = [
    `ORDER ${order.orderNumber}`,
    formatTimestamp(order.createdAt),
    "",
    ...order.items.map((line) => {
      const menuItem = menuItemsById.get(line.menuItemId);
      return `${line.quantity}  ${(menuItem?.name ?? line.menuItemId).toUpperCase()}`;
    })
  ];

  return concatEscPos([initializePrinter(), doubleSizeOn(), encodeLines(lines), doubleSizeOff(), cutPaper()]);
}

function initializePrinter(): Uint8Array {
  return Uint8Array.from([0x1b, 0x40]);
}

function doubleSizeOn(): Uint8Array {
  return Uint8Array.from([0x1d, 0x21, 0x11]);
}

function doubleSizeOff(): Uint8Array {
  return Uint8Array.from([0x1d, 0x21, 0x00]);
}

function cutPaper(): Uint8Array {
  return Uint8Array.from([0x1d, 0x56, 0x42, 0x08]);
}

function encodeLines(lines: string[]): Uint8Array {
  return encoder.encode(`${lines.join("\r\n")}\r\n\r\n\r\n`);
}

function concatEscPos(chunks: Uint8Array[]): Uint8Array {
  const length = chunks.reduce((total, chunk) => total + chunk.length, 0);
  const output = new Uint8Array(length);
  let offset = 0;

  for (const chunk of chunks) {
    output.set(chunk, offset);
    offset += chunk.length;
  }

  return output;
}

function formatCurrency(cents: number): string {
  return `$${(cents / 100).toFixed(2)}`;
}

function formatTimestamp(value: string): string {
  return new Date(value).toLocaleString("en-US", {
    month: "numeric",
    day: "numeric",
    hour: "numeric",
    minute: "2-digit"
  });
}
