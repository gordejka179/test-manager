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
	AddConfig(ctx context.Context, config *core.Config) error
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

	c.JSON(http.StatusOK, gin.H{
		"name": name,
	})
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
	h.service.AddConfig(c, &config)

	c.JSON(http.StatusOK, gin.H{
		"name":   ConfigName,
		"config": fileText,
	})
}
