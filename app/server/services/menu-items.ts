import { menuItems } from "../routes/menu";

export const menuItemsById = new Map(menuItems.map((item) => [item.id, item]));
