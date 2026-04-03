package search

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type Entry struct {
	Hanzi    string    `json:"hanzi"`
	Pinyin   string    `json:"pinyin,omitempty"`
	Meanings []Meaning `json:"meanings"`
}

type Meaning struct {
	Index int    `json:"index"`
	Text  string `json:"text"`
}

type Result struct {
	Data  []Entry `json:"data"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
}

type Searcher struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Searcher {
	return &Searcher{db: db}
}

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

func (s *Searcher) ByHanzi(hanzi string, limit int, extended bool) []Entry {
	// Всегда точное совпадение
	exact := s.byHanziExact(hanzi, limit)

	if !extended {
		return exact
	}

	// Extended: добавляем prefix если есть место
	if len(exact) < limit {
		remaining := limit - len(exact)
		prefix := s.byHanziPrefix(hanzi, remaining)

		// Убираем дубликаты
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

func scanEntries(db *sqlx.DB, rows *sqlx.Rows) []Entry {
	var entryIDs []int
	entryMap := make(map[int]*Entry)

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

	if len(entryIDs) > 0 {
		loadMeanings(db, entryIDs, entryMap)
	}

	result := make([]Entry, 0, len(entryMap))
	for _, e := range entryMap {
		result = append(result, *e)
	}
	return result
}

func loadMeanings(db *sqlx.DB, entryIDs []int, entryMap map[int]*Entry) {
	if len(entryIDs) == 0 {
		return
	}

	placeholders := strings.Repeat("?,", len(entryIDs))
	placeholders = placeholders[:len(placeholders)-1]
	query := `SELECT entry_id, order_num, text FROM meanings WHERE entry_id IN (` + placeholders + `) ORDER BY entry_id, order_num`

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

	for rows.Next() {
		var m struct {
			EntryID  int
			OrderNum int
			Text     string
		}
		if err := rows.Scan(&m.EntryID, &m.OrderNum, &m.Text); err != nil {
			continue
		}
		if e, ok := entryMap[m.EntryID]; ok {
			e.Meanings = append(e.Meanings, Meaning{Index: len(e.Meanings) + 1, Text: m.Text})
		}
	}
}
