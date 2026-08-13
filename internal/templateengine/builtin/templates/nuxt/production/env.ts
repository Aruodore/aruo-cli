import { z } from "zod";

const serverEnvironment = z.object({
  databaseUrl: z.string().url().startsWith("postgres"),
  logLevel: z.enum(["debug", "info", "warn", "error"]).default("info"),
});

export type ServerEnvironment = z.infer<typeof serverEnvironment>;

export function parseServerEnvironment(input: unknown): ServerEnvironment {
  return serverEnvironment.parse(input);
}

export function useServerEnvironment(
  event: Parameters<typeof useRuntimeConfig>[0],
): ServerEnvironment {
  const config = useRuntimeConfig(event);
  return parseServerEnvironment({
    databaseUrl: config.databaseUrl,
    logLevel: config.logLevel,
  });
}
