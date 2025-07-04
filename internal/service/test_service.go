package service

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
	GetLogs(ctx context.Context, testID string, configID string) error
}

type TestRunner interface {
	Run(test *core.Test, configName string) (*core.Log, error)
}

type Service struct {
	repo   TestRepository
	runner TestRunner
}

func NewService(repo TestRepository, runner TestRunner) *Service {
	return &Service{repo: repo, runner: runner}
}

func (s *Service) CreateTest(ctx context.Context, test *core.Test) error {
	return s.repo.CreateTest(ctx, test)
}

func (s *Service) GetByID(ctx context.Context, testId string) (*core.Test, error) {
	return s.repo.GetTestByID(ctx, testId)
}

func (s *Service) GetAll(ctx context.Context) ([]core.Test, error) {
	return s.repo.GetAllTests(ctx)
}

func (s *Service) DeleteTest(ctx context.Context, id string) error {
	return s.repo.DeleteTest(ctx, id)
}

func (s *Service) AddConfig(ctx context.Context, testID string, config *core.Config) error {
	return s.repo.AddConfig(ctx, testID, config)
}

func (s *Service) DeleteConfig(ctx context.Context, testID string) error {
	return s.repo.DeleteConfig(ctx, testID)
}

func (s *Service) GetConfigByID(ctx context.Context, testID string, configID string) (*core.Config, error) {
	return s.repo.GetConfigByID(ctx, testID, configID)
}
