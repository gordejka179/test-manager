package core

import "time"

type Test struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Binary    []byte    `json:"description"`
	CreatedAt time.Time `json:"created_at"`
}

type Config struct {
	ID        string    `json:"id"`
	TestID    string    `json:"test_id"`
	Name      string    `json:"name"`
	Config    string    `json:"config"`
	CreatedAt time.Time `json:"created_at"`
}

type Log struct {
	ID        string    `json:"id"`
	TestID    string    `json:"test_id"`
	ConfigID  string    `json:"config_id"`
	Output    string    `json:"output"`
	CreatedAt time.Time `json:"created_at"`
}
