package pkg

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/ssh"
)

func СonnectSSH() string {
	sshConfig := &ssh.ClientConfig{
		User: "t-bmstu",
		Auth: []ssh.AuthMethod{
			PublicKeyFile("/home/ivan/.ssh/key2"), //путь к приватному ключу, ВАЖНО: полный путь
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	fmt.Println("maamamamamamamaam")

	client, err := ssh.Dial("tcp", "195.19.40.45:22", sshConfig)
	if err != nil {
		fmt.Println("adadaddadad")
		log.Fatalf("Не удалось подключиться: %v", err)
	}
	defer client.Close()

	// Копируем бинарник на сервер
	session, err := client.NewSession()
	if err != nil {
		log.Fatalf("Не удалось создать сессию: %v", err)
	}
	defer session.Close()

	localBinary := "tmp"
	remoteBinary := "/home/t-bmstu/tmpTest-Manager"

	if err := copyFile(client, localBinary, remoteBinary); err != nil {
		log.Fatalf("Ошибка копирования бинарника: %v", err)
	}

	localConfig := "tmp.toml"
	remoteConfig := "/home/t-bmstu/tmpConfigTest-Manager"

	if err := copyFile(client, localConfig, remoteConfig); err != nil {
		log.Fatalf("Ошибка копирования конфига: %v", err)
	}

	// Выполняем команду на сервере
	cmd := fmt.Sprintf("chmod +x %s && %s && cat %s ", remoteBinary, remoteBinary, remoteConfig)
	output, err := runCommand(client, cmd)
	if err != nil {
		log.Fatalf("Команда выполнилась с ошибкой: %v", err)
	}

	fmt.Printf("Вывод программы:\n%s\n", output)
	return output
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
