export class ApplicationError extends Error {
  constructor(
    readonly status: number,
    readonly code: string,
    message: string,
  ) {
    super(message);
  }
}

export function publicError(error: unknown, requestId: string) {
  if (error instanceof ApplicationError) {
    return {
      status: error.status,
      body: { error: { code: error.code, message: error.message, requestId } },
    };
  }
  return {
    status: 500,
    body: {
      error: {
        code: "internal_error",
        message: "The request could not be completed.",
        requestId,
      },
    },
  };
}
