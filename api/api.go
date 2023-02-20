package api

import (
	"celadon-service/db"
	"celadon-service/ipc"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func RegisterApiRoutes(r *gin.Engine) {
	apisubdomain := r.Group("/celadon")
	apisubdomain.GET("/regioninfo", regionInfo)  // RegionInfo
	apisubdomain.GET("/gameinfo", gameInfo)      // GameInfo
	apisubdomain.GET("/priceinfo", priceInfo)    // PriceInfo
	apisubdomain.GET("/hotwords", queryHotWords) // HotWords

	//apisubdomain.GET("/gameprice", gamePrice)		// GamePrice
	apisubdomain.GET("/search", searchPrice)        // []ContentItem - ok
	apisubdomain.GET("/recommend", queryRecommend)  // []ContentItem - ok
	apisubdomain.GET("/queryplist", queryPriceList) // []ContentItem - ok

	apisubdomain.GET("/cateindex", queryCateIndex) // CateIndex
	apisubdomain.GET("/discover", queryDiscover)   // []ContentItem - ok
}

func PrepareToRun() {
	db.ReCheckCheapGames()
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
				"id":    r.Region_id,
				"name":  r.Name,
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
			d = formGameInfo(c, r)
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

	p, err := db.QueryPriceInfo(id)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if p != nil {
			d = formPrice(c, *p)
		}
		c.JSON(200, formResult(0, "", d))
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
			var items []gin.H
			var ids []string
			for _, e := range *r {
				items = append(items, formContentItem(c, NewContentItemGamePrice(e)))
				ids = append(ids, e.Id)
			}
			if items != nil {
				d = gin.H{
					"items": items,
				}
			}
			invokeIpcTask(ids)
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func queryRecommend(c *gin.Context) {
	sz := c.DefaultQuery("sz", "20")
	no := c.DefaultQuery("no", "0")
	pageSize, err := strconv.Atoi(sz)
	if err != nil {
		pageSize = 20
	}
	pageNo, err := strconv.Atoi(no)
	if err != nil {
		pageNo = 0
	}
	startPos := pageSize * pageNo

	r, err := db.QueryCheapGames(startPos, pageSize)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			var items []gin.H
			var ids []string
			for _, e := range *r {
				items = append(items, formContentItem(c, NewContentItemGamePrice(e)))
				ids = append(ids, e.Id)
			}
			if items != nil {
				d = gin.H{
					"items": items,
				}
			}
			invokeIpcTask(ids)
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func queryPriceList(c *gin.Context) {
	ids := c.Query("ids")
	if len(ids) == 0 {
		c.JSON(200, formResult(301, string("invalid param ids"), gin.H{}))
		return
	}
	s := strings.Split(ids, ",")
	r, err := db.QueryPriceListByIds(s)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			var items []gin.H
			var ids []string
			for _, e := range *r {
				items = append(items, formContentItem(c, NewContentItemGamePrice(e)))
				ids = append(ids, e.Id)
			}
			if items != nil {
				d = gin.H{
					"items": items,
				}
			}
			invokeIpcTask(ids)
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func queryHotWords(c *gin.Context) {
	hwd := GetCurrentHotWords()
	if hwd == nil {
		c.JSON(202, formResult(300, "something wrong ...", gin.H{}))
	} else {
		d := gin.H{
			"hotwords": hwd.Words,
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func queryCateIndex(c *gin.Context) {
	name := c.Query("name")

	if strings.EqualFold(name, "discover") {
		queryDiscoverCateIndex(c)
	} else {
		c.JSON(200, formResult(203, string("no match"), gin.H{}))
	}
}

func invokeIpcTask(id []string) {
	go ipc.AddTask(id)
}
