package search

import (
	"strings"

	"github.com/jmoiron/sqlx"
)

// scanEntries — простая и понятная версия
// Берёт строки из базы и собирает []Entry с значениями (meanings)
func scanEntries(db *sqlx.DB, rows *sqlx.Rows) []Entry {
	var entries []Entry

	// 1. Читаем основные данные слов (id, hanzi, pinyin, frequency)
	for rows.Next() {
		var e Entry
		
		err := rows.Scan(&e.ID, &e.Hanzi, &e.Pinyin)
		if err != nil {
			continue
		}

		e.Meanings = []Meaning{} // пока пустой слайс
		entries = append(entries, e)
	}

	if err := rows.Err(); err != nil {
		return []Entry{}
	}

	if len(entries) == 0 {
		return []Entry{}
	}

	// 2. Загружаем русские значения для всех найденных слов
	loadMeaningsSimple(db, entries)

	return entries
}


func loadMeaningsSimple(db *sqlx.DB, entries []Entry) {
	if len(entries) == 0 {
		return
	}

	// Собираем все ID слов
	entryIDs := make([]int, len(entries))
	for i, e := range entries {
		entryIDs[i] = e.ID
	}

	// Создаём плейсхолдеры для IN (?, ?, ?)
	placeholders := strings.Repeat("?,", len(entryIDs))
	placeholders = placeholders[:len(placeholders)-1]

	query := `
		SELECT entry_id, id, order_num, text 
		FROM meanings 
		WHERE entry_id IN (` + placeholders + `)
		ORDER BY entry_id, order_num ASC
	`

	rows, err := db.Query(query, toInterfaceSlice(entryIDs)...)
	if err != nil {
		return
	}
	defer rows.Close()

	// Группируем значения по ID слова
	meaningMap := make(map[int][]Meaning)

	for rows.Next() {
		var entryID int
		var m Meaning

		err := rows.Scan(&entryID, &m.ID, &m.Index, &m.Text)
		if err != nil {
			continue
		}

		meaningMap[entryID] = append(meaningMap[entryID], m)
	}

	// Присваиваем значения каждому Entry
	for i := range entries {
		if meanings, ok := meaningMap[entries[i].ID]; ok {
			entries[i].Meanings = meanings
		}
	}
}

// toInterfaceSlice помогает использовать []int в запросе с IN (?,?,?)
func toInterfaceSlice(ids []int) []interface{} {
	result := make([]interface{}, len(ids))
	for i, id := range ids {
		result[i] = id
	}
	return result
}