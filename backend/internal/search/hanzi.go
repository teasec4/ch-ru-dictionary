package search

// ByHanzi searches for entries by Chinese characters
// exact: only exact match, extended: includes prefix matches
func (s *Searcher) ByHanzi(hanzi string, limit int, extended bool) []Entry {
	// Always get exact match first
	exact := s.byHanziExact(hanzi, limit)

	// If not extended mode, return only exact matches
	if !extended {
		return exact
	}

	// Extended mode: add prefix matches if there's room
	if len(exact) < limit {
		remaining := limit - len(exact)
		prefix := s.byHanziPrefix(hanzi, remaining)

		// Remove duplicates
		existing := make(map[string]bool)
		for _, e := range exact {
			existing[e.Hanzi] = true
		}

		for _, p := range prefix {
			if !existing[p.Hanzi] && len(exact) < limit {
				exact = append(exact, p)
				existing[p.Hanzi] = true
			}
		}
	}

	return exact
}

// byHanziExact returns entries with exact hanzi match
func (s *Searcher) byHanziExact(hanzi string, limit int) []Entry {
	rows, err := s.db.Queryx(`
		SELECT id, hanzi, pinyin FROM entries 
		WHERE hanzi = ?
		ORDER BY LENGTH(hanzi)
		LIMIT ?`, hanzi, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanEntries(s.db, rows)
}

// byHanziPrefix returns entries where hanzi starts with the search term
func (s *Searcher) byHanziPrefix(hanzi string, limit int) []Entry {
	rows, err := s.db.Queryx(`
		SELECT id, hanzi, pinyin FROM entries 
		WHERE hanzi LIKE ? AND hanzi != ?
		ORDER BY LENGTH(hanzi)
		LIMIT ?`, hanzi+"%", hanzi, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanEntries(s.db, rows)
}

// byHanziContains returns entries where hanzi contains the search term (not used currently)
func (s *Searcher) byHanziContains(hanzi string, limit int) []Entry {
	rows, err := s.db.Queryx(`
		SELECT id, hanzi, pinyin FROM entries 
		WHERE hanzi LIKE ? AND hanzi NOT LIKE ? AND hanzi != ?
		ORDER BY LENGTH(hanzi)
		LIMIT ?`, "%"+hanzi+"%", hanzi+"%", hanzi, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()
	return scanEntries(s.db, rows)
}
