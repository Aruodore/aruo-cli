import { defineConfig } from "drizzle-kit";

if (!process.env.NUXT_DATABASE_URL)
  throw new Error("NUXT_DATABASE_URL is required");

export default defineConfig({
  dialect: "postgresql",
  schema: "./server/db/schema.ts",
  out: "./server/db/migrations",
  dbCredentials: { url: process.env.NUXT_DATABASE_URL },
  strict: true,
  verbose: true,
});
