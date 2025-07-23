package pkg

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"strings"
)

func ParseStructsFromFile(src []byte, rootStructName string) (map[string]interface{}, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("failed to parse file: %v", err)
	}

	// Собираем информацию о структурах
	structs := make(map[string]*ast.StructType)

	// Сбор структур
	for _, decl := range f.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						structs[typeSpec.Name.Name] = structType
					}
				}
			}
		}
	}

	// Проверка корневой структуры
	if _, exists := structs[rootStructName]; !exists {
		return nil, fmt.Errorf("root struct '%s' not found in file", rootStructName)
	}

	// построение map
	var buildMap func(structType *ast.StructType) map[string]interface{}
	buildMap = func(structType *ast.StructType) map[string]interface{} {
		result := make(map[string]interface{})

		for _, field := range structType.Fields.List {
			if len(field.Names) == 0 {
				continue
			}

			tag := getMapstructureTag(field)

			switch ft := field.Type.(type) {
			case *ast.StructType:
				// Для вложенных анонимных структур
				result[tag] = buildMap(ft)
			case *ast.Ident:
				// Базовые типы
				switch ft.Name {
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
					// Вложенная структура, но не анонимная
					if nestedStruct, exists := structs[ft.Name]; exists {
						result[tag] = buildMap(nestedStruct)
					} else {
						result[tag] = nil
					}
				}
			case *ast.ArrayType:
				result[tag] = []interface{}{}
			case *ast.StarExpr:
				// Указатели
				result[tag] = nil
			case *ast.MapType:
				result[tag] = make(map[string]interface{})
			case *ast.SelectorExpr:
				// Селекторы
				result[tag] = nil
			default:
				result[tag] = nil
			}
		}

		return result
	}

	return buildMap(structs[rootStructName]), nil
}

// Функция для извлечения тега mapstructure или toml
func getMapstructureTag(field *ast.Field) string {
	if field.Tag == nil {
		return field.Names[0].Name
	}

	tagStr := strings.Trim(field.Tag.Value, "`")
	tag := reflect.StructTag(tagStr)
	//Если в бинарнике напротив поля написано toml
	if tagValue := tag.Get("toml"); tagValue != "" {
		parts := strings.Split(tagValue, ",")
		return parts[0]
	}
	//Если в бинарнике напротив поля написано mapstructure
	if tagValue := tag.Get("mapstructure"); tagValue != "" {
		parts := strings.Split(tagValue, ",")
		return parts[0]
	}

	return field.Names[0].Name
}
