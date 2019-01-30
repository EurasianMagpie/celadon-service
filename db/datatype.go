package db

type Region struct {
	Name string
	Cname string
	Region_id string
	Logo string
}

type GameInfo struct {
	Id string
	Name string
	Cname string
	Desc string
	Language string
	Cover string
	ReleaseDate string
	Status int

	Ref string
	CoverUrl string
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

type GamePrice struct {
	Id string
	Name string
	Region string
	Price string
}