package pkg

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
)

// Преобразуем FormData в map[string]interface{}
func ConvertToMap(v url.Values) map[string]interface{} {
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
