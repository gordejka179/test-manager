package core

import "encoding/json"

type Test struct {
	Name       string          `json:"name"`
	ConfigType string          `json:"config_type"`
	Binary     []byte          `json:"binary"`
	Template   json.RawMessage `json:"template"`
}

type Config struct {
	ID         int             `json:"id"`
	TestName   string          `json:"test_name"`
	Name       string          `json:"name"`
	ConfigType string          `json:"config_type"`
	Content    json.RawMessage `json:"content"`
}

type Log struct {
	ID       int    `json:"id"`
	TestName string `json:"test_name"`
	ConfigID string `json:"config_id"`
	Output   string `json:"output"`
}
