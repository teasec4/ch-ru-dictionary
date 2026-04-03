import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ params, url }) => {
  try {
    const word = params.word;

    if (!word || typeof word !== "string") {
      return new Response(JSON.stringify({
        error: 'Invalid word parameter',
        message: 'Word parameter is required',
        data: [],
        count: 0
      }), {
        status: 400,
        headers: { "Content-Type": "application/json" },
      });
    }

    const cleanedWord = word.trim();
    if (cleanedWord.length === 0) {
      return new Response(JSON.stringify({
        error: 'Invalid word parameter',
        message: 'Word parameter is required and must be a non-empty string',
        data: [],
        count: 0
      }), {
        status: 400,
        headers: { "Content-Type": "application/json" },
      });
    }

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 5000);

    let queryParam = 'hanzi';
    if (/[a-zA-Z]/.test(cleanedWord) && !/[\u4e00-\u9fff]/.test(cleanedWord)) {
      queryParam = 'pinyin';
    }

    let res;
    try {
      res = await fetch(`http://localhost:8080/api/entries?${queryParam}=${encodeURIComponent(cleanedWord)}&limit=20`, {
        signal: controller.signal
      });
    } catch (e) {
      throw new Error('Backend request timeout');
    } finally {
      clearTimeout(timeoutId);
    }

    if (!res.ok) {
      const errorText = await res.text();
      console.error(`Backend error (${res.status}):`, errorText);
      
      return new Response(JSON.stringify({ 
        error: 'Backend error',
        message: `Backend responded with status: ${res.status}`,
        data: [],
        count: 0
      }), { 
        status: 502,
        headers: { "Content-Type": "application/json" }
      });
    }

    let data;
    try {
      data = await res.json();
    } catch (parseError) {
      console.error('Failed to parse backend response:', parseError);
      
      return new Response(JSON.stringify({ 
        error: 'Invalid response format',
        message: 'Backend returned invalid JSON',
        data: [],
        count: 0
      }), { 
        status: 502,
        headers: { "Content-Type": "application/json" }
      });
    }

    const response = {
      data: (data.data || []).map((entry: any) => ({
        chinese: entry.hanzi,
        pinyin: entry.pinyin || '',
        meanings: (entry.meanings || []).map((m: any) => `${m.index}) ${m.text}`).join('; ')
      })),
      count: data.total || 0,
      message: data.total === 0 ? 'no results' : undefined
    };

    return new Response(JSON.stringify(response), {
      headers: { "Content-Type": "application/json" },
    });
  } catch (e) {
    console.error('Unexpected error:', e);
    return new Response(JSON.stringify({
      error: 'Internal server error',
      message: 'An unexpected error occurred',
      data: [],
      count: 0
    }), {
      status: 500,
      headers: { "Content-Type": "application/json" }
    });
  }
};