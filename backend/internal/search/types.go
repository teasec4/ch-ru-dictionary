package search

import "github.com/jmoiron/sqlx"

// Entry represents a dictionary entry with hanzi, pinyin and meanings
type Entry struct {
	Hanzi    string    `json:"hanzi"`
	Pinyin   string    `json:"pinyin,omitempty"`
	Meanings []Meaning `json:"meanings"`
}

// Meaning represents a single meaning/translation of an entry
type Meaning struct {
	ID    int      `json:"-"`
	Index int      `json:"index"`
	Text  string   `json:"text"`
	Refs  []string `json:"refs,omitempty"`
}

// Result represents the search result with pagination info
type Result struct {
	Data  []Entry `json:"data"`
	Total int     `json:"total"`
	Page  int     `json:"page"`
	Limit int     `json:"limit"`
}

// Searcher handles all search operations against the dictionary database
type Searcher struct {
	db *sqlx.DB
}

// New creates a new Searcher instance with the given database connection
func New(db *sqlx.DB) *Searcher {
	return &Searcher{db: db}
}
