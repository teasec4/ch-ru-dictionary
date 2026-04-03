type Response = {
  chinese: string;
  pinyin: string;
  meanings: string;
}

type SearchResponse = {
  data: Response[];
  count: number;
  message?: string;
}

export const load = async ({ fetch, params }) => {
  const response = await fetch(`/api/search/${params.word}`);
  const data: SearchResponse = await response.json();
  return {
    data
  }
};