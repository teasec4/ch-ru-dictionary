package service

import (
	"backend/internal/model"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type DbService struct {
	db *sqlx.DB
}

func NewDbService(db *sqlx.DB) *DbService {
	return &DbService{
		db: db,
	}
}

func (s *DbService) SearchWords(search string) ([]model.Response, error) {
	// Определяем тип запроса
	queryType := detectQueryType(search)

	var results []struct {
		Hanzi    string `db:"hanzi"`
		Pinyin   string `db:"pinyin"`
		Meanings string `db:"meanings"`
	}

	var err error

	switch queryType {
	case "chinese":
		err = s.searchChinese(search, &results)
	case "pinyin":
		err = s.searchPinyin(search, &results)
	default:
		err = s.searchMeaning(search, &results)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to search words: %w", err)
	}

	// Преобразуем результаты в Response
	var responses []model.Response
	for _, r := range results {
		meaningsList := splitMeanings(r.Meanings)

		responses = append(responses, model.Response{
			Chinese:          r.Hanzi,
			Pinyin:           r.Pinyin,
			PinyinNormalized: normalizePinyin(r.Pinyin),
			Meanings:         meaningsList,
		})
	}

	return responses, nil
}

// detectQueryType определяет тип поискового запроса
func detectQueryType(query string) string {
	// Проверяем китайские иероглифы
	for _, r := range query {
		if (r >= 0x4E00 && r <= 0x9FFF) || // Основные иероглифы
			(r >= 0x3400 && r <= 0x4DBF) || // Расширение A
			(r >= 0x20000 && r <= 0x2A6DF) { // Расширение B
			return "chinese"
		}
	}

	// Проверяем пиньинь (содержит латинские буквы)
	hasLatin := false
	for _, r := range query {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLatin = true
			break
		}
	}

	if hasLatin {
		return "pinyin"
	}

	// По умолчанию - поиск по переводу
	return "meaning"
}

// searchChinese ищет по китайским иероглифам
func (s *DbService) searchChinese(search string, results *[]struct {
	Hanzi    string `db:"hanzi"`
	Pinyin   string `db:"pinyin"`
	Meanings string `db:"meanings"`
}) error {
	query := `
		SELECT w.hanzi, w.pinyin, GROUP_CONCAT(m.meaning, '; ') as meanings
		FROM words w
		LEFT JOIN meanings m ON w.id = m.word_id
		WHERE w.hanzi = ?              -- exact match
		   OR w.hanzi LIKE ?           -- prefix match
		   OR w.hanzi LIKE ?           -- fuzzy match
		GROUP BY w.id
		ORDER BY
			CASE
				WHEN w.hanzi = ? THEN 1
				WHEN w.hanzi LIKE ? THEN 2
				ELSE 3
			END,
			w.frequency DESC,
			LENGTH(w.hanzi) ASC,
			w.hanzi ASC
		LIMIT 50`

	exact := search
	prefix := search + "%"
	fuzzy := "%" + search + "%"

	return s.db.Select(results, query,
		exact, prefix, fuzzy, // WHERE conditions
		exact, prefix) // ORDER BY conditions
}

// searchPinyin ищет по пиньиню
func (s *DbService) searchPinyin(search string, results *[]struct {
	Hanzi    string `db:"hanzi"`
	Pinyin   string `db:"pinyin"`
	Meanings string `db:"meanings"`
}) error {
	// Нормализуем пиньинь для поиска (убираем тоны)
	normalizedSearch := normalizePinyin(search)

	query := `
		SELECT w.hanzi, w.pinyin, GROUP_CONCAT(m.meaning, '; ') as meanings
		FROM words w
		LEFT JOIN meanings m ON w.id = m.word_id
		WHERE w.pinyin = ?              -- exact pinyin match
		   OR w.pinyin LIKE ?           -- prefix pinyin match
		   OR w.pinyin LIKE ?           -- fuzzy pinyin match
		   OR w.pinyin_norm LIKE ?      -- normalized pinyin match
		GROUP BY w.id
		ORDER BY
			CASE
				WHEN w.pinyin = ? THEN 1
				WHEN w.pinyin LIKE ? THEN 2
				WHEN w.pinyin_norm LIKE ? THEN 3
				ELSE 4
			END,
			w.frequency DESC,
			LENGTH(w.hanzi) ASC,
			w.hanzi ASC
		LIMIT 50`

	exact := search
	prefix := search + "%"
	fuzzy := "%" + search + "%"
	normalizedFuzzy := "%" + normalizedSearch + "%"

	return s.db.Select(results, query,
		exact, prefix, fuzzy, normalizedFuzzy, // WHERE conditions
		exact, prefix, normalizedFuzzy) // ORDER BY conditions
}

// searchMeaning ищет по переводу
func (s *DbService) searchMeaning(search string, results *[]struct {
	Hanzi    string `db:"hanzi"`
	Pinyin   string `db:"pinyin"`
	Meanings string `db:"meanings"`
}) error {
	query := `
		SELECT w.hanzi, w.pinyin, GROUP_CONCAT(m.meaning, '; ') as meanings
		FROM words w
		LEFT JOIN meanings m ON w.id = m.word_id
		WHERE w.id IN (
			SELECT word_id FROM meanings WHERE meaning = ?              -- exact match
			UNION
			SELECT word_id FROM meanings WHERE meaning LIKE ? || ',%'   -- first in list
			UNION
			SELECT word_id FROM meanings WHERE meaning LIKE ? || ';%'   -- first in list (alt)
			UNION
			SELECT word_id FROM meanings WHERE meaning LIKE '% ' || ? || ' %'  -- word in middle
			UNION
			SELECT word_id FROM meanings WHERE meaning LIKE ? || ' %'   -- word at start
			UNION
			SELECT word_id FROM meanings WHERE meaning LIKE '% ' || ?   -- word at end
			UNION
			SELECT word_id FROM meanings WHERE meaning LIKE '%' || ? || '%'    -- fuzzy match
		)
		GROUP BY w.id
		ORDER BY
			CASE
				WHEN w.id IN (SELECT word_id FROM meanings WHERE meaning = ?) THEN 1
				WHEN w.id IN (SELECT word_id FROM meanings WHERE meaning LIKE ? || ',%') THEN 2
				WHEN w.id IN (SELECT word_id FROM meanings WHERE meaning LIKE ? || ';%') THEN 3
				WHEN w.id IN (SELECT word_id FROM meanings WHERE meaning LIKE '% ' || ? || ' %') THEN 4
				WHEN w.id IN (SELECT word_id FROM meanings WHERE meaning LIKE ? || ' %') THEN 5
				WHEN w.id IN (SELECT word_id FROM meanings WHERE meaning LIKE '% ' || ?) THEN 6
				ELSE 7
			END,
			w.frequency DESC,
			LENGTH(w.hanzi) ASC,
			w.hanzi ASC
		LIMIT 50`

	return s.db.Select(results, query,
		search, search, search, search, search, search, search, // WHERE conditions
		search, search, search, search, search, search) // ORDER BY conditions
}

func normalizePinyin(pinyin string) string {
	if pinyin == "" {
		return ""
	}

	// Убираем тоны (цифры 1-5) и пробелы
	var result strings.Builder
	for _, r := range pinyin {
		if !(r >= '1' && r <= '5') && r != ' ' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

func splitMeanings(meanings string) []string {
	if meanings == "" {
		return []string{}
	}

	// Разделяем по "; " и очищаем
	var result []string
	parts := strings.Split(meanings, "; ")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}

	return result
}
