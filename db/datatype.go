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
	Desc string
	Language string
	Cover string
	ReleaseDate time.Time
	Status int

	Ref string
	CoverUrl string
}

func NewGameInfo(id string, name string, ref string) GameInfo {
	gameInfo := GameInfo{Id:id,Name:name,Cname:"",Desc:"",Language:"",Cover:"",Status:0}
	dt, _ := time.Parse("2006-01-02", "2018-01-01")
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
}

func NewPrice(id string, pr string, lp string, lr string, hp string, hr string) Price {
	price := Price{Id:id, Price:pr, Discount:"", LPrice:lp, LRegion:lr, HPrice:hp, HRegion:hr}
	return price
}

type GamePrice struct {
	Id string
	Name string
	Region string
	Price string
}