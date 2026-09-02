import { handlePlaceOrder } from "#lib/server/routes/orders";
import type { RequestHandler } from "./$types";

export const POST: RequestHandler = ({ request }) => handlePlaceOrder(request);
