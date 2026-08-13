import postgres from "postgres";
import { drizzle } from "drizzle-orm/postgres-js";
import type { ServerEnvironment } from "../utils/env";

export function createDatabase(environment: ServerEnvironment) {
  const client = postgres(environment.databaseUrl, { max: 10 });
  return { client, database: drizzle(client) };
}
