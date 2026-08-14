import { database } from "@/server/database";
import { environment } from "@/server/env";
import { publicError } from "@/server/errors";
import { log } from "@/server/logger";

export async function GET(request: Request) {
  const requestId = request.headers.get("x-request-id") || crypto.randomUUID();
  const { client } = database(environment());
  try {
    await client`select 1`;
    return Response.json({ status: "ready" });
  } catch (error) {
    log("error", "readiness_failed", {
      requestId,
      errorName: error instanceof Error ? error.name : "unknown",
    });
    const response = publicError(error, requestId);
    return Response.json(response.body, { status: 503 });
  } finally {
    await client.end();
  }
}
