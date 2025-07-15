package storage

import (
	"fmt"

	"github.com/gordejka179/test-manager/internal/core"
)

type Runner struct {
}

func NewRunner() *Runner {
	return &Runner{}
}

func (r *Runner) Run(configId string) (*core.Log, error) {
	fmt.Print("Hello")
	return nil, nil
}
