package service

import (
	"github.com/jmoiron/sqlx"
)

type Result struct{
	Hanzi string `json:"chinese"`
	Pinyin string `json:"pinyin"`
	Meaning string `json:"meanings"`
}

type DbService struct {
	db *sqlx.DB
}

func NewDbService(db *sqlx.DB) *DbService {
	return &DbService{
		db: db,
	}
}

func (s *DbService)Searh(query string)([]Result, error){
	if containsHanzi(query){
		return s.SearchByHanzi(query)
	}
	
	if isLatin(query){
		return s.SearchByPinyin(query)
	}
	
	return s.SearchByMeaning(query)
}

func (s *DbService) SearchByHanzi(query string) ([]Result, error){
	rows, err := s.db.Query(`
			SELECT w.hanzi, w.pinyin, GROUP_CONCAT(m.meaning, '; ')
			FROM words w
			JOIN meanings m ON w.id = m.word_id
			WHERE w.hanzi LIKE ?
			GROUP BY w.id
			LIMIT 20;
		`, query+"%")
	
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()
	
	var results []Result
	
	for rows.Next(){
		var r Result
		if err := rows.Scan(&r.Hanzi, &r.Pinyin, &r.Meaning); err != nil{
			return nil, err
		}
		results = append(results, r)
	}
	
	return results, nil
}

func (s *DbService) SearchByPinyin(query string) ([]Result, error) {
	rows, err := s.db.Query(`
		SELECT w.hanzi, w.pinyin, GROUP_CONCAT(m.meaning, '; ')
		FROM words w
		JOIN meanings m ON w.id = m.word_id
		WHERE w.pinyin LIKE ?
		GROUP BY w.id
		LIMIT 20;
	`, query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Result

	for rows.Next() {
		var r Result
		if err := rows.Scan(&r.Hanzi, &r.Pinyin, &r.Meaning); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

func (s *DbService)SearchByMeaning(query string) ([]Result, error) {
	rows, err := s.db.Query(`
		SELECT w.hanzi, w.pinyin, GROUP_CONCAT(m.meaning, '; ')
		FROM words w
		JOIN meanings m ON w.id = m.word_id
		WHERE m.meaning LIKE ?
		GROUP BY w.id
		LIMIT 20;
	`, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []Result

	for rows.Next() {
		var r Result
		if err := rows.Scan(&r.Hanzi, &r.Pinyin, &r.Meaning); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

// helpers
func containsHanzi(s string) bool {
	for _, r := range s {
		if r >= 0x4E00 && r <= 0x9FFF {
			return true
		}
	}
	return false
}

func isLatin(s string) bool {
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			return true
		}
	}
	return false
}

