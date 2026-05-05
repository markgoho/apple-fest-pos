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

export type PlaceOrderResponse = {
  order: {
    id: string;
    orderNumber: number;
    status: "accepted" | "printed" | "print_failed";
    subtotalCents: number;
    taxCents: number;
    totalCents: number;
    createdAt: string;
  };
  print: {
    customer: "queued" | "printed" | "failed" | "disabled";
    kitchen: "queued" | "printed" | "failed" | "disabled";
  };
};

export type HealthResponse = {
  ok: true;
  serverTime: string;
};
