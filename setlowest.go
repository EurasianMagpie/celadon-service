package main

import "fmt"
import "flag"

import "github.com/EurasianMagpie/celadon-service/db"

func main()  {

	id := flag.String("id", "", "game id")
	prc := flag.String("price", "", "lowest price")
	flag.Parse()

	gameId := *id
	lowestPrice := *prc
	if len(gameId) > 0 && len(lowestPrice) > 0 {
		if db.UpdateLowestPrice(gameId, lowestPrice, "") {
			fmt.Println("update succeeded")
		} else {
			fmt.Println("update failed")
		}
	} else {
		fmt.Println("invalid params")
	}
}