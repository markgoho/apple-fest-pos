import { afterEach, beforeEach, describe, expect, test } from "bun:test";
import { closeDatabase } from "../../app/server/db/sqlite";
import { handlePreflight, withCors } from "../../app/server/http/cors";
import { handlePlaceOrder } from "../../app/server/routes/orders";

const sqlitePath = `${import.meta.dir}/tmp-orders.sqlite`;

const validOrder = {
  clientOrderId: "client-1",
  deviceId: "device-1",
  payment: { method: "cash" as const },
  items: [{ menuItemId: "potato-pancake", quantity: 2 }]
};

beforeEach(async () => {
  process.env.SQLITE_PATH = sqlitePath;
  process.env.PRINTER_ENABLED = "false";
  await deleteIfPresent(sqlitePath);
  await deleteIfPresent(`${sqlitePath}-shm`);
  await deleteIfPresent(`${sqlitePath}-wal`);
  closeDatabase();
});

afterEach(() => {
  closeDatabase();
});

describe("handlePlaceOrder", () => {
  test("persists an order and disables printing when printer is unavailable", async () => {
    const response = await handlePlaceOrder(createRequest(validOrder));
    const body = await response.json();

    expect(response.status).toBe(201);
    expect(body.order.orderNumber).toBe(100);
    expect(body.order.totalCents).toBe(2000);
    expect(body.print).toEqual({
      customer: "disabled",
      kitchen: "disabled"
    });
  });

  test("returns the same order for duplicate client ids", async () => {
    const first = await handlePlaceOrder(createRequest(validOrder));
    const second = await handlePlaceOrder(createRequest(validOrder));
    const firstBody = await first.json();
    const secondBody = await second.json();

    expect(second.status).toBe(201);
    expect(secondBody).toEqual(firstBody);
  });

  test("rejects unknown menu items", async () => {
    const response = await handlePlaceOrder(
      createRequest({
        ...validOrder,
        clientOrderId: "client-2",
        items: [{ menuItemId: "unknown-item", quantity: 1 }]
      })
    );

    expect(response.status).toBe(400);
    expect(await response.json()).toEqual({ error: "Unknown menu item: unknown-item" });
  });
});

describe("cors helpers", () => {
  test("answers api preflight requests", () => {
    process.env.PWA_ORIGIN = "https://pos.example";
    const response = handlePreflight(new Request("http://localhost:3000/api/orders", {
      method: "OPTIONS",
      headers: { origin: "https://pos.example" }
    }));

    expect(response.status).toBe(204);
    expect(response.headers.get("Access-Control-Allow-Methods")).toBe("GET,POST,OPTIONS");
  });

  test("adds cors headers to api responses", () => {
    process.env.PWA_ORIGIN = "https://pos.example";
    const response = withCors(
      new Request("http://localhost:3000/api/orders", {
        method: "POST",
        headers: { origin: "https://pos.example" }
      }),
      Response.json({ ok: true })
    );

    expect(response.headers.get("Access-Control-Allow-Headers")).toBe("Content-Type");
  });
});

function createRequest(body: unknown): Request {
  return new Request("http://localhost:3000/api/orders", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });
}

async function deleteIfPresent(path: string): Promise<void> {
  const file = Bun.file(path);
  if (await file.exists()) {
    await file.delete();
  }
}
