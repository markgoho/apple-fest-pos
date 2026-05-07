const allowedOrigin = process.env.PWA_ORIGIN;

export function createCorsHeaders(request: Request): Headers {
  const headers = new Headers();
  const origin = request.headers.get("origin");

  if (allowedOrigin && origin === allowedOrigin) {
    headers.set("Access-Control-Allow-Origin", allowedOrigin);
    headers.set("Vary", "Origin");
  }

  headers.set("Access-Control-Allow-Methods", "GET,POST,OPTIONS");
  headers.set("Access-Control-Allow-Headers", "Content-Type");

  return headers;
}

export function withCors(request: Request, response: Response): Response {
  const headers = new Headers(response.headers);
  const corsHeaders = createCorsHeaders(request);

  for (const [key, value] of corsHeaders.entries()) {
    headers.set(key, value);
  }

  return new Response(response.body, {
    status: response.status,
    statusText: response.statusText,
    headers
  });
}

export function handlePreflight(request: Request): Response {
  return new Response(null, {
    status: 204,
    headers: createCorsHeaders(request)
  });
}
