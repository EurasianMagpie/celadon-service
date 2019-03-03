package api

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"


func formGamePrice(g db.GamePrice) gin.H {
	return gin.H {
		"id": g.Id,
		"name": g.Name,
		"cname": g.Cname,
		"cover": g.Cover,
		"region": g.Region,
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