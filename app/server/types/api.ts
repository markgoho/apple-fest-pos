export type PrintGroup = "kitchen" | "customer";

export type MenuItem = {
  id: string;
  name: string;
  category: string;
  priceCents: number;
  sortOrder: number;
  printGroup: PrintGroup;
};

export type CartLine = {
  menuItemId: string;
  quantity: number;
  notes?: string;
};

export type PlaceOrderRequest = {
  clientOrderId: string;
  deviceId: string;
  payment: {
    method: "cash";
  };
  notes?: string;
  items: CartLine[];
};

export type PrintStatus = "queued" | "printed" | "failed" | "disabled";

export type OrderStatus = "accepted" | "printed" | "print_failed";

export type ReceiptOrder = {
  orderId: string;
  orderNumber: number;
  createdAt: string;
  subtotalCents: number;
  totalCents: number;
  items: CartLine[];
};

export type PlaceOrderResponse = {
  order: {
    id: string;
    orderNumber: number;
    status: OrderStatus;
    subtotalCents: number;
    taxCents: number;
    totalCents: number;
    createdAt: string;
  };
  print: {
    customer: PrintStatus;
    kitchen: PrintStatus;
  };
};

export type HealthResponse = {
  ok: true;
  serverTime: string;
};
