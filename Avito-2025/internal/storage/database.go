package storage

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func NewDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}
	// Проверяем соединение
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
