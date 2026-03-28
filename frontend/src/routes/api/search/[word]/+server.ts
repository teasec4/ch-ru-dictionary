import type { RequestHandler } from "../$types";

export const GET: RequestHandler = async ({ params }) => {
  try {
    const { word } = params;
  
  // 1. param validation
  if (!word || typeof word !== "string") {
    return new Response(JSON.stringify({
      error: 'Invalid word parameter',
      message: 'Word parameter is required and must be a string',
      data: [],
      count: 0
    }), {
      status: 400,
      headers: { "Content-Type": "application/json" },
    });
  }
  
  // 2. clean the word parameter
  const cleanedWord = word.trim();
  if (cleanedWord.length === 0) {
    return new Response(
      JSON.stringify({
        error: 'Invalid word parameter',
        message: 'Word parameter is required and must be a non-empty string',
        data: [],
        count: 0
      }),
      {
        status: 400,
        headers: { "Content-Type": "application/json" },
      }
    )
  }
  
  // 3. Запрос к бэкенду с таймаутом
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 5000); // 5 секунд таймаут

    let res;
    try {
      res = await fetch("http://localhost:8080/", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({ word: cleanedWord }),
        signal: controller.signal
      });
    } catch (e) {
      if (e === 'AbortError') {
        throw new Error('Backend request timeout');
      }
      throw e;
    } finally {
      clearTimeout(timeoutId);
    }
  
    // 4. Проверка статуса ответа от бэкенда
    if (!res.ok) {
      const errorText = await res.text();
      console.error(`Backend error (${res.status}):`, errorText);
      
      return new Response(
        JSON.stringify({ 
          error: 'Backend error',
          message: `Backend responded with status: ${res.status}`,
          data: [],
          count: 0
        }),
          { 
            status: 502, // Bad Gateway
            headers: { "Content-Type": "application/json" }
          }
        );
      }

    // 5. Парсинг и валидация ответа
        let data;
        try {
          data = await res.json();
        } catch (parseError) {
          console.error('Failed to parse backend response:', parseError);
          
          return new Response(
            JSON.stringify({ 
              error: 'Invalid response format',
              message: 'Backend returned invalid JSON',
              data: [],
              count: 0
            }),
            { 
              status: 502,
              headers: { "Content-Type": "application/json" }
            }
          );
        }

    // 6. Базовая валидация структуры ответа
        if (!data || typeof data !== 'object') {
          return new Response(
            JSON.stringify({ 
              error: 'Invalid response structure',
              message: 'Backend returned invalid data structure',
              data: [],
              count: 0
            }),
            { 
              status: 502,
              headers: { "Content-Type": "application/json" }
            }
          );
        }
  
  return new Response(JSON.stringify(data), {
    headers: { "Content-Type": "application/json" },
  });
  } catch (e) {
    console.error('Unexpected error:', e);
    return new Response(
      JSON.stringify({
        error: 'Internal server error',
        message: 'An unexpected error occurred',
        data: [],
        count: 0
      }),
      {
        status: 500,
        headers: { "Content-Type": "application/json" }
      }
    );
  }
};