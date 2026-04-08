package search

import (
	"sort"
	"strings"
	"unicode"

	"github.com/jmoiron/sqlx"
)

type QueryType int

const (
	QueryTypeUnknown QueryType = iota
	QueryTypeHanzi
	QueryTypePinyin
	QueryTypeMeaning
)

func detectQueryType(query string) QueryType {
	query = strings.TrimSpace(query)
	if query == "" {
		return QueryTypeUnknown
	}

	runes := []rune(query)
	for _, r := range runes {
		if unicode.Is(unicode.Han, r) {
			return QueryTypeHanzi
		}
	}

	if strings.ContainsAny(query, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return QueryTypePinyin
	}

	return QueryTypeMeaning
}

func normalizeQuery(query string) string {
	query = strings.TrimSpace(query)
	query = strings.ToLower(query)
	query = strings.ReplaceAll(query, " ", "")
	query = strings.ReplaceAll(query, "\u00A0", "")
	query = removeDiacritics(query)
	query = strings.ReplaceAll(query, "'", "")
	return query
}

func removeDiacritics(s string) string {
	toneMarks := map[rune]rune{
		'à': 'a', 'á': 'a', 'ǎ': 'a', 'ā': 'a',
		'è': 'e', 'é': 'e', 'ě': 'e', 'ē': 'e',
		'ì': 'i', 'í': 'i', 'ǐ': 'i', 'ī': 'i',
		'ò': 'o', 'ó': 'o', 'ǒ': 'o', 'ō': 'o',
		'ù': 'u', 'ú': 'u', 'ǔ': 'u', 'ū': 'u',
		'ǜ': 'ü', 'ǚ': 'ü', 'ǘ': 'ü', 'ǖ': 'ü',
	}

	runes := []rune(s)
	for i, r := range runes {
		if replacement, ok := toneMarks[r]; ok {
			runes[i] = replacement
		}
	}

	return string(runes)
}

func (s *Searcher) Search(query string) (Result, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return Result{Data: []Entry{}, Total: 0, Page: 1, Limit: 50}, nil
	}

	queryType := detectQueryType(query)

	switch queryType {
	case QueryTypeHanzi:
		return s.searchByHanzi(query)
	case QueryTypePinyin:
		normalized := normalizeQuery(query)
		return s.searchByPinyin(normalized)
	case QueryTypeMeaning:
		return s.searchByMeaning(query)
	default:
		return Result{Data: []Entry{}, Total: 0, Page: 1, Limit: 50}, nil
	}
}

// func (s *Searcher) searchByHanzi(query string) (Result, error) {
// 	rows, err := s.db.Queryx(`
// 		SELECT id, headword, pinyin, frequency, rank FROM (
// 			SELECT id, headword, pinyin, frequency, 1 as rank
// 			FROM entries
// 			WHERE headword = ?
// 			UNION ALL
// 			SELECT id, headword, pinyin, frequency, 2 as rank
// 			FROM entries
// 			WHERE headword LIKE ? || '%' AND headword != ?
// 			UNION ALL
// 			SELECT id, headword, pinyin, frequency, 3 as rank
// 			FROM entries
// 			WHERE headword LIKE '%' || ? || '%'
// 			  AND headword NOT LIKE ? || '%'
// 			  AND headword != ?
// 			  AND LENGTH(headword) <= LENGTH(?) + 3
// 		)
// 		ORDER BY rank ASC, LENGTH(headword) ASC, frequency DESC
// 		LIMIT 50`,
// 		query, query, query, query, query, query, query)
// 	if err != nil {
// 		return Result{}, err
// 	}
// 	defer rows.Close()

// 	entries := scanEntries(s.db, rows)
// 	return Result{Data: entries, Total: len(entries), Page: 1, Limit: 50}, nil
// }

func (s *Searcher) searchByPinyin(query string) (Result, error) {
	rows, err := s.db.Queryx(`
		SELECT id, headword, pinyin, frequency, rank FROM (
			SELECT id, headword, pinyin, frequency, 1 as rank
			FROM entries
			WHERE pinyin_normalized = ?
			UNION ALL
			SELECT id, headword, pinyin, frequency, 2 as rank
			FROM entries
			WHERE pinyin_normalized LIKE ? || '%' AND pinyin_normalized != ?
			UNION ALL
			SELECT id, headword, pinyin, frequency, 3 as rank
			FROM entries
			WHERE pinyin_normalized LIKE '%' || ? || '%'
			  AND pinyin_normalized NOT LIKE ? || '%'
			  AND pinyin_normalized != ?
			  AND LENGTH(headword) <= LENGTH(?) + 3
		)
		ORDER BY rank ASC, LENGTH(headword) ASC, frequency DESC
		LIMIT 50`,
		query, query, query, query, query, query, query)
	if err != nil {
		return Result{}, err
	}
	defer rows.Close()

	entries := scanEntries(s.db, rows)
	return Result{Data: entries, Total: len(entries), Page: 1, Limit: 50}, nil
}

