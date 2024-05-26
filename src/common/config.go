package common

import (
	"encoding/json"
	"os"
)

type Config interface {
}

func createNewConfig[T Config](c Config, path string) {

	jsonData, err := json.Marshal(c)

	if err != nil {
		panic("Error => JSON encoding error.")
	}

	file, err := os.Create(path)

	if err != nil {
		panic("Error => Config file can't be created.")
	}

	file.Write(jsonData)

	defer file.Close()
}

func InitializeConfig[T Config](c Config, path string) {

	if _, e := os.Stat(path); e != nil {
		createNewConfig[T](c, path)
	}

	cbytes, er := os.ReadFile(path)

	if er != nil {
		panic("Error => Can't read config file!")
	}

	err := json.Unmarshal(cbytes, c)

	if err != nil {
		panic("Error => Can't read config file!")
	}
}
