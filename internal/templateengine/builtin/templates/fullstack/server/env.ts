import { z } from "zod";

const schema = z.object({
  NODE_ENV: z
    .enum(["development", "test", "production"])
    .default("development"),
  PORT: z.coerce.number().int().min(1).max(65535).default(3001),
  DATABASE_URL: z.string().url().startsWith("postgresql://"),
  LOG_LEVEL: z.enum(["debug", "info", "warn", "error"]).default("info"),
});

export type Environment = z.infer<typeof schema>;

let cached: Environment | undefined;

export function environment(
  input: NodeJS.ProcessEnv = process.env,
): Environment {
  if (input === process.env && cached) return cached;
  const parsed = schema.safeParse(input);
  if (!parsed.success) {
    const fields = parsed.error.issues
      .map((issue) => issue.path.join("."))
      .join(", ");
    throw new Error(`Invalid server environment: ${fields}`);
  }
  if (input === process.env) cached = parsed.data;
  return parsed.data;
}
