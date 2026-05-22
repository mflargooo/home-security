package store

import "github.com/jackc/pgx/v5/pgxpool"

type ClipStore struct {
	db *pgxpool.Pool
}

func NewClipStore(db *pgxpool.Pool) *ClipStore {
	return &ClipStore{db: db}
}
