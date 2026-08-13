import type { H3Event } from "h3";
import { getRequestId } from "./request-id";

export class ApplicationError extends Error {
  constructor(
    readonly statusCode: number,
    readonly code: string,
    message: string,
  ) {
    super(message);
  }
}

export function sendApplicationError(event: H3Event, error: unknown) {
  const known = error instanceof ApplicationError;
  const statusCode = known ? error.statusCode : 500;
  setResponseStatus(event, statusCode);
  return {
    error: {
      code: known ? error.code : "internal_error",
      message: known ? error.message : "The request could not be completed.",
      requestId: getRequestId(event),
    },
  };
}
