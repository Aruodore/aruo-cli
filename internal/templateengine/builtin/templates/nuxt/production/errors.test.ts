import { describe, expect, it } from "vitest";
import { ApplicationError } from "../server/utils/errors";

describe("ApplicationError", () => {
  it("carries a safe public error contract", () => {
    const error = new ApplicationError(
      400,
      "invalid_request",
      "Input is invalid.",
    );
    expect({
      statusCode: error.statusCode,
      code: error.code,
      message: error.message,
    }).toEqual({
      statusCode: 400,
      code: "invalid_request",
      message: "Input is invalid.",
    });
  });
});
