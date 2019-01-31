package api

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"


func RegisterApiRoutes(r *gin.Engine) {
	apisubdomain := r.Group("/api")
	apisubdomain.GET("/rinfo", regionInfo)
	apisubdomain.GET("/gp", gamePrice)
	apisubdomain.GET("/sp", searchPrice)
}

func regionInfo(c *gin.Context) {
	id := c.Query("id")
	r, err := db.QueryRegionInfo(id)
	if err == nil && r != nil {
		c.JSON(200, gin.H{
			"api": "rinfo",
			"id": r.Region_id,
			"name": r.Name,
			"cname": r.Cname,
		})
	} else {
		c.JSON(200, gin.H{
			"api": "rinfo",
			"error": "no result",
		})
	}//*/
}

func gamePrice(c *gin.Context) {
	id := c.Query("id")
	r, err := db.QueryGamePrice(id)
	if err != nil {
		c.JSON(200, gin.H{
			"api": "gp",
			"error": "no result",
		})
	} else {
		c.JSON(200, gin.H{
			"api": "gp",
			"id": r.Id,
			"name": r.Name,
			"region": r.Region,
			"price": r.Price,
		})
	}
}

func searchPrice(c *gin.Context) {
	name := c.Query("name")
	r, err := db.QuerySearchGamePrice(name)
	if err != nil {
		c.JSON(200, gin.H{
			"api": "sp",
			"error": string(err.Error()),
		})
	} else {
		var s string
		for _, p := range *r {
			s = s + p.Name
		}
		c.JSON(200, gin.H{
			"api": "sp",
			"r":  r,
		})
	}
}