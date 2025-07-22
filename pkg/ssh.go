package pkg

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
)

func СonnectSSH(serverIp string, username string, commandTemplate string) (string, error) {
	cmd := exec.Command("bash", "-c", "echo $HOME")

	outputHome, err := cmd.CombinedOutput()
	homeDir := strings.TrimSpace(string(outputHome))
	if err != nil {
		return "", fmt.Errorf("ошибка выполения команды: %w", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			PublicKeyFile(homeDir + "/.ssh/key2"), //путь к приватному ключу, ВАЖНО: полный путь
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", serverIp+":22", sshConfig)
	if err != nil {
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer client.Close()

	// Копируем бинарник на сервер
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Не удалось создать сессию: %v", err)
	}
	defer session.Close()

	//Узнаём адрес директории на сервере
	var stdout bytes.Buffer
	session.Stdout = &stdout
	err = session.Run("echo $HOME")

	if err != nil {
		log.Fatalf("Не удалось узнать адрес директории: %v", err)
	}

	localBinary := "tmp"
	remoteBinary := strings.TrimSpace(stdout.String()) + "/tmpBinaryTM"

	if err := copyFile(client, localBinary, remoteBinary); err != nil {
		log.Fatalf("Ошибка копирования бинарника: %v", err)
	}

	localConfig := "tmp.toml"
	remoteConfig := strings.TrimSpace(stdout.String()) + "/tmpConfigTM"

	if err := copyFile(client, localConfig, remoteConfig); err != nil {
		log.Fatalf("Ошибка копирования конфига: %v", err)
	}

	// Выполняем команду на сервере
	command := strings.ReplaceAll(commandTemplate, "{BIN_FILE}", remoteBinary)
	command = strings.ReplaceAll(command, "{CONFIG}", remoteConfig)

	output, err := runCommand(client, command)
	if err != nil {
		return "", fmt.Errorf("команда выполнилась с ошибкой: %w", err)
	}

	return output, nil
}

// возвращает AuthMethod для аутентификации по ключу
func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}

	return ssh.PublicKeys(key)
}

// копирует файл на удаленный сервер через SCP
func copyFile(client *ssh.Client, localPath, remotePath string) error {
	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	file, err := os.Open(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// Канал для передачи ошибок
	errChan := make(chan error, 1)

	// Запускаем SCP приемник на сервере
	go func() {
		w, err := session.StdinPipe()
		if err != nil {
			errChan <- err
			return
		}
		defer w.Close()

		// Отправляем заголовок SCP
		fmt.Fprintf(w, "C0644 %d %s\n", stat.Size(), filepath.Base(remotePath))
		if _, err := io.Copy(w, file); err != nil {
			errChan <- err
			return
		}
		fmt.Fprint(w, "\x00") // Передаем нулевой байт для завершения
		errChan <- nil
	}()

	// Выполняем команду scp на сервере
	if err := session.Run(fmt.Sprintf("scp -t %s", filepath.Dir(remotePath))); err != nil {
		return fmt.Errorf("failed to run scp: %v", err)
	}

	if err := <-errChan; err != nil {
		return fmt.Errorf("ошибка команды scp: %v", err)
	}

	return nil
}

// runCommand выполняет команду на удаленном сервере и возвращает вывод
func runCommand(client *ssh.Client, cmd string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(cmd)
	if err != nil {
		return stderr.String(), err
	}

	return stdout.String(), nil
}
