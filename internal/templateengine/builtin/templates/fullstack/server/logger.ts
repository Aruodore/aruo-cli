const blocked = /authorization|cookie|password|secret|token|database_url/i;

export function log(
  level: "info" | "warn" | "error",
  event: string,
  fields: Record<string, unknown> = {},
) {
  const safe = Object.fromEntries(
    Object.entries(fields).filter(([key]) => !blocked.test(key)),
  );
  process.stdout.write(
    `${JSON.stringify({ timestamp: new Date().toISOString(), level, event, ...safe })}\n`,
  );
}
