package service

import (
	"context"

	"github.com/gordejka179/test-manager/internal/core"
)

type TestRepository interface {
	AddTest(ctx context.Context, test *core.Test) error
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

type TestService struct {
	repo TestRepository
}

func NewTestService(repo TestRepository) *TestService {
	return &TestService{repo: repo}
}

func (s *TestService) AddTest(ctx context.Context, test *core.Test) error {
	return s.repo.AddTest(ctx, test)
}

func (s *TestService) GetTestByID(ctx context.Context, testId string) (*core.Test, error) {
	return s.repo.GetTestByID(ctx, testId)
}

func (s *TestService) GetAllTests(ctx context.Context) ([]core.Test, error) {
	return s.repo.GetAllTests(ctx)
}

func (s *TestService) DeleteTest(ctx context.Context, id string) error {
	return s.repo.DeleteTest(ctx, id)
}

func (s *TestService) AddConfig(ctx context.Context, testID string, config *core.Config) error {
	return s.repo.AddConfig(ctx, testID, config)
}

func (s *TestService) GetConfigByID(ctx context.Context, testID string, configID string) (*core.Config, error) {
	return s.repo.GetConfigByID(ctx, testID, configID)
}

func (s *TestService) GetAllConfigs(ctx context.Context) ([]core.Config, error) {
	return s.repo.GetAllConfigs(ctx)
}

func (s *TestService) GetAllConfigsToTest(ctx context.Context, testID string) ([]core.Config, error) {
	return s.repo.GetAllConfigsToTest(ctx, testID)
}

func (s *TestService) DeleteConfig(ctx context.Context, testID string) error {
	return s.repo.DeleteConfig(ctx, testID)
}

func (s *TestService) GetLogs(ctx context.Context, testID string, configID string) ([]core.Log, error) {
	return s.repo.GetLogs(ctx, testID, configID)
}
