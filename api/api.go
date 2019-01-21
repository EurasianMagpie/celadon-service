package api

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"


func RegisterApiRoutes(r *gin.Engine) {
	apisubdomain := r.Group("/api")
	apisubdomain.GET("/regionInfo", regionInfo)
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