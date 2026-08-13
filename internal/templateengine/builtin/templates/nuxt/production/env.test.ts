import { describe, expect, it } from "vitest";
import { parseServerEnvironment } from "../server/utils/env";

describe("server environment", () => {
  it("accepts the documented configuration", () => {
    expect(
      parseServerEnvironment({
        databaseUrl: "postgres://localhost/app",
        logLevel: "info",
      }),
    ).toEqual({
      databaseUrl: "postgres://localhost/app",
      logLevel: "info",
    });
  });

  it("rejects a missing database URL", () => {
    expect(() => parseServerEnvironment({ databaseUrl: "" })).toThrow();
  });
});
