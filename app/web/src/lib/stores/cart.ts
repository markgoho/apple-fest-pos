import type { MenuItem, PlaceOrderRequest } from "$lib/types/api";
import { createId } from "$lib/utils/id";

const potatoPancakeId = "potato-pancake";

export type CartItem = {
  id: string;
  menuItem: MenuItem;
  quantity: number;
  notes: string;
};

export type CartState = {
  items: CartItem[];
};

export function createEmptyCart(): CartState {
  return {
    items: []
  };
}

export function addMenuItem(cart: CartState, menuItem: MenuItem): CartState {
  if (menuItem.id === potatoPancakeId) {
    return addCartLine(cart, menuItem);
  }

  const existingItem = cart.items.find((item) => item.menuItem.id === menuItem.id && item.notes === "");

  if (existingItem) {
    return updateQuantity(cart, existingItem.id, existingItem.quantity + 1);
  }

  return addCartLine(cart, menuItem);
}


export function updateQuantity(cart: CartState, cartItemId: string, quantity: number): CartState {
  const updatedItems = cart.items
    .map((item) => {
      if (item.id !== cartItemId) {
        return item;
      }

      return { ...item, quantity };
    })
    .filter((item) => item.quantity > 0);

  return {
    ...cart,
    items: updatedItems
  };
}

export function updateLineNotes(cart: CartState, cartItemId: string, notes: string): CartState {
  const updatedItems = cart.items.map((item) => {
    if (item.id !== cartItemId) {
      return item;
    }

    return { ...item, notes };
  });

  return {
    ...cart,
    items: updatedItems
  };
}

function addCartLine(cart: CartState, menuItem: MenuItem): CartState {
  return {
    ...cart,
    items: [...cart.items, { id: createId(), menuItem, quantity: 1, notes: "" }]
  };
}

export function getSubtotalCents(cart: CartState): number {
  return cart.items.reduce((total, item) => total + item.menuItem.priceCents * item.quantity, 0);
}

export function getTotalCents(cart: CartState): number {
  return getSubtotalCents(cart);
}

export function buildOrderRequest(cart: CartState, clientOrderId: string, deviceId: string): PlaceOrderRequest {
  return {
    clientOrderId,
    deviceId,
    payment: {
      method: "cash"
    },
    items: cart.items.map((item) => ({
      menuItemId: item.menuItem.id,
      quantity: item.quantity,
      notes: item.notes || undefined
    }))
  };
}
