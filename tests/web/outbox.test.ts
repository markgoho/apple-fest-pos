import { describe, expect, test } from "bun:test";
import type { PlaceOrderRequest, PlaceOrderResponse } from "../../app/web/src/lib/types/api";

const order: PlaceOrderRequest = {
  clientOrderId: "client-1",
  deviceId: "device-1",
  payment: {
    method: "cash"
  },
  items: [
    {
      menuItemId: "apple-crisp",
      quantity: 2
    }
  ]
};

type ReplayQueue = {
  queued: PlaceOrderRequest[];
  replay: (submit: (order: PlaceOrderRequest) => Promise<PlaceOrderResponse>) => Promise<number>;
};

function createReplayQueue(orders: PlaceOrderRequest[]): ReplayQueue {
  const queued = [...orders];

  return {
    queued,
    async replay(submit: (order: PlaceOrderRequest) => Promise<PlaceOrderResponse>) {
      let syncedCount = 0;

      while (queued.length > 0) {
        const queuedOrder = queued[0];

        await submit(queuedOrder);
        queued.shift();
        syncedCount += 1;
      }

      return syncedCount;
    }
  };
}

describe("outbox replay contract", () => {
  test("preserves client order id during retry", async () => {
    const replayQueue = createReplayQueue([order]);
    const submittedIds: string[] = [];

    await replayQueue.replay(async (submittedOrder) => {
      submittedIds.push(submittedOrder.clientOrderId);
      return createResponse(submittedOrder.clientOrderId);
    });

    expect(submittedIds).toEqual(["client-1"]);
    expect(replayQueue.queued).toHaveLength(0);
  });

  test("replays orders in fifo order", async () => {
    const replayQueue = createReplayQueue([
      order,
      { ...order, clientOrderId: "client-2" },
      { ...order, clientOrderId: "client-3" }
    ]);
    const submittedIds: string[] = [];

    await replayQueue.replay(async (submittedOrder) => {
      submittedIds.push(submittedOrder.clientOrderId);
      return createResponse(submittedOrder.clientOrderId);
    });

    expect(submittedIds).toEqual(["client-1", "client-2", "client-3"]);
  });
});

function createResponse(id: string): PlaceOrderResponse {
  return {
    order: {
      id,
      orderNumber: 100,
      status: "accepted",
      subtotalCents: 1200,
      taxCents: 0,
      totalCents: 1200,
      createdAt: "2026-05-05T00:00:00.000Z"
    },
    print: {
      customer: "queued",
      kitchen: "queued"
    }
  };
}
