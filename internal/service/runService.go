package service

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"github.com/gordejka179/test-manager/internal/core"
	"github.com/gordejka179/test-manager/pkg"
)

type RunService struct {
	repo TestRepository
}

func NewRunService(repo TestRepository) *RunService {
	return &RunService{repo: repo}
}

func (s *RunService) RunTest(ctx context.Context, configId int, serverIp string, username string, commandTemplate string) error {
	config, err := s.repo.GetConfigByID(ctx, configId)
	if err != nil {
		log.Fatal("Ошибка метода RunTest:", err)
	}

	testName := config.TestName
	test, err := s.repo.GetTestByName(ctx, testName)

	if err != nil {
		log.Fatal("Ошибка метода RunTest:", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(config.Content, &data); err != nil {
		return err
	}

	pkg.SaveToTOML("tmp.toml", data)

	file, err := os.Create("tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(test.Binary)
	if err != nil {
		log.Fatal("Ошибка метода RunTest", err)
	}

	output := pkg.СonnectSSH(serverIp, username, commandTemplate)

	log := core.Log{Output: output, ConfigID: configId}
	s.repo.AddLog(ctx, &log)
	return nil
}
