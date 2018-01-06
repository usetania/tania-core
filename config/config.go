package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	File struct {
		AreaPhoto string
	}
}

func Init() Configuration {
	file, _ := os.Open("conf.json")
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		panic("Please setup your config file")
	}

	return configuration
}
