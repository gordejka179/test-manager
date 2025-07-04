package service

import (
	"context"

	"github.com/gordejka179/test-manager/internal/core"
)

type TestRepository interface {
	CreateTest(ctx context.Context, test *core.Test) error
	GetByID(ctx context.Context, testID string) (*core.Test, error)
	GetAll(ctx context.Context) ([]*core.Test, error)
	DeleteTest(ctx context.Context, id string) error
	AddConfig(ctx context.Context, testID string, config *core.Config) error
	DeleteConfig(ctx context.Context, testID string, configID string) error
	GetConfigByID(ctx context.Context, testID string, configID string) (*core.Config, error)
}

type TestRunner interface {
	Run(test *core.Test, configName string) (*core.log, error)
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
	return s.repo.GetByID(ctx, testId)
}

func (s *Service) GetAll(ctx context.Context) ([]*core.Test, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) DeleteTest(ctx context.Context, id string) error {
	return s.repo.DeleteTest(ctx, id)
}

func (s *Service) AddConfig(ctx context.Context, testID string, config *core.Config) error {
	return s.repo.AddConfig(ctx, testID, config)
}

func (s *Service) DeleteConfig(ctx context.Context, testID string, configID string) error {
	return s.repo.DeleteConfig(ctx, testID, configID)
}

func (s *Service) GetConfigByID(ctx context.Context, testID string, configID string) (*core.Config, error) {
	return s.repo.GetConfigByID(ctx, testID, configID)
}
