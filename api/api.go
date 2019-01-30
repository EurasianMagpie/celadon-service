package api

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"


func RegisterApiRoutes(r *gin.Engine) {
	apisubdomain := r.Group("/api")
	apisubdomain.GET("/regionInfo", regionInfo)
	apisubdomain.GET("/gamePrice", gamePrice)
}

func regionInfo(c *gin.Context) {
	id := c.Query("id")
	r, err := db.QueryRegionInfo(id)
	if err == nil && r != nil {
		c.JSON(200, gin.H{
			"api": "regionInfo",
			"id": r.Region_id,
			"name": r.Name,
			"cname": r.Cname,
		})
	} else {
		c.JSON(200, gin.H{
			"api": "regionInfo",
			"error": "no result",
		})
	}//*/
}

func gamePrice(c *gin.Context) {
	id := c.Query("id")
	r, err := db.QueryGamePrice(id)
	if err != nil {
		c.JSON(200, gin.H{
			"api": "gamePrice",
			"error": "no result",
		})
	} else {
		c.JSON(200, gin.H{
			"api": "gamePrice",
			"id": r.Id,
			"name": r.Name,
			"region": r.Region,
			"price": r.Price,
		})
	}
}