import { handleHealth } from "./routes/health";
import { handleMenu } from "./routes/menu";
import { handlePlaceOrder } from "./routes/orders";

const port = Number(process.env.PORT ?? 3000);
const webRoot = `${import.meta.dir}/../web/build`;
const indexFile = Bun.file(`${webRoot}/index.html`);

Bun.serve({
  port,
  async fetch(request): Promise<Response> {
    const url = new URL(request.url);

    if (url.pathname === "/api/health" && request.method === "GET") {
      return handleHealth();
    }

    if (url.pathname === "/api/menu" && request.method === "GET") {
      return handleMenu();
    }

    if (url.pathname === "/api/orders" && request.method === "POST") {
      return handlePlaceOrder(request);
    }

    if (url.pathname.startsWith("/api/")) {
      return Response.json({ error: "Not found" }, { status: 404 });
    }

    return serveStaticFile(url.pathname);
  },
  development: {
    hmr: true,
    console: true
  }
});

console.log(`Apple Fest POS server listening on http://localhost:${port}`);

async function serveStaticFile(pathname: string): Promise<Response> {
  if (pathname.includes("..")) {
    return new Response("Bad request", { status: 400 });
  }

  if (pathname !== "/") {
    const staticFile = Bun.file(`${webRoot}${pathname}`);
    if (await staticFile.exists()) {
      return new Response(staticFile, { headers: contentHeaders(pathname) });
    }
  }

  if (await indexFile.exists()) {
    return new Response(indexFile, { headers: { "Content-Type": "text/html; charset=utf-8" } });
  }

  return new Response("Run bun run build:web before starting the server.", { status: 503 });
}

function contentHeaders(pathname: string): HeadersInit {
  if (pathname.endsWith(".js")) {
    return { "Content-Type": "text/javascript; charset=utf-8" };
  }

  if (pathname.endsWith(".css")) {
    return { "Content-Type": "text/css; charset=utf-8" };
  }

  if (pathname.endsWith(".json") || pathname.endsWith(".webmanifest")) {
    return { "Content-Type": "application/json; charset=utf-8" };
  }

  if (pathname.endsWith(".svg")) {
    return { "Content-Type": "image/svg+xml" };
  }

  return {};
}
