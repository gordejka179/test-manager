package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
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
	GetConfigByID(ctx context.Context, configID string) (*core.Config, error)
	GetAllConfigs(ctx context.Context) ([]core.Config, error)
	GetAllConfigsToTest(ctx context.Context, testName string) ([]core.Config, error)
	DeleteConfig(ctx context.Context, id string) error
	GetLogsToConfig(ctx context.Context, configID string) ([]core.Log, error)
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
	data := convertToMap(c.Request.PostForm)

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

	h.service.AddConfig(c, &config)

	c.JSON(http.StatusOK, config)

	//saveToTOML("tmp.toml", data)
}

func convertToMap(v url.Values) map[string]interface{} {
	// Преобразуем FormData в map[string]interface{}
	data := make(map[string]interface{})
	for key, values := range v {
		if len(values) == 0 {
			continue
		}
		value := values[0]
		setNestedValue(&data, key, value)
	}
	return data
}
func setNestedValue(m *map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, ".")
	for i, k := range keys {
		if i == len(keys)-1 {
			(*m)[k] = parseValue(value.(string))
		} else {
			// Если вложенного map нет — создаём
			if _, exists := (*m)[k]; !exists {
				(*m)[k] = make(map[string]interface{})
			}
			nested := (*m)[k].(map[string]interface{})
			setNestedValue(&nested, strings.Join(keys[i+1:], "."), value)
			return
		}
	}
}

func parseValue(s string) interface{} {
	var jsonValue interface{}
	if err := json.Unmarshal([]byte(s), &jsonValue); err == nil {
		return jsonValue
	}

	if s == "[]" {
		return []interface{}{} // Пустой массив
	}
	if s == "true" {
		return true
	}
	if s == "false" {
		return false
	}
	if i, err := strconv.Atoi(s); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(s, 64); err == nil {
		return f
	}
	return s
}

func saveToTOML(filename string, data interface{}) error {
	file, err := os.Create(filename) // Создает файл с указанным именем
	if err != nil {
		return err // Если произошла ошибка при создании файла, возвращается ошибка
	}
	defer file.Close() // Убедитесь, что файл закроется после завершения функции

	encoder := toml.NewEncoder(file) // Создает новый TOML-энкодер для записи в файл
	return encoder.Encode(data)      // Кодирует данные в TOML и записывает их в файл
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
	configId := c.PostForm("config_id")
	logs, err := h.service.GetLogsToConfig(c.Request.Context(), configId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Fatal(err)
		return
	}
	c.JSON(http.StatusOK, logs)
}
