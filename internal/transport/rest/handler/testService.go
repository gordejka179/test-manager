package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
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

	data, err := ParseStructsFromFile(configFileBytes, structureName)

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

func ParseStructsFromFile(src []byte, rootStructName string) (map[string]interface{}, error) {

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %v", err)
	}

	// Собираем информацию о структурах и типах
	structs := make(map[string]*ast.StructType)
	typeDefs := make(map[string]ast.Expr)
	constValues := make(map[string]interface{})

	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			if d.Tok == token.TYPE {
				for _, spec := range d.Specs {
					typeSpec := spec.(*ast.TypeSpec)
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						structs[typeSpec.Name.Name] = structType
					} else {
						typeDefs[typeSpec.Name.Name] = typeSpec.Type
					}
				}
			} else if d.Tok == token.CONST {
				for _, spec := range d.Specs {
					valueSpec := spec.(*ast.ValueSpec)
					for i, name := range valueSpec.Names {
						if len(valueSpec.Values) > i {
							if basicLit, ok := valueSpec.Values[i].(*ast.BasicLit); ok {
								switch basicLit.Kind {
								case token.INT:
									val, _ := strconv.Atoi(basicLit.Value)
									constValues[name.Name] = val
								}
							}
						}
					}
				}
			}
		}
	}

	// Проверяем, что корневая структура существует
	if _, exists := structs[rootStructName]; !exists {
		return nil, fmt.Errorf("root struct '%s' not found in file", rootStructName)
	}

	// Рекурсивная функция для построения map
	var buildMap func(structName string) map[string]interface{}
	buildMap = func(structName string) map[string]interface{} {
		// Проверяем, является ли тип структурой
		if structType, exists := structs[structName]; exists {
			result := make(map[string]interface{})
			for _, field := range structType.Fields.List {
				if len(field.Names) == 0 {
					continue
				}

				fieldName := field.Names[0].Name
				tag := ""
				if field.Tag != nil {
					tagValue := strings.Trim(field.Tag.Value, "`")
					tagParts := strings.Split(tagValue, ":")
					if len(tagParts) > 1 {
						tag = strings.Trim(strings.Split(tagParts[1], "\"")[1], "\"")
					}
				}

				if tag == "" {
					tag = strings.ToLower(fieldName)
				}

				// Обрабатываем тип поля
				switch t := field.Type.(type) {
				case *ast.Ident:
					if _, isStruct := structs[t.Name]; isStruct {
						result[tag] = buildMap(t.Name)
					} else if _, isType := typeDefs[t.Name]; isType {
						// Пользовательский тип (например, RunMode)
						result[tag] = 0 // нулевое значение для числового типа
					} else {
						switch t.Name {
						case "string":
							result[tag] = ""
						case "int", "int8", "int16", "int32", "int64",
							"uint", "uint8", "uint16", "uint32", "uint64":
							result[tag] = 0
						case "float32", "float64":
							result[tag] = 0.0
						case "bool":
							result[tag] = false
						default:
							result[tag] = nil
						}
					}
				case *ast.ArrayType:
					result[tag] = []interface{}{}
				case *ast.SelectorExpr:
					result[tag] = nil
				case *ast.StarExpr:
					result[tag] = nil
				case *ast.MapType:
					result[tag] = make(map[string]interface{})
				default:
					result[tag] = nil
				}
			}
			return result
		}
		return nil
	}

	return buildMap(rootStructName), nil
}

func (h *TestServiceHandler) AddConfig(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.String(http.StatusBadRequest, "Bad request")
		return
	}

	//fmt.Println(c.Request.PostForm)
	data := convertToMap(c.Request.PostForm)

	testName, ok := data["test_name"].(string)
	fmt.Printf("%T", data["test_name"])
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

	//fmt.Printf("%#v\n", data)

	//saveToTOML("tmp.toml", data)
}

func convertToMap(v url.Values) map[string]interface{} {
	// Преобразуем FormData в map[string]interface{}
	data := make(map[string]interface{})
	for key, values := range v {
		if len(values) == 0 {
			continue
		}
		value := values[0] // Берём первое значение
		setNestedValue(&data, key, value)
	}
	return data
}
func setNestedValue(m *map[string]interface{}, key string, value interface{}) {
	keys := strings.Split(key, ".")
	for i, k := range keys {
		if i == len(keys)-1 {
			// Последний ключ — записываем значение
			(*m)[k] = parseValue(value.(string))
		} else {
			// Если вложенного map нет — создаём
			if _, exists := (*m)[k]; !exists {
				(*m)[k] = make(map[string]interface{})
			}
			// Получаем указатель на вложенную map
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

func GenerateForm(config map[string]interface{}, prefix string, w *strings.Builder) {
	for key, value := range config {
		fullKey := key
		if prefix != "" {
			fullKey = prefix + "." + key
		}

		switch v := value.(type) {
		case map[string]interface{}: // Вложенная структура
			w.WriteString(fmt.Sprintf("<fieldset><legend>%s</legend>\n", key))
			GenerateForm(v, fullKey, w)
			w.WriteString("</fieldset>\n")

		default: // Обычное поле
			w.WriteString(fmt.Sprintf(
				"<label for='%s'>%s</label>\n"+
					"<input type='text' id='%s' name='%s' value='%v'><br>\n",
				fullKey, key, fullKey, fullKey, v,
			))
		}
	}
}

func CreateHTMLForm(config map[string]interface{}) string {
	var html strings.Builder
	html.WriteString("<form method='POST'>\n")

	GenerateForm(config, "", &html)

	html.WriteString("<input type='submit' value='Save'>\n</form>")
	return html.String()
}
