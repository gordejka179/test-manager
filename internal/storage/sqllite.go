package storage

import (
	"database/sql"
	"fmt"
)

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("DB ping failed: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS tests (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			binary BLOB NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS test_configs (
			id TEXT PRIMARY KEY,
			test_id TEXT NOT NULL,
			name TEXT NOT NULL,
			config TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(test_id) REFERENCES tests(id) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS test_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			test_id TEXT NOT NULL,
			config_id TEXT NOT NULL,
			output TEXT,         
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			
			FOREIGN KEY (test_id) REFERENCES tests(id) ON DELETE CASCADE,

			FOREIGN KEY (config_id) REFERENCES test_configs(id) ON DELETE CASCADE,
		);
	`)
	return err
}
