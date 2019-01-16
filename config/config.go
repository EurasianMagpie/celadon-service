package config

import (
	"encoding/json"
	"os"
	"fmt"
)

type Config struct {
	Db DB `json:"db"`
}

type DB struct {
	Host string `json:"host"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
}

var cfg Config
var loaded = false

func init() {
	load()
}

func load() {
	file, err := os.Open("config/conf.json")
	if err != nil {
		fmt.Println("config.load | error:", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println("config.load | error:", err)
	}
	loaded = true
}

func GetConfig() Config {
	if !loaded {
		load()
	}
	return cfg
}