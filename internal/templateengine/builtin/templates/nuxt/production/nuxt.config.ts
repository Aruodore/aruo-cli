export default defineNuxtConfig({
  compatibilityDate: "2026-08-01",
  devtools: { enabled: false },
  css: ["~/assets/css/main.css"],
  runtimeConfig: {
    databaseUrl: "",
    logLevel: "info",
    public: {},
  },
  typescript: { strict: true, typeCheck: true },
});
