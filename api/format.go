package api

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"


func formGamePrice(c *gin.Context, g db.GamePrice) gin.H {

	p, err := getGameCoverFilePath(g.Id)
	if err == nil {
		p = "http://" + c.Request.Host + p
	}
	regionName := g.Region
	r, err := db.FindRegion(g.Region)
	if err == nil {
		regionName = r.Cname
	}

	return gin.H {
		"id": g.Id,
		"name": g.Name,
		"cname": g.Cname,
		"cover": p,
		"region": regionName,
		"price": g.Price,
	}
}

func formResult(errno int, errmsg string, data gin.H) gin.H {
	return gin.H {
		"errno": errno,
		"errmsg": errmsg,
		"data" : data,
	}
}