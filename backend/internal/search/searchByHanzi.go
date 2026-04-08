package search

import (
	"strings"

	"github.com/jmoiron/sqlx"
)


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

	entires := scanEntries2(s.db, rows)
	
	// Entry not entry ?!?!
	result := Result{
		Data: entires,
		Total: len(entires),
		Page: 1,
		Limit: 10,
	}
	
	return result, nil
}

func scanEntries2(db *sqlx.DB, rows *sqlx.Rows) []Entry {
	var entries []Entry

	// read all rows from SQL result
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.ID, &e.Hanzi, &e.Pinyin); err != nil {
            continue
        }
		e.Meanings = []Meaning{}
		entries = append(entries, e)
	}

	// Проверяем ошибки после цикла
	if err := rows.Err(); err != nil {
		return []Entry{} // возвращаем пустой результат при ошибке
	}

	if len(entries) == 0 {
		return []Entry{}
	}

	loadMeaningsForEntires(db, entries)

	return entries
}

func loadMeaningsForEntires(db *sqlx.DB, entries []Entry) {
	if len(entries) == 0 {
		return
	}

	entryIDs := make([]int, len(entries))
	for i, e := range entries {
		entryIDs[i] = e.ID
	}

	// Создаём плейсхолдеры: ?,?,? ...
	placeholders := strings.Repeat("?,", len(entryIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := `
		SELECT entry_id, text
		FROM meanings
		WHERE entry_id IN (` + placeholders + `)
		ORDER BY entry_id, order_num
	`

	rows, err := db.Query(query, toInterfaceSlice(entryIDs)...)
	if err != nil {
		return
	}
	defer rows.Close()

	meaningMap := make(map[int][]Meaning)

	for rows.Next() {
		var entryID int
        var text string

        
        if err := rows.Scan(&entryID, &text); err != nil {
            continue
        }
        
        meaningMap[entryID] = append(meaningMap[entryID], Meaning{
            Text: text,
            // Index можно добавить позже, если нужно
        })
	}
	
	for i := range entries {
        if meanings, ok := meaningMap[entries[i].ID]; ok {
            entries[i].Meanings = meanings
        }
    }
}

func toInterfaceSlice(ids []int) []interface{} {
	result := make([]interface{}, len(ids))
	for i, id := range ids {
		result[i] = id
	}
	return result
}
