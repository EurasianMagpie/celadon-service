package operation

import (
	"encoding/json"
	"os"
	"fmt"
)
import "github.com/EurasianMagpie/celadon/db"


type LocalData struct {
	Cname map[string]string `json:"cname"`
}

var localData LocalData

func UpdateCnameFromLocalData() {
	file, err := os.Open("data/cname.json")
	if err != nil {
		fmt.Println("config.load | error:", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&localData)
	if err != nil {
		fmt.Println("config.load | error:", err)
	}

	for k, v := range localData.Cname {
		fmt.Println("Update CName > ", k, v)
		db.UpdateGameCname(k, v)
	}
}