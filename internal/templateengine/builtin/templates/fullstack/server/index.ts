import { serve } from "@hono/node-server";
import { serveStatic } from "@hono/node-server/serve-static";
import { Hono } from "hono";
import { secureHeaders } from "hono/secure-headers";
import { requestId } from "hono/request-id";

import { database } from "./database.js";
import { environment } from "./env.js";
import { publicError } from "./errors.js";
import { log } from "./logger.js";
import { allowRequest } from "./rate-limit.js";

const env = environment();
const app = new Hono();

app.use("*", requestId());
app.use("*", secureHeaders());
app.use("/api/*", async (context, next) => {
  const key =
    context.req.header("x-forwarded-for")?.split(",")[0]?.trim() || "unknown";
  if (!allowRequest(key))
    return context.json(
      {
        error: {
          code: "rate_limited",
          message: "Too many requests.",
          requestId: context.get("requestId"),
        },
      },
      429,
    );
  await next();
});
app.use("*", async (context, next) => {
  const started = performance.now();
  await next();
  log("info", "http_request", {
    method: context.req.method,
    path: context.req.path,
    status: context.res.status,
    durationMs: Math.round(performance.now() - started),
    requestId: context.get("requestId"),
  });
});

app.get("/api/health/live", (context) => context.json({ status: "ok" }));
app.get("/api/health/ready", async (context) => {
  const { client } = database(env);
  try {
    await client`select 1`;
    return context.json({ status: "ready" });
  } finally {
    await client.end();
  }
});

app.use("*", serveStatic({ root: "./dist" }));
app.get("*", serveStatic({ path: "./dist/index.html" }));

app.onError((error, context) => {
  const requestIdValue = context.get("requestId") || crypto.randomUUID();
  log("error", "http_error", {
    requestId: requestIdValue,
    errorName: error.name,
  });
  const response = publicError(error, requestIdValue);
  return context.json(response.body, response.status as 500);
});

const server = serve({ fetch: app.fetch, port: env.PORT });
log("info", "server_started", { port: env.PORT });

function shutdown(signal: string) {
  log("info", "server_stopping", { signal });
  server.close((error) => {
    if (error) {
      log("error", "server_stop_failed", { errorName: error.name });
      process.exitCode = 1;
    }
  });
}
process.once("SIGTERM", () => shutdown("SIGTERM"));
process.once("SIGINT", () => shutdown("SIGINT"));
