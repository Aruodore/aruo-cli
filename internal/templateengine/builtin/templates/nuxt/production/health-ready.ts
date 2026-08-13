import { createDatabase } from "../../db/client";

export default defineEventHandler(async (event) => {
  try {
    const environment = useServerEnvironment(event);
    const { client } = createDatabase(environment);
    try {
      await client`select 1`;
    } finally {
      await client.end();
    }
    return { status: "ready" };
  } catch {
    setResponseStatus(event, 503);
    return { status: "unavailable" };
  }
});
