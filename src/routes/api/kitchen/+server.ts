import { json } from "@sveltejs/kit";
import { getKitchenBoard } from "#lib/server/services/kitchen";
import type { RequestHandler } from "./$types";

export const GET: RequestHandler = () => json(getKitchenBoard());
