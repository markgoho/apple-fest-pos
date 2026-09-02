import type { MenuItem } from "#lib/types/api";

export const menuItems: MenuItem[] = [
  { id: "potato-pancake", name: "Potato Pancake", category: "Menu", priceCents: 1000, sortOrder: 10, printGroup: "kitchen" },
  { id: "og-toastie", name: "The OG Toastie", category: "Grilled Cheese", priceCents: 500, sortOrder: 20, printGroup: "kitchen" },
  { id: "pizza-toastie", name: "The Pizza Toastie", category: "Grilled Cheese", priceCents: 600, sortOrder: 30, printGroup: "kitchen" },
  { id: "harvest-toastie", name: "The Harvest Toastie", category: "Grilled Cheese", priceCents: 800, sortOrder: 40, printGroup: "kitchen" }
];

export const sortedMenuItems: MenuItem[] = menuItems.toSorted((a, b) => a.sortOrder - b.sortOrder);

export const menuItemsById = new Map(menuItems.map((item) => [item.id, item]));
