package service

import (
	"fmt"

	"github.com/gordejka179/test-manager/internal/core"
)

type TestRunner interface {
	Run(configId string) (*core.Log, error)
}

type RunService struct {
	runner TestRunner
}

func NewRunService(runner TestRunner) *RunService {
	return &RunService{runner: runner}
}

func (s *RunService) Run(configId string) (*core.Log, error) {
	fmt.Println("Run")
	log := core.Log{}
	return &log, nil
}
