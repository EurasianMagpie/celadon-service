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
	if err == nil {
		d := gin.H{}
		if r != nil {
			d = gin.H{
				"id": r.Region_id,
				"name": r.Name,
				"cname": r.Cname,
			}
		}
		c.JSON(200, formResult(0, "", d))
	} else {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	}
}

func gamePrice(c *gin.Context) {
	id := c.Query("id")
	r, err := db.QueryGamePrice(id)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			d = gin.H {
				"id": r.Id,
				"name": r.Name,
				"region": r.Region,
				"price": r.Price,
			}
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func searchPrice(c *gin.Context) {
	name := c.Query("name")
	r, err := db.QuerySearchGamePrice(name)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			d = gin.H {
				"games" : r,
			}
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func formResult(errno int, errmsg string, data gin.H) gin.H {
	return gin.H {
		"errno": errno,
		"errmsg": errmsg,
		"data" : data,
	}
}