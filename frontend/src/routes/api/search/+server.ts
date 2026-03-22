import type { RequestHandler } from "./$types";

export const POST: RequestHandler = async ({ request }) => {
  const { word } = await request.json();

  const res = await fetch("http://localhost:8080/", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ word }),
  });

  const data = await res.json();

  return new Response(JSON.stringify(data), {
    headers: { "Content-Type": "application/json" },
  });
};