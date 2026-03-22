import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ params }) => {
  // let res = await fetch(`http://localhost:8080/search/${params.word}`);
  console.log(params.word)
  let res = await fetch('http://localhost:8080/', {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ word: params.word }),
  })
  
  // console.log(res)
  
  if (!res.ok) {
    return new Response(JSON.stringify({ error: "Failed to fetch" }), {
      status: res.status,
      headers: {
        "Content-Type": "application/json",
      },
      
    });
  }
  
  let data = await res.json();
  console.log(data)
  
  return new Response(JSON.stringify( data ),
    {
      headers: {
        "Content-Type": "application/json",
      },
    }
  );
};