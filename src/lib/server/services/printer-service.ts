import type { PlaceOrderResponse, ReceiptOrder } from "#lib/types/api";
import { buildCustomerReceipt, buildKitchenTicket } from "./escpos";

type PrintStatus = PlaceOrderResponse["print"];

export type PrinterConfig = {
  enabled: boolean;
  host?: string;
  port: number;
};

export async function printOrder(order: ReceiptOrder): Promise<PrintStatus> {
  const config = getPrinterConfig();
  if (!config.enabled || !config.host) {
    return {
      customer: "disabled",
      kitchen: "disabled"
    };
  }

  const customer = await sendEscPos(config, buildCustomerReceipt(order));
  const kitchen = await sendEscPos(config, buildKitchenTicket(order));

  return {
    customer: customer ? "printed" : "failed",
    kitchen: kitchen ? "printed" : "failed"
  };
}

function getPrinterConfig(): PrinterConfig {
  return {
    enabled: process.env.PRINTER_ENABLED === "true",
    host: process.env.PRINTER_HOST,
    port: Number(process.env.PRINTER_PORT ?? 9100)
  };
}

async function sendEscPos(config: PrinterConfig, payload: Uint8Array): Promise<boolean> {
  let socket: ReturnType<typeof Bun.connect> | null = null;

  try {
    socket = await Bun.connect({
      hostname: config.host!,
      port: config.port,
      socket: {
        data() {},
        open() {}
      }
    });

    socket.write(payload);
    socket.end();
    return true;
  } catch {
    if (socket) {
      socket.end();
    }
    return false;
  }
}
