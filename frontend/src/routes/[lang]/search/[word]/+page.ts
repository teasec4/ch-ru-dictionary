type Meaning = {
  index: number;
  text: string;
  refs?: string[];
}

type Response = {
  hanzi: string;
  pinyin: string;
  meanings: Meaning[];
}

type SearchResponse = {
  data: Response[];
  total: number;
  count?: number;
  message?: string;
  error?: string;
}

export const load = async ({ fetch, params }) => {
  const response = await fetch(`/api/search/${params.word}`);
  const data: SearchResponse = await response.json();
  return { 
    data: {
      data: data.data || [],
      count: data.total || data.count || 0,
      total: data.total || data.count || 0,
      message: data.message
    }
  };
};