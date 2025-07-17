package pkg

import (
	"os"

	"github.com/BurntSushi/toml"
)

func SaveToTOML(filename string, data interface{}) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(data)
}
