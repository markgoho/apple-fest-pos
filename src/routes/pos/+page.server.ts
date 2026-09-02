import { sortedMenuItems } from "#lib/server/services/menu-items";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = () => ({ menuItems: sortedMenuItems });
