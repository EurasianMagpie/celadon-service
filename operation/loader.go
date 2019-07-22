package operation

import (
	"encoding/json"
	"os"
	"fmt"
)
import "github.com/EurasianMagpie/celadon-service/db"

type GameData struct {
	Cname string `json:"cname"`
	RealCard int `json:"card"`
	Publisher string `json:"publisher"`
	ReleaseTime string `json:"time"`
	Lowest string `json:"lowest"`
}

type OpData struct {
	//Cname map[string]string `json:"cname"`
	//RealCard map[string]int `json:"realcard"`
	Data map[string]GameData `json:"data"`
}

var dataFileName = "data/opdata.json"
var opData OpData

func init() {
	load()
}

func load() {
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
}

func LoadUpdateOpData() {

	for k, data := range opData.Data {
		fmt.Println("Update OpData > ", k, data)
		if len(data.Cname) > 0 {
			db.UpdateGameCname(k, data.Cname)
		}
		if data.RealCard > 0 {
			db.UpdateGameRealCard(k, data.RealCard)
		}
	}
}

func HasCname(id string) bool {
	if data, ok := opData.Data[id]; ok {
		return len(data.Cname) > 0
	}
	return false
}