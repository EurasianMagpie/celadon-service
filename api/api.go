package api

import (
	"strings"
	"strconv"
)

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"
import "github.com/EurasianMagpie/celadon/ipc"


func RegisterApiRoutes(r *gin.Engine) {
	apisubdomain := r.Group("/celadon")
	apisubdomain.GET("/rinfo", regionInfo)
	apisubdomain.GET("/ginfo", gameInfo)
	apisubdomain.GET("/pinfo", priceInfo)
	apisubdomain.GET("/gp", gamePrice)
	apisubdomain.GET("/sp", searchPrice)
	apisubdomain.GET("/recommend", queryRecommend)
	apisubdomain.GET("/plist", queryPriceList)
	apisubdomain.GET("/discover", queryDiscover)
	apisubdomain.GET("/cateindex", queryCateIndex)
	apisubdomain.GET("/hotwords", queryHotWords)
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
			var games []gin.H
			var ids []string
			for _, e := range *r {
				games = append(games, formGamePrice(c, e))
				ids = append(ids, e.Id)
			}
			if games != nil {
				d = gin.H {
					"games" : games,
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
		pageNo = 0;
	}
	startPos := pageSize * pageNo

	r, err := db.QueryCheapGames(startPos, pageSize)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		d := gin.H{}
		if r != nil {
			var games []gin.H
			var ids []string
			for _, e := range *r {
				games = append(games, formGamePrice(c, e))
				ids = append(ids, e.Id)
			}
			if games != nil {
				d = gin.H {
					"games" : games,
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
			var games []gin.H
			var ids []string
			for _, e := range *r {
				games = append(games, formGamePrice(c, e))
				ids = append(ids, e.Id)
			}
			if games != nil {
				d = gin.H {
					"games" : games,
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
			"hotwords" : hwd.Words,
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