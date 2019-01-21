package main

import "fmt"

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/debug"
import "github.com/EurasianMagpie/celadon/api"

func main()  {
	debug.Info()
	fmt.Println("api.main")

	r := gin.Default()
	api.RegisterApiRoutes(r)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong - gin",
		})
	})
	r.Run()
}