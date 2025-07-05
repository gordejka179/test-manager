package storage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/gordejka179/test-manager/internal/core"
)

type SQLiteStorage struct {
	DB *sql.DB
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

	return &SQLiteStorage{DB: db}, nil
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

//TODO: везде обработка ошибок

// Tests
func (s *SQLiteStorage) AddTest(ctx context.Context, test *core.Test) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO tests (id, name, binary, created_at) 
		VALUES (?, ?, ?, ?)`,
		test.ID, test.Name, test.Binary, test.CreatedAt)
	return err
}

func (s *SQLiteStorage) GetTestByID(ctx context.Context, id string) (*core.Test, error) {
	var test core.Test
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, name, binary, created_at 
		FROM tests WHERE id = ?`, id).Scan(
		&test.ID, &test.Name, &test.Binary, &test.CreatedAt)

	// TODO: Обработка ошибок
	return &test, err
}

func (s *SQLiteStorage) GetAllTests(ctx context.Context) ([]core.Test, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, name, description, created_at, updated_at FROM tests`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []core.Test
	for rows.Next() {
		var test core.Test
		if err := rows.Scan(
			&test.ID, &test.Name, &test.Binary, &test.CreatedAt); err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}

	return tests, nil
}

func (s *SQLiteStorage) DeleteTest(ctx context.Context, testID string) error {
	_, err := s.DB.ExecContext(ctx,
		`DELETE FROM tests WHERE id = ?`, testID)
	if err != nil {
		return err
	}

	// TODO: обработка ошибок

	return nil

}

// Configs

func (s *SQLiteStorage) AddConfig(ctx context.Context, testID string, config *core.Config) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO test_configs (id, test_id, name, config, created_at)
		VALUES (?, ?, ?, ?, ?)`,
		config.ID, config.TestID, config.Name, config.Config, config.CreatedAt)
	return err
}

func (s *SQLiteStorage) GetConfigByID(ctx context.Context, testID string, configID string) (*core.Config, error) {
	var config core.Config
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, test_id, name, config, created_at
		FROM test_configs WHERE id = ?`,
		configID).Scan(
		&config.ID, &config.TestID, &config.Name, &config.Config, &config.CreatedAt)

	return &config, err
}

func (s *SQLiteStorage) GetAllConfigs(ctx context.Context) ([]core.Config, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, test_id, name, config, created_at
		FROM test_configs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []core.Config
	for rows.Next() {
		var config core.Config
		if err := rows.Scan(
			&config.ID, &config.TestID, &config.Name, &config.Config, &config.CreatedAt,
		); err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (s *SQLiteStorage) GetAllConfigsToTest(ctx context.Context, testID string) ([]core.Config, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, test_id, name, config, created_at
		FROM test_configs WHERE test_id = ?`,
		testID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []core.Config
	for rows.Next() {
		var config core.Config
		if err := rows.Scan(
			&config.ID, &config.TestID, &config.Name, &config.Config, &config.CreatedAt,
		); err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil

}

func (s *SQLiteStorage) DeleteConfig(ctx context.Context, testID string) error {
	_, err := s.DB.ExecContext(ctx,
		`DELETE FROM test_configs WHERE id = ?`, testID)
	if err != nil {
		return err
	}

	// TODO: обработка ошибок

	return nil
}

func (s *SQLiteStorage) GetLogs(ctx context.Context, testID string, configID string) ([]core.Log, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, test_id, config_id, output, created_at
		FROM test_logs WHERE id = ? AND test_id =Y`,
		configID, testID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []core.Log
	for rows.Next() {
		var log core.Log
		if err := rows.Scan(
			&log.ID, &log.TestID, &log.ConfigID, &log.Output, &log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}
