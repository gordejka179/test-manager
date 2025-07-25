package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"
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

	if _, err := os.Stat("test-manager"); os.IsNotExist(err) {
		err = os.Mkdir("test-manager", 0755)
		if err != nil {
			return fmt.Errorf("ошибка создания папки test-manager на локальной машине: %v", err)
		}
	}

	pkg.SaveToTOML("./test-manager/tmpConfig.toml", data)

	file, err := os.Create("./test-manager/tmpBinary")
	if err != nil {
		return fmt.Errorf("не получилось создать файл tmpBinary локально: %v", err)
	}

	_, err = file.Write(test.Binary)
	file.Close()
	if err != nil {
		return fmt.Errorf("ошибка метода RunTest: %v", err)
	}
	//Чтобы успеть закрыть файл
	time.Sleep(100 * time.Millisecond)

	cmd := exec.Command("chmod", "+x", "./test-manager/tmpBinary")
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("не получилось выдать право бинарнику на исполнение: %v", err)
	}

	var output string
	if serverIp != "localhost" {
		output, err = pkg.СonnectSSH(serverIp, username, commandTemplate)
	} else {
		localBinary := "./test-manager/tmpBinary"
		localConfig := "./test-manager/tmpConfig.toml"
		command := strings.ReplaceAll(commandTemplate, "{BIN_FILE}", localBinary)
		command = strings.ReplaceAll(command, "{CONFIG}", localConfig)

		cmd := exec.Command("bash", "-c", command)
		var stdout, stderr bytes.Buffer
		var otherError string
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			switch err := err.(type) {
			case *exec.ExitError:
				// Ошибка выполнения команды (ненулевой код или сигнал)

				exitCode := err.ExitCode()
				// Если завершено сигналом
				if status, ok := err.Sys().(syscall.WaitStatus); ok && status.Signaled() {
					signal := status.Signal()
					otherError = fmt.Sprintf("Команда убита сигналом: %v (код: %d)\n", signal, exitCode)
				} else {
					//Обычный ненулевой код возврата
					otherError = fmt.Sprintf("Команда завершилась с кодом ошибки %d\n.", exitCode)
				}

			case *exec.Error:
				// Ошибка запуска (например, команда не найдена)

				otherError = fmt.Sprintf("Ошибка запуска: %v\n", err)

			default:
				// Все остальные ошибки

				otherError = fmt.Sprintf("Неизвестная ошибка: %T %v\n", err, err)
			}
		}

		output = fmt.Sprintf("\n Поток ошибок: %s\n Поток вывода: %s\n Другая ошибка: %s\n", stderr.String(), stdout.String(), otherError)
	}

	log := core.Log{Output: output, ConfigID: configId, TestName: testName}
	err = s.repo.AddLog(ctx, &log)
	if err != nil {
		return fmt.Errorf("ошибка метода RunTest при добавлении лога: %v", err)
	}
	return err
}
