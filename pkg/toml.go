package pkg

import (
	"os"

	"github.com/pelletier/go-toml"
)

func fixInts(data map[string]interface{}) map[string]interface{} {
	for k, v := range data {
		switch val := v.(type) {
		case float64:
			// Если число выглядит как целое (например, 5.0), сохраняем как int
			if val == float64(int(val)) {
				data[k] = int(val)
			}
		case []interface{}:
			// Обрабатываем слайсы
			for i, item := range val {
				if f, ok := item.(float64); ok && f == float64(int(f)) {
					val[i] = int(f)
				}
			}
		case map[string]interface{}:
			// Рекурсивно обрабатываем вложенные мапы
			data[k] = fixInts(val)
		}
	}
	return data
}

func SaveToTOML(filename string, data map[string]interface{}) error {
	data = fixInts(data)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(data)
}
