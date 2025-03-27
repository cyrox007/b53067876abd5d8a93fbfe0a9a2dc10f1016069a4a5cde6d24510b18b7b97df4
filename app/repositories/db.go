package repositories

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver
	"gopkg.in/reform.v1"
)

type DB struct {
	*reform.DB
}

func NewDB(dsn string) (*DB, error) {
	sqlDB, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Настройка пула соединений
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetConnMaxLifetime(0)

	// Создание реформ-базы данных
	db := reform.NewDB(sqlDB, nil, nil)
	return &DB{db}, nil
}
