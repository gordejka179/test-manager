package storage

import (
	"context"

	"github.com/gordejka179/test-manager/internal/core"
)

type Storage interface {
	// Tests
	CreateTest(ctx context.Context, test core.Test) error
	GetTest(ctx context.Context, id string) (core.Test, error)
	ListTests(ctx context.Context) ([]core.Test, error)
	UpdateTest(ctx context.Context, id string, test core.Test) error
	DeleteTest(ctx context.Context, id string) error

	// Configs
	CreateConfig(ctx context.Context, config core.Config) error

	//Log
	GetLogs(ctx context.Context, testID string, configID string) error
}
