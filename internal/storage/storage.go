package storage

import (
	"context"

	"github.com/gordejka179/test-manager/internal/core"
)

type Storage interface {
	// Tests
	CreateTest(ctx context.Context, test *core.Test) error
	GetTestByID(ctx context.Context, testID string) (*core.Test, error)
	GetAllTests(ctx context.Context) ([]core.Test, error)
	DeleteTest(ctx context.Context, id string) error

	// Configs
	AddConfig(ctx context.Context, testID string, config *core.Config) error
	GetConfigByID(ctx context.Context, testID string, configID string) (*core.Config, error)
	GetAllConfigs(ctx context.Context) ([]core.Config, error)
	GetAllConfigsToTest(ctx context.Context, testID string) ([]core.Config, error)
	DeleteConfig(ctx context.Context, testID string) error

	//Log
	GetLogs(ctx context.Context, testID string, configID string) error
}
