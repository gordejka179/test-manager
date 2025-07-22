package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/gordejka179/test-manager/internal/core"
	_ "github.com/mattn/go-sqlite3"
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
			name TEXT NOT NULL PRIMARY KEY,
			config_type TEXT NOT NULL,
			binary BLOB NOT NULL,
			template JSON,
			binary_name TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS test_configs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			test_name TEXT NOT NULL,
			name TEXT NOT NULL,
			config_type TEXT NOT NULL,
			content JSON,
			FOREIGN KEY(test_name) REFERENCES tests(name) ON DELETE CASCADE
		);

		CREATE TABLE IF NOT EXISTS test_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			config_id INTEGER,
			number INTEGER,
			output TEXT,       
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (config_id) REFERENCES test_configs(id) ON DELETE CASCADE
		);
	`)
	return err
}

// Tests

//TODO: проверить был ли такой тест раньше

func (s *SQLiteStorage) AddTest(ctx context.Context, test *core.Test) error {
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO tests (name, config_type, binary, template, binary_name) 
		VALUES (?, ?, ?, ?, ?)`,
		test.Name, test.ConfigType, test.Binary, test.Template, test.BinaryName)

	if err != nil {
		log.Fatalf("Ошибка метода AddTest: %v", err)
	}
	return err
}

func (s *SQLiteStorage) GetTestByName(ctx context.Context, name string) (*core.Test, error) {
	var test core.Test
	err := s.DB.QueryRowContext(ctx,
		`SELECT name, config_type, binary, template, binary_name 
		FROM tests WHERE name = ?`, name).Scan(
		&test.Name, &test.ConfigType, &test.Binary, &test.Template, &test.BinaryName)

	if err != nil {
		log.Fatalf("Ошибка метода GetTestByName: %v", err)
	}
	return &test, err
}

func (s *SQLiteStorage) GetAllTests(ctx context.Context) ([]core.Test, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT name, config_type, binary, template, binary_name FROM tests`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []core.Test
	for rows.Next() {
		var test core.Test
		if err := rows.Scan(&test.Name, &test.ConfigType, &test.Binary, &test.Template, &test.BinaryName); err != nil {
			return nil, err
		}
		tests = append(tests, test)
		if err != nil {
			log.Fatal("Ошибка метода GetAllTests: ", err)
		}
	}

	return tests, nil
}

func (s *SQLiteStorage) DeleteTest(ctx context.Context, name string) error {
	_, err := s.DB.ExecContext(ctx,
		`DELETE FROM tests WHERE name = ?`, name)
	if err != nil {
		return err
	}

	return nil

}

// Configs
func (s *SQLiteStorage) AddConfig(ctx context.Context, config *core.Config) (int64, error) {
	result, err := s.DB.ExecContext(ctx,
		`INSERT INTO test_configs (test_name, name, config_type, content)
			VALUES (?, ?, ?, ?)`,
		config.TestName, config.Name, config.ConfigType, config.Content)

	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

func (s *SQLiteStorage) GetConfigByID(ctx context.Context, configID int) (*core.Config, error) {
	var config core.Config
	err := s.DB.QueryRowContext(ctx,
		`SELECT id, test_name, name, config_type, content
		FROM test_configs WHERE id = ?`,
		configID).Scan(
		&config.ID, &config.TestName, &config.Name, &config.ConfigType, &config.Content)

	return &config, err
}

func (s *SQLiteStorage) GetAllConfigs(ctx context.Context) ([]core.Config, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, test_name, name, config_type, content
		FROM test_configs`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []core.Config
	for rows.Next() {
		var config core.Config
		if err := rows.Scan(
			&config.ID, &config.TestName, &config.Name, &config.ConfigType, &config.Content,
		); err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func (s *SQLiteStorage) GetAllConfigsToTest(ctx context.Context, testName string) ([]core.Config, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, test_name, name, config_type, content
		FROM test_configs WHERE test_name = ?`,
		testName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var configs []core.Config
	for rows.Next() {
		var config core.Config
		if err := rows.Scan(
			&config.ID, &config.TestName, &config.Name, &config.ConfigType, &config.Content,
		); err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil

}

func (s *SQLiteStorage) DeleteConfig(ctx context.Context, id string) error {
	_, err := s.DB.ExecContext(ctx,
		`DELETE FROM test_configs WHERE id = ?`, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *SQLiteStorage) GetLogsToConfig(ctx context.Context, configID int) ([]core.Log, error) {
	rows, err := s.DB.QueryContext(ctx,
		`SELECT id, config_id, created_at, output
		FROM test_logs WHERE config_id = ?`,
		configID)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []core.Log
	for rows.Next() {
		var log core.Log
		if err := rows.Scan(
			&log.ID, &log.ConfigID, &log.CreatedAt, &log.Output,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	return logs, nil
}

func (s *SQLiteStorage) AddLog(ctx context.Context, log *core.Log) error {
	//считаем, сколько раз логов относится к конфигу
	num := 0

	row := s.DB.QueryRowContext(ctx, `
	SELECT MAX(number) FROM test_logs WHERE config_id = ?`, log.ConfigID)
	if row == nil {
		num = 0
	} else {
		row.Scan(&num)
	}

	num += 1
	_, err := s.DB.ExecContext(ctx,
		`INSERT INTO test_logs (config_id, number, output)
			VALUES (?, ?, ?)`,
		log.ConfigID, num, log.Output)

	if err != nil {
		return err
	}
	return nil
}
