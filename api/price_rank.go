package api

import "strings"
import "fmt"
import "strconv"
import "container/list"

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/db"

type regionPrice struct {
	Index int
	FP float64
	Price string
}

var rankLength = 10
func calcRegionPriceRank(prices string) ([]regionPrice, error) {
	var priceList []regionPrice
	rank := list.New()
	s := strings.Split(prices, ",")
	for i, e := range s {
		f, err := strconv.ParseFloat(e, 64)
		if err != nil {
			f = 0.0
		}
		cur := regionPrice{i, f, e}

		_insert := false
		for r := rank.Front(); r != nil; r = r.Next() {
			er := r.Value.(regionPrice)
			if cur.FP < er.FP {
				rank.InsertBefore(cur, r)
				_insert = true
				break
			}
		}
		if !_insert {
			rank.PushBack(cur)
		}
	}
	
	for r := rank.Front(); r != nil; r = r.Next() {
		er := r.Value.(regionPrice)
		priceList = append(priceList, er)
		if len(priceList) == rankLength - 1 {
			break
		}
	}
	if rank.Len() >= rankLength {
		priceList = append(priceList, rank.Back().Value.(regionPrice))
	}
	fmt.Println(priceList)
	return priceList, nil
}

func formPriceRank(priceList []regionPrice) []gin.H {
	var r []gin.H
	for _, p := range priceList {
		region, err := db.FindRegionByIndex(p.Index)
		regionName := ""
		if err == nil {
			regionName = region.Cname
		}
		r = append(r, gin.H{
			"price": p.Price,
			"region": regionName,
		})
	}
	return r
}

func QueryPriceRank(c *gin.Context) {
	id := c.Query("id")
	if len(id) == 0 {
		c.JSON(200, formResult(301, string("invalid param id"), gin.H{}))
		return
	}
	r, err := db.QueryGameFullPrice(id)
	if err != nil {
		c.JSON(200, formResult(300, string(err.Error()), gin.H{}))
	} else {
		rankData, err := calcRegionPriceRank(r)
		d := gin.H{}
		if err == nil {
			if (rankData != nil) {
				d = gin.H {
					"prices": formPriceRank(rankData),
				}
			}
		}
		c.JSON(200, formResult(0, "", d))
	}
}