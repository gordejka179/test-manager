package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

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
		return fmt.Errorf("ошибка метода RunTest: %v", err)
	}

	testName := config.TestName
	test, err := s.repo.GetTestByName(ctx, testName)

	if err != nil {
		return fmt.Errorf("ошибка метода RunTest: %v", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(config.Content, &data); err != nil {
		return err
	}

	pkg.SaveToTOML("tmp.toml", data)

	file, err := os.Create("tmp")
	if err != nil {
		return fmt.Errorf("не получилось создать файл tmp локально: %v", err)
	}

	_, err = file.Write(test.Binary)
	file.Close()
	if err != nil {
		return fmt.Errorf("ошибка метода RunTest: %v", err)
	}
	//Чтобы успеть закрыть файл
	time.Sleep(100 * time.Millisecond)

	cmd := exec.Command("chmod", "+x", "tmp")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("не получилось выдать право бинарнику на исполнение: %v", err)
	}

	var output string
	if serverIp != "localhost" {
		output, err = pkg.СonnectSSH(serverIp, username, commandTemplate)
	} else {
		localBinary := "./tmp"
		localConfig := "tmp.toml"
		command := strings.ReplaceAll(commandTemplate, "{BIN_FILE}", localBinary)
		command = strings.ReplaceAll(command, "{CONFIG}", localConfig)

		cmd := exec.Command("bash", "-c", command)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			output = stderr.String()
		} else {
			output = stdout.String()
		}
	}

	log := core.Log{Output: output, ConfigID: configId}
	s.repo.AddLog(ctx, &log)
	return err
}
