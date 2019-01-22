package db

type Region struct {
	Name string
	Cname string
	Region_id string
	Logo string
}

type GamePrice struct {
	Id string
	Name string
	Region string
	Price string
}