import { getKitchenBoard } from "#lib/server/services/kitchen";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = () => ({ board: getKitchenBoard() });
