import { getAdminSales } from "#lib/server/services/sales";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = ({ url }) => {
  const pinnedDate = url.searchParams.get("date");
  return { pinnedDate, sales: getAdminSales(pinnedDate) };
};
