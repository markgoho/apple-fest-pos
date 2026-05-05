import { describe, expect, test } from "bun:test";
import {
  addMenuItem,
  buildOrderRequest,
  createEmptyCart,
  getSubtotalCents,
  getTotalCents,
  updateLineNotes,
  updateQuantity
} from "../../app/web/src/lib/stores/cart";
import type { MenuItem } from "../../app/web/src/lib/types/api";

const potatoPancake: MenuItem = {
  id: "potato-pancake",
  name: "Potato Pancake",
  category: "Menu",
  priceCents: 1000,
  sortOrder: 10,
  printGroup: "kitchen"
};

const pizzaToastie: MenuItem = {
  id: "pizza-toastie",
  name: "The Pizza Toastie",
  category: "Grilled Cheese",
  priceCents: 600,
  sortOrder: 30,
  printGroup: "kitchen"
};

describe("cart math", () => {
  test("adds items and totals cents", () => {
    let cart = createEmptyCart();
    cart = addMenuItem(cart, potatoPancake);
    cart = addMenuItem(cart, potatoPancake);
    cart = addMenuItem(cart, pizzaToastie);

    expect(getSubtotalCents(cart)).toBe(2600);
    expect(getTotalCents(cart)).toBe(2600);
    expect(cart.items).toHaveLength(3);
  });

  test("removes a line when quantity reaches zero", () => {
    let cart = addMenuItem(createEmptyCart(), pizzaToastie);
    cart = updateQuantity(cart, cart.items[0].id, 0);

    expect(cart.items).toEqual([]);
  });

  test("builds stable order payload with notes", () => {
    let cart = addMenuItem(createEmptyCart(), potatoPancake);
    cart = updateLineNotes(cart, cart.items[0].id, "Applesauce");

    expect(buildOrderRequest(cart, "client-1", "device-1")).toEqual({
      clientOrderId: "client-1",
      deviceId: "device-1",
      payment: {
        method: "cash"
      },
      items: [
        {
          menuItemId: "potato-pancake",
          quantity: 1,
          notes: "Applesauce"
        }
      ]
    });
  });

  test("keeps potato pancakes as separate lines", () => {
    let cart = createEmptyCart();
    cart = addMenuItem(cart, potatoPancake);
    cart = addMenuItem(cart, potatoPancake);
    cart = updateLineNotes(cart, cart.items[0].id, "Sour cream");
    cart = updateLineNotes(cart, cart.items[1].id, "Ketchup");

    expect(buildOrderRequest(cart, "client-1", "device-1").items).toEqual([
      {
        menuItemId: "potato-pancake",
        quantity: 1,
        notes: "Sour cream"
      },
      {
        menuItemId: "potato-pancake",
        quantity: 1,
        notes: "Ketchup"
      }
    ]);
  });
});