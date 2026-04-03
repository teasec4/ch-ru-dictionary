import type { RequestHandler } from "./$types";

export const GET: RequestHandler = async ({ params, url }) => {
  try {
    const word = params.word;
    const mode = url.searchParams.get("mode") || "exact";

    if (!word || typeof word !== "string") {
      return new Response(JSON.stringify({
        error: 'Invalid word parameter',
        data: [],
        count: 0
      }), { status: 400, headers: { "Content-Type": "application/json" } });
    }

    const cleanedWord = word.trim();
    if (cleanedWord.length === 0) {
      return new Response(JSON.stringify({
        error: 'Invalid word parameter',
        data: [],
        count: 0
      }), { status: 400, headers: { "Content-Type": "application/json" } });
    }

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 5000);

    let res;
    try {
      res = await fetch(`http://localhost:8080/api/entries?word=${encodeURIComponent(cleanedWord)}&limit=20&mode=${mode}`, {
        signal: controller.signal
      });
    } catch {
      throw new Error('Backend request timeout');
    } finally {
      clearTimeout(timeoutId);
    }

    if (!res.ok) {
      return new Response(JSON.stringify({ 
        error: 'Backend error',
        data: [],
        count: 0
      }), { status: 502, headers: { "Content-Type": "application/json" } });
    }

    let data;
    try {
      data = await res.json();
    } catch {
      return new Response(JSON.stringify({ 
        error: 'Invalid response format',
        data: [],
        count: 0
      }), { status: 502, headers: { "Content-Type": "application/json" } });
    }

    const response = {
      data: (data.data || []).map((entry: any) => ({
        hanzi: entry.hanzi,
        pinyin: entry.pinyin || '',
        meanings: (entry.meanings || []).map((m: any) => ({
          index: m.index,
          text: m.text,
          refs: m.refs || []
        }))
      })),
      total: data.total || 0,
      count: data.total || 0,
      message: data.total === 0 ? 'no results' : undefined
    };

    return new Response(JSON.stringify(response), {
      headers: { "Content-Type": "application/json" },
    });
  } catch {
    return new Response(JSON.stringify({
      error: 'Internal server error',
      data: [],
      count: 0
    }), { status: 500, headers: { "Content-Type": "application/json" } });
  }
};