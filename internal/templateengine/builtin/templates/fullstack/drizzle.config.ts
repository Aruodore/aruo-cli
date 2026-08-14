import { defineConfig } from "drizzle-kit";
import { environment } from "./server/env.js";

const env = environment();
export default defineConfig({
  dialect: "postgresql",
  schema: "./server/db/schema.ts",
  out: "./server/db/migrations",
  dbCredentials: { url: env.DATABASE_URL },
});
