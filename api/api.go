package api

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"


func RegisterApiRoutes(r *gin.Engine) {
	apisubdomain := r.Group("/api")
	apisubdomain.GET("/rinfo", regionInfo)
	apisubdomain.GET("/ginfo", gameInfo)
	apisubdomain.GET("/pinfo", priceInfo)
	apisubdomain.GET("/gp", gamePrice)
	apisubdomain.GET("/sp", searchPrice)
}

func regionInfo(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(200, formResult(301, string("invalid param id"), gin.H{}))
		return
	}

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

func gameInfo(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(200, formResult(301, string("invalid param id"), gin.H{}))
		return
	}

	r, err := db.QueryGameInfo(id)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			d = gin.H {
				"id": r.Id,
				"name": r.Name,
				"cname": r.Cname,
				"desc": r.Desc,
				"lan": r.Language,
				"cover": r.Cover,
			}
			c.JSON(200, formResult(0, "", d))
		} else {
			c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
		}
	}
}

func priceInfo(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(200, formResult(301, string("invalid param id"), gin.H{}))
		return
	}

	r, err := db.QueryPriceInfo(id)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			d = gin.H {
				"id": r.Id,
				"discount": r.Discount,
				"price": r.Price,
				"lprice": r.LPrice,
				"lregion": r.LRegion,
				"hprice": r.HPrice,
				"hregion": r.HRegion,
			}
			c.JSON(200, formResult(0, "", d))
		} else {
			c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
		}
	}
}

func gamePrice(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(200, formResult(301, string("invalid param id"), gin.H{}))
		return
	}

	r, err := db.QueryGamePrice(id)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			d = formGamePrice(c, *r)
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func searchPrice(c *gin.Context) {
	name := c.Query("name")
	if len(name) == 0 {
		c.JSON(200, formResult(301, string("invalid param name"), gin.H{}))
		return
	}
	r, err := db.QuerySearchGamePrice(name)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			var games []gin.H
			for _, e := range *r {
				games = append(games, formGamePrice(c, e))
			}
			if games != nil {
				d = gin.H {
					"games" : games,
				}
			}
		}
		c.JSON(200, formResult(0, "", d))
	}
}