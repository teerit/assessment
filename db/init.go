package db

import (
	"database/sql"
	"fmt"
	"os"
)

func InitDB() (*sql.DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		title TEXT,
		amount FLOAT,
		note TEXT,
		tags TEXT[]
	);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		return nil, fmt.Errorf("can't create table: %w", err)
	}

	return db, nil
}
