package api

import (
	"fmt"

	"github.com/EurasianMagpie/celadon-service/db"
	"github.com/EurasianMagpie/celadon-service/operation"
	"github.com/gin-gonic/gin"
)

func formCoverUrl(c *gin.Context, id string) (string, error) {
	p, err := getGameCoverRefPath(id)
	if err == nil {
		p = "http://" + c.Request.Host + "/" + p
		//p = "http://192.3.80.174/" + p
	}
	fmt.Println("formCoverUrl", p)
	return p, nil
}

func formGameInfo(c *gin.Context, g *db.GameInfo) gin.H {
	p, _ := formCoverUrl(c, g.Id)

	return gin.H{
		"id":        g.Id,
		"name":      g.Name,
		"cname":     g.Cname,
		"publisher": g.Publisher,
		"desc":      g.Desc,
		"cover":     p,
		"lan":       g.Language,
		"tags":      g.Tags,
		"realcard":  g.RealCard,
	}
}

func formGamePrice(c *gin.Context, g db.GamePrice) gin.H {
	p, _ := formCoverUrl(c, g.Id)

	regionName := g.Region
	r, err := db.FindRegionByAbbr(g.Region)
	if err == nil {
		regionName = r.Cname
	}

	if !operation.HasCname(g.Id) {
		g.IsLowest = 0
	}
	return gin.H{
		"id":     g.Id,
		"name":   g.Name,
		"cname":  g.Cname,
		"cover":  p,
		"region": regionName,
		"price":  g.Price,
		"lowest": g.IsLowest,
	}
}

func formPrice(c *gin.Context, p db.Price) gin.H {
	rankData, _ := calcRegionPriceRank(p.Price)
	if !operation.HasCname(p.Id) {
		p.LowestPrice = ""
		p.IsLowest = 0
	}
	return gin.H{
		"id":          p.Id,
		"rank":        formPriceRank(rankData),
		"discount":    p.Discount,
		"lprice":      p.LPrice,
		"lregion":     p.LRegion,
		"hprice":      p.HPrice,
		"hregion":     p.HRegion,
		"lowestprice": p.LowestPrice,
		"islowest":    p.IsLowest,
	}
}

func formContentItem(c *gin.Context, item ContentItem) gin.H {
	if item.Type == ContentType_GamePrice {
		return gin.H{
			"type": item.Type,
			"gp":   formGamePrice(c, item.Data.(db.GamePrice)),
		}
	} else {
		return gin.H{
			"type": item.Type,
		}
	}
}

func formResult(errno int, errmsg string, data gin.H) gin.H {
	return gin.H{
		"errno":  errno,
		"errmsg": errmsg,
		"data":   data,
	}
}
