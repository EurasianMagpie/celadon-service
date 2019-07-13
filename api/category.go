package api

import (
	"strconv"
	"errors"
)

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"

type CateHandlerFunc func(int, int) (*[]db.GamePrice, error)

type Cate struct {
	Name string
	Type string
	Param string
	IsDefault int

	Handler CateHandlerFunc
}

var discoverCates = [...]Cate {
	Cate{"最新发布", "gplist", "latest", 1, queryLatestHandler},
	Cate{"经典游戏", "gplist", "classic", 0, queryClassicHandler},
}

var mapDiscoverCates map[string]Cate
var cateIndex gin.H

func init() {
	mapDiscoverCates = make(map[string]Cate)
	var cates []gin.H
	for _, cate := range discoverCates {
		mapDiscoverCates[cate.Param] = cate
		cates = append(cates, gin.H{
			"name" : cate.Name,
			"type" : cate.Type,
			"param" : cate.Param,
			"default" : cate.IsDefault,
		})
	}
	cateIndex = gin.H{
		"cates" : cates,
	}
}

func queryDiscoverCateIndex(c *gin.Context) {
	c.JSON(200, formResult(0, "", cateIndex))
}

func queryDiscover(c *gin.Context) {
	cate := c.Query("cate")
	if len(cate) == 0 {
		c.JSON(200, formResult(301, string("invalid param cate"), gin.H{}))
		return
	}

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

	var handler CateHandlerFunc = unknownCateHandler
	if val, ok := mapDiscoverCates[cate]; ok {
		handler = val.Handler
	} 
	r, err := handler(startPos, pageSize)

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
				d = gin.H {
					"items" : items,
				}
			}
			invokeIpcTask(ids)
		}
		c.JSON(200, formResult(0, "", d))
	}
}

func queryLatestHandler(startPos int, pageSize int) (*[]db.GamePrice, error) {
	if startPos <= 100 {
		return db.QueryLatestGames(startPos, pageSize)
	} else {
		return nil, nil
	}
}

func queryClassicHandler(startPos int, pageSize int) (*[]db.GamePrice, error) {
	return db.QueryRealCardGames(startPos, pageSize)
}

func unknownCateHandler(startPos int, pageSize int) (*[]db.GamePrice, error) {
	return nil, errors.New("unknown cate")
}