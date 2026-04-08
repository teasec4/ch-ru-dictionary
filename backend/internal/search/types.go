package search

import "github.com/jmoiron/sqlx"

// added an id for consistnecy with DB 
type Entry struct {
	ID int 				`json:"id"`
	Hanzi    string    `json:"hanzi"`
	Pinyin   string    `json:"pinyin,omitempty"`
	Meanings []Meaning `json:"meanings"`
}

type Meaning struct {
	ID    int      `json:"-"`
	Index int      `json:"index"`
	Text  string   `json:"text"`
	Refs  []string `json:"refs,omitempty"`
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
