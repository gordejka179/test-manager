package storage

import (
	"database/sql"
	"fmt"

	"github.com/gordejka179/test-manager/internal/core"
)

type Runner struct {
	DB *sql.DB
}

func NewRunner(db *sql.DB) *Runner {
	return &Runner{DB: db}
}

func (r *Runner) Run(test *core.Test, configName string) (*core.Log, error) {
	fmt.Print("Hello")
	return nil, nil
}
