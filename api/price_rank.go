package api

import "strings"
import "fmt"
import "strconv"
import "container/list"

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon-service/db"

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
			continue
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