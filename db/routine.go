package db

import "strconv"
import "time"
import "container/list"


////////////////////////////////////////////////////////////
var mapCheckGameDetail map[string]bool

func ReCheckGameDetail() bool {
	d := getdb()
	if d == nil {
		return false
	}

	var m map[string]bool
	m = make(map[string]bool)
	rows, err := d.Query("select game_id, DATE(release_date), description from game")
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var id, date, desc string
		err := rows.Scan(&id, &date, &desc)
		if err != nil {
			return false
		}
		m[id] = len(desc)>0
		if DefaultReleaseDate == date && id > "1900" {
			m[id] = false
		}
		//fixGameDesc(id, desc)	// done
		//fixUnEscapedDesc(id, desc)	// done
	}
	mapCheckGameDetail = m
	return true
}

func IsGameDetialed(id string) bool {
	if val, ok := mapCheckGameDetail[id]; ok {
		return val
	}
	return false
}

func MarkGameDetailed(id string) {
	mapCheckGameDetail[id] = true
}

func FindAnyUnDetailedGames(count int) []string {
	var ids []string
	for k, v := range mapCheckGameDetail {
		if !v {
			ids = append(ids, k)
			if len(ids) > count {
				break
			}
		}
	}
	return ids
}

////////////////////////////////////////////////////////////
var curCheapGames []GamePrice
var checkCheapTime time.Time

type cheapGame struct {
	CheapRate float64
	GP GamePrice
}

func ReCheckCheapGames() bool {
	d := getdb()
	if d == nil {
		return false
	}

	rows, err := d.Query(`
	select 
		price.game_id, t1.name, t1.cname, t1.cover, price.lregion, price.lprice, price.islowest, price.lowestprice 
	from 
		price 
		inner join 
			(select game_id, name, cname, cover from game where cname!="") as t1 
		on price.game_id=t1.game_id
	order by t1.cname
	`)
	if err != nil {
		return false
	}
	defer rows.Close()

	var cheapGames []GamePrice
	cheapList := list.New()
	for rows.Next() {
		var p GamePrice
		var lowestPrice = ""
		err := rows.Scan(&p.Id, &p.Name, &p.Cname, &p.Cover, &p.Region, &p.Price, &p.IsLowest, &lowestPrice)
		if err != nil {
			return false
		}
		
		pr, err := strconv.ParseFloat(p.Price, 64)
		if err != nil || pr < 0.01 {
			continue
		}

		if (p.IsLowest == 1) {
			cheapGames = append(cheapGames, p)
		} else {
			lpr, err := strconv.ParseFloat(lowestPrice, 64)
			if err != nil {
				continue
			}
			if lpr < 0.01 {
				lpr = 0.01
			}
			lrate := (pr - lpr) / lpr
			chgm := cheapGame{CheapRate:lrate, GP:p}

			_insert := false
			for r := cheapList.Front(); r != nil; r = r.Next() {
				er := r.Value.(cheapGame)
				if chgm.CheapRate < er.CheapRate {
					cheapList.InsertBefore(chgm, r)
					_insert = true
					break
				}
			}
			if !_insert {
				cheapList.PushBack(chgm)
			}
		}
	}

	for r := cheapList.Front(); r != nil; r = r.Next() {
		er := r.Value.(cheapGame)
		cheapGames = append(cheapGames, er.GP)
	}

	curCheapGames = cheapGames
	checkCheapTime = time.Now()

	return true
}

func QueryCheapGames(startPos int, pageSize int) (*[]GamePrice, error) {
	duration := time.Since(checkCheapTime)
	if duration.Hours() > 6 {
		go ReCheckCheapGames()
		checkCheapTime = time.Now()
	}

	var gamePrices []GamePrice
	if startPos < len(curCheapGames) {
		if startPos + pageSize < len(curCheapGames) {
			gamePrices = curCheapGames[startPos:startPos+pageSize]
		} else {
			gamePrices = curCheapGames[startPos:]
		}
	}
	return &gamePrices, nil
}