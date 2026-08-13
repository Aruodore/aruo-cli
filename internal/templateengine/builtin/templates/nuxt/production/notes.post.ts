import { z } from "zod";
import { createDatabase } from "../db/client";
import { notes } from "../db/schema";

const inputSchema = z
  .object({ title: z.string().trim().min(1).max(200) })
  .strict();

export default defineEventHandler(async (event) => {
  try {
    const input = inputSchema.safeParse(await readBody(event));
    if (!input.success)
      throw new ApplicationError(
        400,
        "invalid_request",
        "A title between 1 and 200 characters is required.",
      );

    const { client, database } = createDatabase(useServerEnvironment(event));
    try {
      const [note] = await database
        .insert(notes)
        .values(input.data)
        .returning();
      setResponseStatus(event, 201);
      return { data: note };
    } finally {
      await client.end();
    }
  } catch (error) {
    return sendApplicationError(event, error);
  }
});
