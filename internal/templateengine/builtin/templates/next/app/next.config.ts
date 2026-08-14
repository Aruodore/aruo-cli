import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  // The compiler API avoids a Next 16.3 CLI-output parsing failure observed
  // with TypeScript 5.9 while preserving Next's production type check.
  experimental: { useTypeScriptCli: false },
};

export default nextConfig;
