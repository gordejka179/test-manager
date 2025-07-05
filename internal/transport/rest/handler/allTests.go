package handler

import (
	"context"

	"github.com/gordejka179/test-manager/internal/core"
)

type TestRepository interface {
	CreateTest(ctx context.Context, test *core.Test) error
	GetTestByID(ctx context.Context, testID string) (*core.Test, error)
	GetAllTests(ctx context.Context) ([]core.Test, error)
	DeleteTest(ctx context.Context, id string) error
	AddConfig(ctx context.Context, testID string, config *core.Config) error
	GetConfigByID(ctx context.Context, testID string, configID string) (*core.Config, error)
	GetAllConfigs(ctx context.Context) ([]core.Config, error)
	GetAllConfigsToTest(ctx context.Context, testID string) ([]core.Config, error)
	DeleteConfig(ctx context.Context, testID string) error
	GetLogs(ctx context.Context, testID string, configID string) ([]core.Log, error)
}

type TestRunner interface {
	Run(test *core.Test, configName string) (*core.Log, error)
}

type Service struct {
	repo   TestRepository
	runner TestRunner
}

type ServiceHandler struct {
	service Service
}

func NewServiceHandler(S Service) *ServiceHandler {
	return &ServiceHandler{service: S}
}
