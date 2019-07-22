package api

import "github.com/EurasianMagpie/celadon-service/db"

const (
	ContentType_GamePrice 	= "1"
	ContentType_H5			= "2"
	ContentType_AD			= "9"
)

type ContentItem struct {
	Type string
	Data interface{}
}

func NewContentItemGamePrice(gamePrice db.GamePrice) ContentItem {
	item := ContentItem{Type:ContentType_GamePrice, Data:gamePrice}
	return item
}