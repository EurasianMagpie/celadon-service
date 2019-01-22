package api

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"


func RegisterApiRoutes(r *gin.Engine) {
	apisubdomain := r.Group("/api")
	apisubdomain.GET("/regionInfo", regionInfo)
	apisubdomain.GET("/gamePrice", regionInfo)
}

func regionInfo(c *gin.Context) {
	id := c.Query("id")
	r := db.QueryRegionInfo(id)
	c.JSON(200, gin.H{
		"api": "regionInfo",
		"id": r.Region_id,
		"name": r.Name,
		"cname": r.Cname,
	})
}

func gamePrice(c *gin.Context) {
	id := c.Query("id")
	r := db.QueryGamePrice(id)
	c.JSON(200, gin.H{
		"api": "gamePrice",
		"id": r.Id,
		"Name": r.Name,
		"Region": r.Region,
		"Price": r.Price,
	})
}