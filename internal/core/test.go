package core

type Test struct {
	Name   string `json:"name"`
	Binary []byte `json:"binary"`
}

type Config struct {
	ID       int    `json:"id"`
	TestName string `json:"test_name"`
	Name     string `json:"name"`
	Config   string `json:"config"`
}

type Log struct {
	ID       int    `json:"id"`
	TestName string `json:"test_name"`
	ConfigID string `json:"config_id"`
	Output   string `json:"output"`
}
