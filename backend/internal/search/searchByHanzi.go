package search

func (s *Searcher) searchByHanzi(query string) (Result, error) {
	rows, err := s.db.Queryx(`
		SELECT id, headword, pinyin
		FROM entries
		WHERE headword = ?
		`, query)
	if err != nil {
		return Result{}, err
	}

	defer rows.Close()

	entires := scanEntries(s.db, rows)
	
	// Entry not entry ?!?!
	result := Result{
		Data: entires,
		Total: len(entires),
		Page: 1,
		Limit: 10,
	}
	
	return result, nil
}

