package service

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sqlx.DB
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{db: db}
}

func (d *DB) InitFTS() error {
	_, err := d.db.Exec(`
		CREATE VIRTUAL TABLE IF NOT EXISTS meanings_fts USING fts5(
			text,
			entry_id UNINDEXED,
			content=meanings,
			content_rowid=id
		);
	`)
	if err != nil {
		return err
	}

	_, err = d.db.Exec(`
		INSERT OR IGNORE INTO meanings_fts(rowid, text, entry_id)
		SELECT id, text, entry_id FROM meanings;
	`)
	return err
}
