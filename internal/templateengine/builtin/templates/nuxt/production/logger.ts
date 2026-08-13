type LogLevel = "debug" | "info" | "warn" | "error";

const weights: Record<LogLevel, number> = {
  debug: 10,
  info: 20,
  warn: 30,
  error: 40,
};

export function writeLog(
  configuredLevel: LogLevel,
  level: LogLevel,
  event: string,
  fields: Record<string, unknown> = {},
): void {
  if (weights[level] < weights[configuredLevel]) return;
  const record = JSON.stringify({
    timestamp: new Date().toISOString(),
    level,
    event,
    ...fields,
  });
  if (level === "error") console.error(record);
  else if (level === "warn") console.warn(record);
  else console.log(record);
}
