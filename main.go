package main

import "fmt"
import "flag"
import "strings"

import "github.com/gin-gonic/gin"

//import "github.com/EurasianMagpie/celadon/debug"
import "github.com/EurasianMagpie/celadon/api"
import "github.com/EurasianMagpie/celadon/mon"

func main()  {
	//debug.Info()
	t := flag.String("type", "", "process type. api or mon")
	flag.Parse()

	if strings.Compare(*t, "api") == 0 {
		fmt.Println("main | api")
		r := gin.Default()
		api.RegisterApiRoutes(r)
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong - gin",
			})
		})
		r.Run()
	} else if strings.Compare(*t, "mon") == 0 {
		fmt.Println("main | mon")
		mon.FetchPage()
	} else {
		fmt.Println("Please specify process type")
	}
}