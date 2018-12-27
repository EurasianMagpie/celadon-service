package main

import "fmt"

import "github.com/EurasianMagpie/celadon/db"
import "github.com/EurasianMagpie/celadon/debug"

func main()  {
	debug.Info()
	db.Getdb()
	fmt.Println("api.main")
}