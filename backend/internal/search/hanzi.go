package search

import "unicode"

// ByHanzi searches for entries by Chinese characters
// exact: only exact match, extended: includes prefix matches
func (s *Searcher) ByHanzi(hanzi string, limit int, extended bool) []Entry {
	// Always get exact match first
	exact := s.byHanziExact(hanzi, limit)

	// If exact match found, return it
	if len(exact) > 0 {
		if !extended {
			return exact
		}
		// Extended mode: add prefix matches if there's room
		if len(exact) < limit {
			remaining := limit - len(exact)
			prefix := s.byHanziPrefix(hanzi, remaining)
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

	// If no exact match, try to find component characters
	// This helps when the full word is not in dictionary but its parts are
	if len(hanzi) > 1 {
		components := s.byHanziComponents(hanzi, limit)
		return components
	}

	return exact
}

// byHanziComponents searches for individual characters that make up the search term
// This helps when a multi-character word is not in dictionary
func (s *Searcher) byHanziComponents(hanzi string, limit int) []Entry {
	runes := []rune(hanzi)
	var charParts []string
	for _, r := range runes {
		if unicode.Is(unicode.Han, r) {
			charParts = append(charParts, string(r))
		}
	}

	if len(charParts) == 0 {
		return nil
	}

	// Search for each component character
	var allEntries []Entry
	existing := make(map[string]bool)

	for _, char := range charParts {
		if len(allEntries) >= limit {
			break
		}
		entries := s.byHanziExact(char, limit-len(allEntries))
		for _, e := range entries {
			if !existing[e.Hanzi] {
				allEntries = append(allEntries, e)
				existing[e.Hanzi] = true
			}
		}
	}

	return allEntries
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
