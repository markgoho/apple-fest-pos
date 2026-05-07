import { describe, expect, test } from "bun:test";
import { buildCustomerReceipt, buildKitchenTicket } from "../../app/server/services/escpos";

const order = {
  orderId: "order-1",
  orderNumber: 101,
  createdAt: "2026-05-07T12:34:00.000Z",
  subtotalCents: 2000,
  totalCents: 2000,
  items: [{ menuItemId: "potato-pancake", quantity: 2 }]
};

describe("escpos builders", () => {
  test("builds a customer receipt with cut command", () => {
    const payload = buildCustomerReceipt(order);

    expect(Array.from(payload.slice(0, 2))).toEqual([0x1b, 0x40]);
    expect(Array.from(payload.slice(-3))).toEqual([0x1d, 0x56, 0x00]);
  });

  test("builds a kitchen ticket with emphasis and cut command", () => {
    const payload = buildKitchenTicket(order);

    expect(Array.from(payload.slice(0, 2))).toEqual([0x1b, 0x40]);
    expect(Array.from(payload.slice(-3))).toEqual([0x1d, 0x56, 0x00]);
    expect(Array.from(payload)).toContain(0x11);
  });
});
