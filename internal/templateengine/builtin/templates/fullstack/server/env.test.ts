import { describe, expect, it } from "vitest";
import { environment } from "./env.js";

describe("environment", () => {
  it("rejects missing database configuration", () => {
    expect(() => environment({ NODE_ENV: "test" })).toThrow("DATABASE_URL");
  });

  it("parses a valid environment", () => {
    expect(
      environment({
        NODE_ENV: "test",
        DATABASE_URL: "postgresql://localhost/app",
      }).PORT,
    ).toBe(3001);
  });
});
