import type { PlaceOrderRequest } from "$lib/types/api";

const databaseName = "apple-fest-pos";
const databaseVersion = 1;
const storeName = "order-outbox";

export type QueuedOrder = {
  clientOrderId: string;
  createdAt: string;
  order: PlaceOrderRequest;
};

export async function enqueueOrder(order: PlaceOrderRequest): Promise<QueuedOrder> {
  const queuedOrder: QueuedOrder = {
    clientOrderId: order.clientOrderId,
    createdAt: new Date().toISOString(),
    order
  };

  await withOutboxDatabase(async (db) => {
    await requestToPromise(db.transaction(storeName, "readwrite").objectStore(storeName).put(queuedOrder));
  });

  return queuedOrder;
}

export async function listQueuedOrders(): Promise<QueuedOrder[]> {
  const orders = await withOutboxDatabase((db) => {
    return requestToPromise<QueuedOrder[]>(db.transaction(storeName).objectStore(storeName).getAll());
  });

  return orders.toSorted((a, b) => a.createdAt.localeCompare(b.createdAt));
}

export async function removeQueuedOrder(clientOrderId: string): Promise<void> {
  await withOutboxDatabase(async (db) => {
    await requestToPromise(db.transaction(storeName, "readwrite").objectStore(storeName).delete(clientOrderId));
  });
}

async function withOutboxDatabase<T>(callback: (db: IDBDatabase) => Promise<T>): Promise<T> {
  const db = await openOutboxDatabase();

  try {
    return await callback(db);
  } finally {
    db.close();
  }
}

function openOutboxDatabase(): Promise<IDBDatabase> {
  return new Promise<IDBDatabase>((resolve, reject) => {
    const request = indexedDB.open(databaseName, databaseVersion);

    request.onupgradeneeded = () => {
      const db = request.result;
      if (!db.objectStoreNames.contains(storeName)) {
        db.createObjectStore(storeName, { keyPath: "clientOrderId" });
      }
    };
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}

function requestToPromise<T = unknown>(request: IDBRequest<T>): Promise<T> {
  return new Promise<T>((resolve, reject) => {
    request.onsuccess = () => resolve(request.result);
    request.onerror = () => reject(request.error);
  });
}
