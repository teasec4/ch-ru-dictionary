package search

func (s *Searcher) ByPinyin(pinyin string, page, limit int) ([]Entry, int) {
	offset := (page - 1) * limit

	var total int
	s.db.Get(&total, `SELECT COUNT(*) FROM entries WHERE pinyin_normalized LIKE ?`, pinyin+"%")

	rows, err := s.db.Queryx(`
		SELECT id, hanzi, pinyin FROM entries 
		WHERE pinyin_normalized LIKE ? 
		ORDER BY LENGTH(hanzi)
		LIMIT ? OFFSET ?`,
		pinyin+"%", limit, offset)
	if err != nil {
		return nil, 0
	}
	defer rows.Close()

	return scanEntries(s.db, rows), total
}
