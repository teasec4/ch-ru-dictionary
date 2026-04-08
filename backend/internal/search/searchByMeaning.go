package search

import "strings"

// searchByMeaning — поиск по русскому значению (простая и понятная версия)
func (s *Searcher) searchByMeaning(query string) (Result, error) {
    query = strings.TrimSpace(query)
    if query == "" {
        return Result{Data: []Entry{}, Total: 0, Page: 1, Limit: 50}, nil
    }

    // Приводим запрос к нижнему регистру для регистронезависимого поиска
    query = strings.ToLower(query)
    likePattern := "%" + query + "%"

    // Выбираем только те поля, которые есть в структуре Entry
    rows, err := s.db.Queryx(`
        SELECT e.id, e.headword, e.pinyin
        FROM entries e
        JOIN meanings m ON e.id = m.entry_id
        WHERE LOWER(m.text) LIKE ?
        LIMIT 20
    `, likePattern)

    if err != nil {
        return Result{}, err
    }
    defer rows.Close()

    entries := scanEntries(s.db, rows)
    return Result{Data: entries, Total: len(entries), Page: 1, Limit: 50}, nil
}
