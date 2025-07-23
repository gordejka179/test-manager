package service

import (
	"context"
	"fmt"

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

type RepService struct {
	repo TestRepository
}

func NewRepService(repo TestRepository) *RepService {
	return &RepService{repo: repo}
}

func (s *RepService) AddTest(ctx context.Context, test *core.Test) error {
	existingTest, err := s.repo.GetTestByName(ctx, test.Name)
	if err != nil {
		return fmt.Errorf("ошибка проверки теста: %v", err)
	}
	if existingTest != nil {
		return fmt.Errorf("тест с таким именем уже есть")
	}
	return s.repo.AddTest(ctx, test)
}

func (s *RepService) GetTestByName(ctx context.Context, name string) (*core.Test, error) {
	return s.repo.GetTestByName(ctx, name)
}

func (s *RepService) GetAllTests(ctx context.Context) ([]core.Test, error) {
	return s.repo.GetAllTests(ctx)
}

func (s *RepService) DeleteTest(ctx context.Context, name string) error {
	return s.repo.DeleteTest(ctx, name)
}

func (s *RepService) AddConfig(ctx context.Context, config *core.Config) (int64, error) {
	return s.repo.AddConfig(ctx, config)
}

func (s *RepService) GetConfigByID(ctx context.Context, configID int) (*core.Config, error) {
	return s.repo.GetConfigByID(ctx, configID)
}

func (s *RepService) GetAllConfigs(ctx context.Context) ([]core.Config, error) {
	return s.repo.GetAllConfigs(ctx)
}

func (s *RepService) GetAllConfigsToTest(ctx context.Context, testID string) ([]core.Config, error) {
	return s.repo.GetAllConfigsToTest(ctx, testID)
}

func (s *RepService) DeleteConfig(ctx context.Context, testID string) error {
	return s.repo.DeleteConfig(ctx, testID)
}

func (s *RepService) GetLogsToConfig(ctx context.Context, configID int) ([]core.Log, error) {
	return s.repo.GetLogsToConfig(ctx, configID)
}

func (s *RepService) AddLog(ctx context.Context, log *core.Log) error {
	return s.repo.AddLog(ctx, log)
}
