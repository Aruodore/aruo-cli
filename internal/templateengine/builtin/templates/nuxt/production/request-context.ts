export default defineEventHandler((event) => {
  const incoming = getHeader(event, "x-request-id");
  const requestId =
    incoming && /^[A-Za-z0-9._-]{1,128}$/.test(incoming)
      ? incoming
      : crypto.randomUUID();
  setRequestId(event, requestId);
  setHeader(event, "x-request-id", requestId);

  const started = performance.now();
  event.node.res.once("finish", () => {
    let logLevel: "debug" | "info" | "warn" | "error" = "info";
    try {
      logLevel = useServerEnvironment(event).logLevel;
    } catch {
      // Environment validation is reported by readiness; request logging must still work.
    }
    writeLog(logLevel, "info", "http_request", {
      requestId,
      method: event.method,
      path: getRequestURL(event).pathname,
      statusCode: event.node.res.statusCode,
      durationMs: Math.round(performance.now() - started),
    });
  });
});
