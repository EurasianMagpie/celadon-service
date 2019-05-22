package db

import "time"

type Region struct {
	Name string
	Cname string
	Region_id string
	Logo string
}

func NewRegion(id string, name string) Region {
	region := Region{Name:name, Cname:"",Region_id:id,Logo:""}
	return region
}

type GameInfo struct {
	Id string
	Name string
	Cname string
	Publisher string
	ReleaseDate time.Time
	Desc string
	Cover string
	Language string
	Tags string
	Ref string	
	RealCard int
	Status int

	CoverUrl string
	CoverType string
}

var DefaultReleaseDate = "2018-01-01"

func NewGameInfo(id string, name string, ref string) GameInfo {
	gameInfo := GameInfo {
		Id:id, Name:name, Cname:"", Publisher:"", Desc:"", Cover:"", Language:"", Tags:"", RealCard:0, Status:0,
	}
	dt, _ := time.Parse("2006-01-02", DefaultReleaseDate)
	gameInfo.ReleaseDate = dt
	gameInfo.Ref = ref
	return gameInfo
}

type Price struct {
	Id string
	Price string
	Discount string
	LPrice string
	LRegion string
	HPrice string
	HRegion string
	LowestPrice string
	LowestRegion string
	IsLowest int
}

func NewPrice(id string, pr string, lp string, lr string, hp string, hr string) Price {
	price := Price{Id:id, Price:pr, Discount:"", LPrice:lp, LRegion:lr, HPrice:hp, HRegion:hr, LowestPrice:"", LowestRegion:"", IsLowest:0}
	return price
}

type GamePrice struct {
	Id string
	Name string
	Cname string
	Cover string
	Region string
	Price string
	IsLowest int
}