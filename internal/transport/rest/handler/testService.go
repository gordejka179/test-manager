package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

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
	configType := c.PostForm("configType")

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
	Test := core.Test{Name: name, ConfigType: configType, Binary: fileBytes}
	h.service.AddTest(c, &Test)

	c.JSON(http.StatusOK, Test)
}

func (h *TestServiceHandler) AddConfig(c *gin.Context) {
	configType := c.PostForm("config_type")
	if configType == "toml" {
		name := c.PostForm("config_name")
		testName := c.PostForm("test_name")
		hosts := c.PostFormArray("hosts")
		user := c.PostForm("user")
		password := c.PostForm("password")
		dbName := c.PostForm("db_name")
		dbConfig := c.PostForm("db_config")
		mode, err := strconv.Atoi(c.PostForm("mode"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "mode должен быть числом"})

		}

		readerGoroutines, err := strconv.Atoi(c.PostForm("reader_goroutines"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "reader_goroutines должен быть числом"})

		}

		consumerGoroutines, err := strconv.Atoi(c.PostForm("consumer_goroutines"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "consumer_goroutines должен быть числом"})

		}

		invalidFieldsLog := c.PostForm("invalid_fields_Log")
		missingDocsImportedLog := c.PostForm("missing_docs_imported_log")

		content := map[string]interface{}{
			"hosts":                     hosts,
			"user":                      user,
			"password":                  password,
			"db_name":                   dbName,
			"db_config":                 dbConfig,
			"mode":                      mode,
			"reader_goroutines":         readerGoroutines,
			"consumer_goroutines":       consumerGoroutines,
			"invalid_fields_log":        invalidFieldsLog,
			"missing_docs_imported_log": missingDocsImportedLog,
		}

		jsonData, err := json.Marshal(content)
		if err != nil {
			log.Fatal("failed to marshal config: ", err)
		}

		config := core.Config{TestName: testName, Name: name, ConfigType: configType, Content: jsonData}

		newConfigId, err := h.service.AddConfig(c, &config)
		if err != nil {
			log.Fatal("Ошибка AddConfig: ", err)
		}

		c.JSON(http.StatusOK, gin.H{"id": newConfigId,
			"test_name":                 testName,
			"name":                      name,
			"config_type":               configType,
			"hosts":                     hosts,
			"user":                      user,
			"password":                  password,
			"db_name":                   dbName,
			"db_config":                 dbConfig,
			"mode":                      mode,
			"reader_goroutines":         readerGoroutines,
			"consumer_goroutines":       consumerGoroutines,
			"invalid_fields_log":        invalidFieldsLog,
			"missing_docs_imported_log": missingDocsImportedLog,
		})
	}

}

func (h *TestServiceHandler) GetAllConfigsToTest(c *gin.Context) {
	TestName := c.PostForm("testName")
	configs, err := h.service.GetAllConfigsToTest(c, TestName)
	if err != nil {
		log.Fatal("Ошибка метода GetAllConfigsToTest: ", err)
	}

	fmt.Println(configs)
	c.JSON(http.StatusOK, configs)
}
