import { json } from "@sveltejs/kit";
import { getAdminSales } from "#lib/server/services/sales";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = ({ url }) => json(getAdminSales(url.searchParams.get("date")));
