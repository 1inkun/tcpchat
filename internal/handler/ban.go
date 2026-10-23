package handler

import (
	"database/sql"

	"github.com/redis/go-redis/v9"
)

type BanHandler struct {
	rdb *redis.Client
	db  *sql.DB
}

func NewBanHandler(rdb *redis.Client, db *sql.DB) *BanHandler {
	return &BanHandler{
		rdb: rdb,
		db:  db,
	}
}
