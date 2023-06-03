package storage

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func dbSession(postgresDSN string) (*sql.DB, error) {
	db, err := sql.Open("pgx", postgresDSN)
	if err != nil {
		return nil, err
	}
	return db, nil
}
