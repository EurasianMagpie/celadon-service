package api

import (
	"encoding/json"
	"os"
	"fmt"
)

type HotWords struct {
	Words []string `json:"hotwords"`
}

var initHotWords = false
var hotWords HotWords

func GetCurrentHotWords() *HotWords {
	if !initHotWords {
		initHotWords = true
		file, err := os.Open("data/hotwords.json")
		if err != nil {
			fmt.Println("config.load | error:", err)
			return nil
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		err = decoder.Decode(&hotWords)
		if err != nil {
			fmt.Println("config.load | error:", err)
			return nil
		}
		fmt.Println("HotWords : ", hotWords)
	}
	return &hotWords
}