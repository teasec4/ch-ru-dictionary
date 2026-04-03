package search

func (s *Searcher) ByMeaning(meaning string, page, limit int) ([]Entry, int) {
	offset := (page - 1) * limit

	var total int
	s.db.Get(&total, `SELECT COUNT(DISTINCT entry_id) FROM meanings WHERE text LIKE ? COLLATE NOCASE`, "%"+meaning+"%")

	rows, err := s.db.Queryx(`
		SELECT DISTINCT e.id, e.hanzi, e.pinyin 
		FROM entries e
		JOIN meanings m ON e.id = m.entry_id
		WHERE m.text LIKE ? COLLATE NOCASE
		ORDER BY LENGTH(e.hanzi)
		LIMIT ? OFFSET ?`,
		"%"+meaning+"%", limit, offset)
	if err != nil {
		return nil, 0
	}
	defer rows.Close()

	return scanEntries(s.db, rows), total
}

