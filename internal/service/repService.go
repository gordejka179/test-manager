package service

import (
	"context"

	"github.com/gordejka179/test-manager/internal/core"
)

type TestRepository interface {
	AddTest(ctx context.Context, test *core.Test) error
	GetTestByName(ctx context.Context, name string) (*core.Test, error)
	GetAllTests(ctx context.Context) ([]core.Test, error)
	DeleteTest(ctx context.Context, name string) error
	AddConfig(ctx context.Context, config *core.Config) (int64, error)
	GetConfigByID(ctx context.Context, configID int) (*core.Config, error)
	GetAllConfigs(ctx context.Context) ([]core.Config, error)
	GetAllConfigsToTest(ctx context.Context, testName string) ([]core.Config, error)
	DeleteConfig(ctx context.Context, id string) error
	GetLogsToConfig(ctx context.Context, configID int) ([]core.Log, error)
	AddLog(ctx context.Context, log *core.Log) error
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

func (s *TestService) GetTestByName(ctx context.Context, name string) (*core.Test, error) {
	return s.repo.GetTestByName(ctx, name)
}

func (s *TestService) GetAllTests(ctx context.Context) ([]core.Test, error) {
	return s.repo.GetAllTests(ctx)
}

func (s *TestService) DeleteTest(ctx context.Context, name string) error {
	return s.repo.DeleteTest(ctx, name)
}

func (s *TestService) AddConfig(ctx context.Context, config *core.Config) (int64, error) {
	return s.repo.AddConfig(ctx, config)
}

func (s *TestService) GetConfigByID(ctx context.Context, configID int) (*core.Config, error) {
	return s.repo.GetConfigByID(ctx, configID)
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

func (s *TestService) GetLogsToConfig(ctx context.Context, configID int) ([]core.Log, error) {
	return s.repo.GetLogsToConfig(ctx, configID)
}

func (s *TestService) AddLog(ctx context.Context, log *core.Log) error {
	return s.repo.AddLog(ctx, log)
}
