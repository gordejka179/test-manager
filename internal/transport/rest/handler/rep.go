package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gordejka179/test-manager/internal/core"
	"github.com/gordejka179/test-manager/pkg"
)

type TestService interface {
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
	configType := c.PostForm("config_type")
	structureName := c.PostForm("structure_name")

	testFileHeader, err := c.FormFile("test_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен"})
		return
	}

	configFileHeader, err := c.FormFile("config_file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен"})
		return
	}

	testFile, _ := testFileHeader.Open()
	defer testFile.Close()

	testFileBytes, err := io.ReadAll(testFile)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка чтения данных"})
		return
	}

	configFile, _ := configFileHeader.Open()
	defer testFile.Close()

	configFileBytes, err := io.ReadAll(configFile)

	data, err := pkg.ParseStructsFromFile(configFileBytes, structureName)

	//fmt.Println(string(configFileBytes))
	// Выводим результат
	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Ошибка чтения данных"},
		)
		return
	}

	Test := core.Test{Name: name, ConfigType: configType, Binary: testFileBytes, Template: jsonData}
	h.service.AddTest(c, &Test)
	c.JSON(http.StatusOK, Test)

	//saveToTOML("tmp2.toml", data)

}

func (h *TestServiceHandler) AddConfig(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	//fmt.Println(c.Request.PostForm)
	data := pkg.ConvertToMap(c.Request.PostForm)

	testName, ok := data["test_name"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Имя теста не должно быть числом",
		})
		return
	}

	configName, ok := data["config_name"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Имя конфига не должно быть числом",
		})
		return
	}

	configType, ok := data["config_type"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Тип конфига не строка",
		})
		return

	}

	delete(data, "test_name")
	delete(data, "config_name")
	delete(data, "config_type")

	jsonData, err := json.Marshal(data)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "Ошибка чтения данных"},
		)
		return
	}

	config := core.Config{TestName: testName, ConfigType: configType, Name: configName, Content: jsonData}

	id, err := h.service.AddConfig(c, &config)
	if err != nil {
		log.Fatal("Ошибка метода AddConfig:", err)
	}
	config.ID = int(id)

	c.JSON(http.StatusOK, config)

	//saveToTOML("tmp.toml", data)
}

func (h *TestServiceHandler) GetAllConfigsToTest(c *gin.Context) {
	TestName := c.PostForm("testName")
	configs, err := h.service.GetAllConfigsToTest(c, TestName)
	if err != nil {
		log.Fatal("Ошибка метода GetAllConfigsToTest: ", err)
	}

	c.JSON(http.StatusOK, configs)
}

func (h *TestServiceHandler) GetLogsToConfig(c *gin.Context) {
	configIdStr := c.PostForm("config_id")
	configId, err := strconv.Atoi(configIdStr)
	if err != nil {
		log.Fatalf("Ошибка метода GetLogsToConfig: ", err)
	}
	logs, err := h.service.GetLogsToConfig(c.Request.Context(), configId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal(err)
		return
	}
	c.JSON(http.StatusOK, logs)
}
