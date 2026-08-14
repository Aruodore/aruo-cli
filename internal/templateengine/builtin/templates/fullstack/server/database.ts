import { drizzle } from "drizzle-orm/postgres-js";
import postgres from "postgres";

import type { Environment } from "./env.js";

export function database(env: Environment) {
  const client = postgres(env.DATABASE_URL, {
    max: 10,
    connect_timeout: 10,
    idle_timeout: 20,
  });
  return { client, db: drizzle(client) };
}
