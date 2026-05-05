import { ApiError, getHealth, placeOrder } from "$lib/services/api";
import { enqueueOrder, listQueuedOrders, removeQueuedOrder } from "$lib/services/outbox";
import type { PlaceOrderRequest, PlaceOrderResponse } from "$lib/types/api";

export type SyncState = {
  online: boolean;
  queuedCount: number;
  syncing: boolean;
  message: string;
};

export function createSyncState(): SyncState {
  return {
    online: false,
    queuedCount: 0,
    syncing: false,
    message: "Checking server..."
  };
}

export async function refreshReachability(state: SyncState): Promise<SyncState> {
  try {
    await getHealth();
    return { ...state, online: true, message: state.queuedCount > 0 ? `${state.queuedCount} order(s) queued` : "Server online" };
  } catch {
    return { ...state, online: false, message: "Server offline; orders will queue" };
  }
}

export async function submitOrQueue(order: PlaceOrderRequest): Promise<{ response: PlaceOrderResponse | null; queued: boolean }> {
  try {
    const response = await placeOrder(order);
    return { response, queued: false };
  } catch (error) {
    if (error instanceof ApiError) {
      throw error;
    }

    await enqueueOrder(order);
    return { response: null, queued: true };
  }
}

export async function retryQueuedOrders(
  submit: (order: PlaceOrderRequest) => Promise<PlaceOrderResponse> = placeOrder
): Promise<number> {
  const queuedOrders = await listQueuedOrders();
  let syncedCount = 0;

  for (const queuedOrder of queuedOrders) {
    try {
      await submit(queuedOrder.order);
      await removeQueuedOrder(queuedOrder.clientOrderId);
      syncedCount += 1;
    } catch (error) {
      if (error instanceof ApiError) {
        await removeQueuedOrder(queuedOrder.clientOrderId);
        syncedCount += 1;
        continue;
      }

      break;
    }
  }

  return syncedCount;
}