func (s *Searcher) searchByMeaning(query string) (Result, error) {
    // Нормализуем запрос для поиска (нижний регистр + убираем лишние пробелы)
    normalized := strings.ToLower(strings.TrimSpace(query))

    rows, err := s.db.Queryx(`
SELECT id, headword, pinyin, frequency, rank FROM (
    -- Точное совпадение целого слова (самый высокий приоритет)
    SELECT e.id, e.headword, e.pinyin, e.frequency, 1 as rank
    FROM entries e
    JOIN meanings m ON e.id = m.entry_id
    WHERE 
        (',' || LOWER(m.text) || ',') LIKE '%,' || ? || ',%' 
        OR LOWER(m.text) = ?

    UNION ALL

    -- Содержит слово (с приоритетом 2)
    SELECT e.id, e.headword, e.pinyin, e.frequency, 2 as rank
    FROM entries e
    JOIN meanings m ON e.id = m.entry_id
    WHERE LOWER(m.text) LIKE '%' || ? || '%'
      AND NOT (',' || LOWER(m.text) || ',') LIKE '%,' || ? || ',%'
      AND LOWER(m.text) != ?
)
ORDER BY rank ASC, LENGTH(headword) ASC, frequency DESC
LIMIT 50
`, normalized, normalized, normalized, normalized, normalized)

    if err != nil {
        return Result{}, err
    }
    defer rows.Close()

    entries := scanEntries(s.db, rows)
    return Result{Data: entries, Total: len(entries), Page: 1, Limit: 50}, nil
}

func scanEntries(db *sqlx.DB, rows *sqlx.Rows) []Entry {
	type rawEntry struct {
		ID        int
		Headword  string
		Pinyin    string
		Frequency int
		Rank      int
	}

	var rawEntries []rawEntry
	for rows.Next() {
		var e rawEntry
		if err := rows.StructScan(&e); err != nil {
			continue
		}
		rawEntries = append(rawEntries, e)
	}

	if err := rows.Err(); err != nil {
		return []Entry{}
	}

	if len(rawEntries) == 0 {
		return []Entry{}
	}

	sort.Slice(rawEntries, func(i, j int) bool {
		if rawEntries[i].Rank != rawEntries[j].Rank {
			return rawEntries[i].Rank < rawEntries[j].Rank
		}
		if len(rawEntries[i].Headword) != len(rawEntries[j].Headword) {
			return len(rawEntries[i].Headword) < len(rawEntries[j].Headword)
		}
		return rawEntries[i].Frequency > rawEntries[j].Frequency
	})

	entryIDs := make([]int, 0, len(rawEntries))
	entryMap := make(map[int]*Entry)
	for _, e := range rawEntries {
		if _, exists := entryMap[e.ID]; !exists {
			entryIDs = append(entryIDs, e.ID)
		}
		entryMap[e.ID] = &Entry{Hanzi: e.Headword, Pinyin: e.Pinyin, Meanings: []Meaning{}}
	}

	if len(entryIDs) > 0 {
		loadMeanings(db, entryIDs, entryMap)
	}

	result := make([]Entry, 0, len(entryIDs))
	for _, id := range entryIDs {
		result = append(result, *entryMap[id])
	}

	return result
}

func loadMeanings(db *sqlx.DB, entryIDs []int, entryMap map[int]*Entry) {
	if len(entryIDs) == 0 {
		return
	}

	placeholders := strings.Repeat("?,", len(entryIDs))
	placeholders = placeholders[:len(placeholders)-1]
	query := `SELECT id, entry_id, order_num, text FROM meanings WHERE entry_id IN (` + placeholders + `) ORDER BY entry_id, order_num`

	args := make([]interface{}, len(entryIDs))
	for i, id := range entryIDs {
		args[i] = id
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	meaningIDs := make([]int, 0)
	for rows.Next() {
		var m struct {
			ID       int
			EntryID  int
			OrderNum int
			Text     string
		}
		if err := rows.Scan(&m.ID, &m.EntryID, &m.OrderNum, &m.Text); err != nil {
			continue
		}
		if e, ok := entryMap[m.EntryID]; ok {
			e.Meanings = append(e.Meanings, Meaning{ID: m.ID, Index: len(e.Meanings) + 1, Text: m.Text})
			meaningIDs = append(meaningIDs, m.ID)
		}
	}

	if len(meaningIDs) > 0 {
		loadRefs(db, meaningIDs, entryMap)
	}
}

func loadRefs(db *sqlx.DB, meaningIDs []int, entryMap map[int]*Entry) {
	if len(meaningIDs) == 0 {
		return
	}

	placeholders := strings.Repeat("?,", len(meaningIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := `SELECT r.meaning_id, COALESCE(r.target_hanzi, e.headword) as ref_hanzi 
		FROM refs r 
		LEFT JOIN entries e ON r.target_entry_id = e.id 
		WHERE r.meaning_id IN (` + placeholders + `)`

	args := make([]interface{}, len(meaningIDs))
	for i, id := range meaningIDs {
		args[i] = id
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return
	}
	defer rows.Close()

	refsByMeaning := make(map[int][]string)
	for rows.Next() {
		var meaningID int
		var refHanzi string
		if err := rows.Scan(&meaningID, &refHanzi); err != nil {
			continue
		}
		if refHanzi != "" {
			refsByMeaning[meaningID] = append(refsByMeaning[meaningID], refHanzi)
		}
	}

	for _, e := range entryMap {
		for i := range e.Meanings {
			meaningID := e.Meanings[i].ID
			if refs, ok := refsByMeaning[meaningID]; ok {
				e.Meanings[i].Refs = refs
			}
		}
	}
}
