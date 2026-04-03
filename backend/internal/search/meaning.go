package search

func (s *Searcher) ByMeaning(meaning string, page, limit int) ([]Entry, int) {
	offset := (page - 1) * limit

	var total int
	s.db.Get(&total, `SELECT COUNT(DISTINCT entry_id) FROM meanings WHERE text LIKE ?`, "%"+meaning+"%")

	rows, err := s.db.Queryx(`
		SELECT DISTINCT e.id, e.hanzi, e.pinyin 
		FROM entries e
		JOIN meanings m ON e.id = m.entry_id
		WHERE m.text LIKE ?
		ORDER BY LENGTH(e.hanzi)
		LIMIT ? OFFSET ?`,
		"%"+meaning+"%", limit, offset)
	if err != nil {
		return nil, 0
	}
	defer rows.Close()

	return scanEntries(s.db, rows), total
}

func (s *Searcher) Autocomplete(prefix string, limit int) ([]Entry, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 50 {
		limit = 50
	}

	rows, err := s.db.Queryx(`
		SELECT id, hanzi, pinyin FROM entries 
		WHERE hanzi LIKE ? 
		ORDER BY LENGTH(hanzi)
		LIMIT ?`,
		prefix+"%", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.Hanzi, &e.Pinyin); err != nil {
			return nil, err
		}
		entries = append(entries, e)
	}

	return entries, nil
}
