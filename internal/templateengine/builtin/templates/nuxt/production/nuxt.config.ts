export default defineNuxtConfig({
  compatibilityDate: "2026-08-01",
  devtools: { enabled: false },
  runtimeConfig: {
    databaseUrl: "",
    logLevel: "info",
    public: {},
  },
  typescript: { strict: true, typeCheck: true },
});
