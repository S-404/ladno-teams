package repository

import (
	"github.com/jmoiron/sqlx"
)

type BaseRepository struct {
	db *sqlx.DB
}

func NewBaseRepository(db *sqlx.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

func (r *BaseRepository) DB() *sqlx.DB {
	return r.db
}
