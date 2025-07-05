package service

import (
	"github.com/gordejka179/test-manager/internal/core"
)

type TestRunner interface {
	Run(test *core.Test, configName string) (*core.Log, error)
}
