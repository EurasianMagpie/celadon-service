package db

import "fmt"
//import "log"
import "errors"

import "database/sql"
import _ "database/sql/driver"
import _ "github.com/go-sql-driver/mysql"

import "github.com/EurasianMagpie/celadon/config"


var edb *sql.DB
var stmtQueryRegion *sql.Stmt
var stmtQueryGamePrice *sql.Stmt
var stmtQuerySearchGamePrice *sql.Stmt

var stmtUpdateRegion *sql.Stmt
var stmtUpdateGame *sql.Stmt
var stmtUpdatePrice *sql.Stmt

func getdb() *sql.DB {
	if edb != nil {
		return edb
	}
	dbcfg := config.GetConfig().Db
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbcfg.User, dbcfg.Pass, dbcfg.Host, dbcfg.Name)
	fmt.Println("Getdb | DSN:", dsn)
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	edb = d
	return edb
}

func QueryRegionInfo(id string) (*Region, error) {
	d := getdb()
	if d == nil {
		return nil, errors.New("db error")
	}

	if stmtQueryRegion == nil {
		stmt, err := d.Prepare("select region_id, name, cname from region where region_id = ?")
		if err != nil {
			stmtQueryRegion = nil
			return nil, err
		}
		stmtQueryRegion = stmt
	}
	var region Region
	err := stmtQueryRegion.QueryRow(id).Scan(&region.Region_id, &region.Name, &region.Cname)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s %s %s\n", region.Region_id, region.Name, region.Cname)
	return &region, nil
}

func QueryGamePrice(id string) (*GamePrice, error) {
	d := getdb()
	if d == nil {
		return nil, errors.New("db error")
	}
	if stmtQueryGamePrice == nil {
		stmt, err := d.Prepare(`select game.game_id, game.name, rpt.rgn, rpt.lp
		from game,
		(select region.cname as rgn, pt.lprice as lp
		from region,
		(select game_id, lregion, lprice from price where game_id=?) as pt
		where region.region_id=pt.lregion) as rpt
		where game.game_id=?`)
		if err != nil {
			return nil, err
		}
		stmtQueryGamePrice = stmt
	}
	var gamePrice GamePrice
	err := stmtQueryGamePrice.QueryRow(id, id).Scan(&gamePrice.Id, &gamePrice.Name, &gamePrice.Region, &gamePrice.Price)
	if err != nil {
		return nil, err
	}
	return &gamePrice, nil
}

func QuerySearchGamePrice(name string) (*[]GamePrice, error) {
	d := getdb()
	if d == nil {
		return nil, errors.New("db error")
	}
	if stmtQuerySearchGamePrice == nil {
		stmt, err := d.Prepare(`
		select 
			price.game_id, t1.name, price.lregion, price.lprice 
		from 
			price 
			inner join 
				(select game_id, name from game where name like ?) as t1 
			on price.game_id=t1.game_id
		`)
		if err != nil {
			return nil, err
		}
		stmtQuerySearchGamePrice = stmt
	}
	var gamePrices []GamePrice
	np := "%" + name + "%"
	rows, err := stmtQuerySearchGamePrice.Query(np)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p GamePrice
		err := rows.Scan(&p.Id, &p.Name, &p.Region, &p.Price)
		if err != nil {
			return nil, errors.New("scan error")
		}
		//fmt.Println(p)
		gamePrices = append(gamePrices, p)
	}
	return &gamePrices, nil
}

func UpdateRegion(region Region) {
	d := getdb()
	if d == nil {
		return
	}

	if stmtUpdateRegion == nil {
		fmt.Println("create stmtUpdateRegion")
		stmt, err := d.Prepare("INSERT INTO region (region_id,name) VALUES(?,?) ON DUPLICATE KEY UPDATE name=?")
		if err != nil {
			//log.Fatal(err)
			return
		}
		stmtUpdateRegion = stmt
	}
	_, err := stmtUpdateRegion.Exec(region.Region_id, region.Name, region.Name)
	if err != nil {
		panic(err)
	}
}

func UpdateGame(gameInfo GameInfo) {
	d := getdb()
	if d == nil {
		return
	}

	if stmtUpdateGame == nil {
		fmt.Println("create stmtUpdateGame")
		stmt, err := d.Prepare("INSERT INTO game (game_id, name, description, release_date) VALUES(?,?,?,?) ON DUPLICATE KEY UPDATE name=?, description=?, release_date=?")
		if err != nil {
			//log.Fatal(err)
			return
		}
		stmtUpdateGame = stmt
	}
	date := ""
	if len(gameInfo.ReleaseDate) > 0 {
		date = "STR_TO_DATE('" + gameInfo.ReleaseDate + "', %M %D, %Y')"
	}
	_, err := stmtUpdateGame.Exec(gameInfo.Id, gameInfo.Name, gameInfo.Desc, date, gameInfo.Name, gameInfo.Desc, date)
	if err != nil {
		panic(err)
	}
}

func UpdatePrice(price Price) {
	d := getdb()
	if d == nil {
		return
	}

	if stmtUpdatePrice == nil {
		fmt.Println("create stmtUpdatePrice")
		stmt, err := d.Prepare("INSERT INTO price (game_id,price,lprice,lregion,hprice,hregion) VALUES(?,?,?,?,?,?) ON DUPLICATE KEY UPDATE price=?,lprice=?,lregion=?,hprice=?,hregion=?")
		if err != nil {
			//log.Fatal(err)
			return
		}
		stmtUpdatePrice = stmt
	}
	_, err := stmtUpdatePrice.Exec(price.Id, price.Price, price.LPrice, price.LRegion, price.HPrice, price.HRegion, price.Price, price.LPrice, price.LRegion, price.HPrice, price.HRegion)
	if err != nil {
		panic(err)
	}
}