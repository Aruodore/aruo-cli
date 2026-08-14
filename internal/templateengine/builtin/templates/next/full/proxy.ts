import { NextResponse, type NextRequest } from "next/server";
import { allowRequest } from "@/server/rate-limit";
import { log } from "@/server/logger";

export function proxy(request: NextRequest) {
  const client =
    request.headers.get("x-forwarded-for")?.split(",")[0]?.trim() || "unknown";
  if (request.nextUrl.pathname.startsWith("/api/") && !allowRequest(client)) {
    return NextResponse.json(
      {
        error: {
          code: "rate_limited",
          message: "Too many requests.",
          requestId: crypto.randomUUID(),
        },
      },
      { status: 429 },
    );
  }
  const response = NextResponse.next();
  log("info", "http_request", {
    method: request.method,
    path: request.nextUrl.pathname,
  });
  response.headers.set("X-Content-Type-Options", "nosniff");
  response.headers.set("Referrer-Policy", "strict-origin-when-cross-origin");
  response.headers.set(
    "Permissions-Policy",
    "camera=(), microphone=(), geolocation=()",
  );
  response.headers.set("X-Frame-Options", "DENY");
  response.headers.set(
    "x-request-id",
    request.headers.get("x-request-id") || crypto.randomUUID(),
  );
  return response;
}

export const config = {
  matcher: "/((?!_next/static|_next/image|favicon.ico).*)",
};
