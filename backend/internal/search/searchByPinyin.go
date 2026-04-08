package search

func (s *Searcher) searchByPinyin(query string) (Result, error) {
	rows, err := s.db.Queryx(`
		SELECT id, headword, pinyin
		FROM entries
		WHERE pinyin_normalized = ?
			OR pinyin_normalized LIKE ? || '%'
			OR pinyin_normalized LIKE '%' || ? || '%'
		ORDER BY
			CASE
				WHEN pinyin_normalized = ? THEN 1
				WHEN pinyin_normalized LIKE ? || '%' THEN 2
				ELSE 3
			END,
			LENGTH(headword) ASC
		LIMIT 230
		`, query, query, query, query, query)
	if err != nil {
		return Result{}, err
	}
	
	defer rows.Close()
	
	entries := scanEntries(s.db, rows)
	
	result := Result{
		Data:  entries,
        Total: len(entries),
        Page:  1,
        Limit: 50,
	}
	
	return result, nil
}
