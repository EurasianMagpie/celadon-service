package operation

import (
	"encoding/json"
	"os"
	"fmt"
)
import "github.com/EurasianMagpie/celadon/db"


type OpData struct {
	Cname map[string]string `json:"cname"`
	RealCard map[string]int `json:"realcard"`
}

var dataFileName = "data/opdata.json"
var opData OpData

func LoadOpData() {
	file, err := os.Open(dataFileName)
	if err != nil {
		fmt.Println("config.load | error:", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&opData)
	if err != nil {
		fmt.Println("config.load | error:", err)
	}

	for k, v := range opData.Cname {
		fmt.Println("Update CName > ", k, v)
		db.UpdateGameCname(k, v)
	}

	for k, v := range opData.RealCard {
		fmt.Println("Update RealCard > ", k, v)
		db.UpdateGameRealCard(k, v)
	}
}