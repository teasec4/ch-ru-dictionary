package service

import (
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

// do search logic

