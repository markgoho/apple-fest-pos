import { getAdminSales } from "#lib/server/services/sales";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = ({ url }) => ({ sales: getAdminSales(url.searchParams.get("date")) });
