import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ params }) => {
  let res = await fetch(`http://localhost:8080/search/${params.word}`);
  if (!res.ok) {
    return new Response(JSON.stringify({ error: "Failed to fetch" }), {
      status: res.status,
      headers: {
        "Content-Type": "application/json",
      },
    });
  }
  
  let data = await res.json();
  
  return new Response(JSON.stringify({ data: data }),
    {
      headers: {
        "Content-Type": "application/json",
      },
    }
  );
};