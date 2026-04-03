package search

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// Search is the main entry point for all search queries
// It determines the search type based on which parameter is provided:
// - hanzi: Chinese characters search (exact match)
// - pinyin: Pinyin search (prefix match)
// - meaning: Russian meaning search (partial match)
// The extended parameter adds prefix matches for hanzi searches
func (s *Searcher) Search(hanzi, pinyin, meaning string, page, limit int, extended bool) (Result, error) {
	if hanzi != "" {
		entries := s.ByHanzi(hanzi, limit, extended)
		return Result{Data: entries, Total: len(entries), Page: page, Limit: limit}, nil
	}
	if pinyin != "" {
		entries, total := s.ByPinyin(pinyin, page, limit)
		return Result{Data: entries, Total: total, Page: page, Limit: limit}, nil
	}
	if meaning != "" {
		entries, total := s.ByMeaning(meaning, page, limit)
		return Result{Data: entries, Total: total, Page: page, Limit: limit}, nil
	}
	return Result{Data: []Entry{}, Total: 0, Page: page, Limit: limit}, nil
}

// scanEntries reads database rows and converts them to Entry structs
// It also loads meanings for each entry
func scanEntries(db *sqlx.DB, rows *sqlx.Rows) []Entry {
	var entryIDs []int
	entryMap := make(map[int]*Entry)

	// First pass: collect all entries
	for rows.Next() {
		var e struct {
			ID     int
			Hanzi  string
			Pinyin string
		}
		if err := rows.StructScan(&e); err != nil {
			continue
		}
		entryMap[e.ID] = &Entry{Hanzi: e.Hanzi, Pinyin: e.Pinyin, Meanings: []Meaning{}}
		entryIDs = append(entryIDs, e.ID)
	}

	// Load meanings for collected entries
	if len(entryIDs) > 0 {
		loadMeanings(db, entryIDs, entryMap)
	}

	// Convert map to slice
	result := make([]Entry, 0, len(entryMap))
	for _, e := range entryMap {
		result = append(result, *e)
	}
	return result
}

// loadMeanings fetches all meanings for the given entry IDs and populates the entry map
func loadMeanings(db *sqlx.DB, entryIDs []int, entryMap map[int]*Entry) {
	if len(entryIDs) == 0 {
		return
	}

	// Build dynamic IN clause: WHERE entry_id IN (?, ?, ?, ...)
	placeholders := strings.Repeat("?,", len(entryIDs))
	placeholders = placeholders[:len(placeholders)-1]
	query := `SELECT id, entry_id, order_num, text FROM meanings WHERE entry_id IN (` + placeholders + `) ORDER BY entry_id, order_num`

	// Convert []int to []interface{} for sqlx
	args := make([]interface{}, len(entryIDs))
	for i, id := range entryIDs {
		args[i] = id
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Printf("loadMeanings error: %v\n", err)
		return
	}
	defer rows.Close()

	// First pass: populate meanings
	meaningIDs := make([]int, 0)
	count := 0
	for rows.Next() {
		var m struct {
			ID       int
			EntryID  int
			OrderNum int
			Text     string
		}
		if err := rows.Scan(&m.ID, &m.EntryID, &m.OrderNum, &m.Text); err != nil {
			fmt.Printf("Scan error: %v\n", err)
			continue
		}
		count++
		if e, ok := entryMap[m.EntryID]; ok {
			e.Meanings = append(e.Meanings, Meaning{ID: m.ID, Index: len(e.Meanings) + 1, Text: m.Text})
			meaningIDs = append(meaningIDs, m.ID)
		}
	}

	// Second pass: load refs for collected meanings
	if len(meaningIDs) > 0 {
		loadRefs(db, meaningIDs, entryMap)
	}
}

// loadRefs fetches refs for the given meaning IDs and populates the entry map
func loadRefs(db *sqlx.DB, meaningIDs []int, entryMap map[int]*Entry) {
	if len(meaningIDs) == 0 {
		return
	}

	placeholders := strings.Repeat("?,", len(meaningIDs))
	placeholders = placeholders[:len(placeholders)-1]
	query := `SELECT r.meaning_id, COALESCE(r.target_hanzi, e.hanzi) as ref_hanzi 
		FROM refs r 
		LEFT JOIN entries e ON r.target_entry_id = e.id 
		WHERE r.meaning_id IN (` + placeholders + `)`

	args := make([]interface{}, len(meaningIDs))
	for i, id := range meaningIDs {
		args[i] = id
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Printf("loadRefs error: %v\n", err)
		return
	}
	defer rows.Close()

	// Group refs by meaning_id
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

	// Attach refs to meanings
	for _, e := range entryMap {
		for i := range e.Meanings {
			meaningID := e.Meanings[i].ID
			if refs, ok := refsByMeaning[meaningID]; ok {
				e.Meanings[i].Refs = refs
			}
		}
	}
}
