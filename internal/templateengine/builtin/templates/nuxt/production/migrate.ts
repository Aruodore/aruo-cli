import { migrate } from "drizzle-orm/postgres-js/migrator";
import { createDatabase } from "./client";
import { parseServerEnvironment } from "../utils/env";

const environment = parseServerEnvironment({
  databaseUrl: process.env.NUXT_DATABASE_URL,
  logLevel: process.env.NUXT_LOG_LEVEL,
});
const { client, database } = createDatabase(environment);

try {
  await migrate(database, { migrationsFolder: "server/db/migrations" });
} finally {
  await client.end();
}
