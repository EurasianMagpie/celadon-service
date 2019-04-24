package api

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"
import "github.com/EurasianMagpie/celadon/operation"


func formGameInfo(c *gin.Context, g *db.GameInfo) gin.H {
	p, err := getGameCoverRefPath(g.Id)
	if err == nil {
		p = "http://" + c.Request.Host + p
	}

	return gin.H {
		"id": g.Id,
		"name": g.Name,
		"cname": g.Cname,
		"publisher": g.Publisher,
		"desc": g.Desc,
		"cover": p,
		"lan": g.Language,
		"tags": g.Tags,
		"realcard": g.RealCard,
	}
}

func formGamePrice(c *gin.Context, g db.GamePrice) gin.H {
	p, err := getGameCoverRefPath(g.Id)
	if err == nil {
		p = "http://" + c.Request.Host + p
	}
	regionName := g.Region
	r, err := db.FindRegionByAbbr(g.Region)
	if err == nil {
		regionName = r.Cname
	}

	if !operation.HasCname(g.Id) {
		g.IsLowest = 0
	}
	return gin.H {
		"id": g.Id,
		"name": g.Name,
		"cname": g.Cname,
		"cover": p,
		"region": regionName,
		"price": g.Price,
		"lowest": g.IsLowest,
	}
}

func formPrice(c *gin.Context, p db.Price) gin.H {
	rankData, _ := calcRegionPriceRank(p.Price)
	if !operation.HasCname(p.Id) {
		p.LowestPrice = ""
		p.IsLowest = 0
	}
	return gin.H {
		"id": p.Id,
		"rank": formPriceRank(rankData),
		"discount": p.Discount,
		"lprice": p.LPrice,
		"lregion": p.LRegion,
		"hprice": p.HPrice,
		"hregion": p.HRegion,
		"lowestprice": p.LowestPrice,
		"islowest": p.IsLowest,
	}
}

func formResult(errno int, errmsg string, data gin.H) gin.H {
	return gin.H {
		"errno": errno,
		"errmsg": errmsg,
		"data" : data,
	}
}