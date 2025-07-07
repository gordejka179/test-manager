package handler

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gordejka179/test-manager/internal/core"
)

type TestService interface {
	AddTest(ctx context.Context, test *core.Test) error
	GetTestByName(ctx context.Context, name string) (*core.Test, error)
	GetAllTests(ctx context.Context) ([]core.Test, error)
	DeleteTest(ctx context.Context, name string) error
	AddConfig(ctx context.Context, config *core.Config) (int64, error)
	GetConfigByID(ctx context.Context, configID string) (*core.Config, error)
	GetAllConfigs(ctx context.Context) ([]core.Config, error)
	GetAllConfigsToTest(ctx context.Context, testName string) ([]core.Config, error)
	DeleteConfig(ctx context.Context, id string) error
	GetLogs(ctx context.Context, testName string, configID string) ([]core.Log, error)
}

type TestServiceHandler struct {
	service TestService
}

func NewTestServiceHandler(S TestService) *TestServiceHandler {
	return &TestServiceHandler{service: S}
}

func (h *TestServiceHandler) GetAllTests(c *gin.Context) {
	tests, err := h.service.GetAllTests(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal(err)
		return
	}

	c.JSON(http.StatusOK, tests)
}

func (h *TestServiceHandler) AddTest(c *gin.Context) {
	name := c.PostForm("name")

	fileHeader, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен"})
		return
	}

	file, _ := fileHeader.Open()
	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
		return
	}
	Test := core.Test{Name: name, Binary: fileBytes}
	h.service.AddTest(c, &Test)

	c.JSON(http.StatusOK, Test)
}

func (h *TestServiceHandler) AddConfig(c *gin.Context) {
	TestName := c.PostForm("testName")
	ConfigName := c.PostForm("configName")

	fileHeader, err := c.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен"})
	}

	file, _ := fileHeader.Open()
	defer file.Close()

	fileBytes, err := io.ReadAll(file)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
	}

	fileText := string(fileBytes)
	config := core.Config{TestName: TestName, Name: ConfigName, Config: fileText}
	id, err := h.service.AddConfig(c, &config)
	if err != nil {
		log.Fatal("Ошибка метода AddConfig: ", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"test_name": TestName,
		"id":        id,
		"name":      ConfigName,
		"config":    fileText,
	})
}

func (h *TestServiceHandler) GetAllConfigsToTest(c *gin.Context) {
	TestName := c.PostForm("testName")
	configs, err := h.service.GetAllConfigsToTest(c, TestName)
	if err != nil {
		log.Fatal("Ошибка метода GetAllConfigsToTest: ", err)
	}

	c.JSON(http.StatusOK, configs)
}
