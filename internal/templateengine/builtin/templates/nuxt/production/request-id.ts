import type { H3Event } from "h3";

const requestIdKey = "requestId";

export function getRequestId(event: H3Event): string {
  const value = event.context[requestIdKey];
  return typeof value === "string" ? value : "unknown";
}

export function setRequestId(event: H3Event, value: string): void {
  event.context[requestIdKey] = value;
}
