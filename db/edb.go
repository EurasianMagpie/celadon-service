package db

import "fmt"
//import "log"
import "errors"
import "strconv"

import "database/sql"
import _ "database/sql/driver"
import _ "github.com/go-sql-driver/mysql"

import "github.com/EurasianMagpie/celadon/config"


var edb *sql.DB
var stmtQueryRegion *sql.Stmt
var stmtQueryGameInfo *sql.Stmt
var stmtQueryPriceInfo *sql.Stmt
var stmtQueryGamePrice *sql.Stmt
var stmtQueryGameFullPrice *sql.Stmt
var stmtQuerySearchGamePrice *sql.Stmt
var stmtQueryRecommend *sql.Stmt
var stmtQueryMultiPrice *sql.Stmt

var stmtUpdateRegion *sql.Stmt
var stmtUpdateGame *sql.Stmt
var stmtUpdatePrice *sql.Stmt
var stmtUpdateQueryLowestPrice *sql.Stmt
var stmtUpdateLowestPrice *sql.Stmt
var stmtUpdateLowestFlag *sql.Stmt

var stmtUpdateGameCname *sql.Stmt

func init() {
	initAllStmts()
}

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

func initAllStmts() {
	d := getdb()
	if d == nil {
		return
	}

	if stmtQueryRegion == nil {
		stmt, err := d.Prepare("select region_id, name, cname from region where region_id = ?")
		if err == nil {
			stmtQueryRegion = stmt
		}
	}

	if stmtQueryGameInfo == nil {
		stmt, err := d.Prepare(`select game_id, name, cname, ref, description, language, cover, status from game where game_id=?`)
		if err == nil {
			stmtQueryGameInfo = stmt
		}
	}

	if stmtQueryPriceInfo == nil {
		stmt, err := d.Prepare(`select game_id, price, discount, lprice, lregion, hprice, hregion from price where game_id=?`)
		if err == nil {
			stmtQueryPriceInfo = stmt
		}
	}

	if stmtQueryGamePrice == nil {
		stmt, err := d.Prepare(`
		select 
			price.game_id, t1.name, t1.cname, t1.cover, price.lregion, price.lprice 
		from
			price
			inner join
				(select game_id, name, cname, cover from game where game_id=?) as t1
			on price.game_id=t1.game_id
		`)
		if err == nil {
			stmtQueryGamePrice = stmt
		}
	}

	if stmtQueryGameFullPrice == nil {
		stmt, err := d.Prepare(`
		select 
			price 
		from
			price where game_id=?`)
		if err == nil {
			stmtQueryGameFullPrice = stmt
		}
	}

	if stmtQuerySearchGamePrice == nil {
		stmt, err := d.Prepare(`
		select 
			price.game_id, t1.name, t1.cname, t1.cover, price.lregion, price.lprice 
		from 
			price 
			inner join 
				(select game_id, name, cname, cover from game where name like ? or cname like ?) as t1 
			on price.game_id=t1.game_id
		order by t1.cname=""
		`)
		if err == nil {
			stmtQuerySearchGamePrice = stmt
		}
	}

	if stmtQueryRecommend == nil {
		stmt, err := d.Prepare(`
		select 
			price.game_id, t1.name, t1.cname, t1.cover, price.lregion, price.lprice 
		from 
			price 
			inner join 
				(select game_id, name, cname, cover from game where cname!="" order by rand() limit ?) as t1 
			on price.game_id=t1.game_id
		order by t1.cname
		`)
		if err == nil {
			stmtQueryRecommend = stmt
		}
	}

	// todo optmz ...
	if stmtQueryMultiPrice == nil {
		stmt, err := d.Prepare(`
		select 
			price.game_id, t1.name, t1.cname, t1.cover, price.lregion, price.lprice 
		from 
			price 
			inner join 
				(select game_id, name, cname, cover from game where game_id=?) as t1 
			on price.game_id=t1.game_id
		`)
		if err == nil {
			stmtQueryMultiPrice = stmt
		}
	}

	if stmtUpdateRegion == nil {
		//fmt.Println("create stmtUpdateRegion")
		stmt, err := d.Prepare(`
			INSERT INTO region (region_id,name,cname,logo) 
			VALUES(?,?,?,?) 
			ON DUPLICATE KEY UPDATE name=?
		`)
		if err == nil {
			//log.Fatal(err)
			stmtUpdateRegion = stmt
		}
	}

	if stmtUpdateGame == nil {
		//fmt.Println("create stmtUpdateGame")
		stmt, err := d.Prepare(`
			INSERT INTO game (game_id, name, cname, ref, description, language, cover, release_date, status) 
			VALUES(?,?,?,?,?,?,?,?,?) 
			ON DUPLICATE KEY UPDATE name=?, ref=?, description=?, release_date=?
		`)
		if err == nil {
			//log.Fatal(err)
			stmtUpdateGame = stmt
		}
	}

	if stmtUpdatePrice == nil {
		//fmt.Println("create stmtUpdatePrice")
		stmt, err := d.Prepare("INSERT INTO price (game_id,price,discount,lprice,lregion,hprice,hregion) VALUES(?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE price=?,lprice=?,lregion=?,hprice=?,hregion=?")
		if err == nil {
			//log.Fatal(err)
			stmtUpdatePrice = stmt
		}
	}

	if stmtUpdateQueryLowestPrice == nil {
		//fmt.Println("create stmtUpdateQueryLowestPrice")
		stmt, err := d.Prepare("select lowestprice from price where game_id=?")
		if err == nil {
			//log.Fatal(err)
			stmtUpdateQueryLowestPrice = stmt
		}
	}
	
	if stmtUpdateLowestPrice == nil {
		//fmt.Println("create stmtUpdateLowestPrice")
		stmt, err := d.Prepare("UPDATE price SET lowestprice=?, lowestregion=?, islowest=1 where game_id=?")
		if err == nil {
			//log.Fatal(err)
			stmtUpdateLowestPrice = stmt
		}
	}

	if stmtUpdateLowestFlag == nil {
		//fmt.Println("create stmtUpdateLowestFlag")
		stmt, err := d.Prepare("UPDATE price SET islowest=? where game_id=?")
		if err == nil {
			//log.Fatal(err)
			stmtUpdateLowestFlag = stmt
		}
	}
}

func QueryRegionInfo(id string) (*Region, error) {
	d := getdb()
	if d == nil {
		return nil, errors.New("db error")
	}

	if stmtQueryRegion == nil {
		return nil, errors.New("db stmt init failed")
	}
	var region Region
	err := stmtQueryRegion.QueryRow(id).Scan(&region.Region_id, &region.Name, &region.Cname)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s %s %s\n", region.Region_id, region.Name, region.Cname)
	return &region, nil
}

func QueryGameInfo(id string) (*GameInfo, error) {
	d := getdb()
	if d == nil {
		return nil, errors.New("db error")
	}
	if stmtQueryGameInfo == nil {
		return nil, errors.New("db stmt init failed")
	}
	var g GameInfo
	err := stmtQueryGameInfo.QueryRow(id).Scan(&g.Id, &g.Name, &g.Cname, &g.Ref, &g.Desc, &g.Language, &g.Cover, &g.Status)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

func QueryPriceInfo(id string) (*Price, error) {
	d := getdb()
	if d == nil {
		return nil, errors.New("db error")
	}
	if stmtQueryPriceInfo == nil {
		return nil, errors.New("db stmt init failed")
	}
	var p Price
	err := stmtQueryPriceInfo.QueryRow(id).Scan(&p.Id, &p.Price, &p.Discount, &p.LPrice, &p.LRegion, &p.HPrice, &p.HRegion)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func QueryGamePrice(id string) (*GamePrice, error) {
	d := getdb()
	if d == nil {
		return nil, errors.New("db error")
	}
	if stmtQueryGamePrice == nil {
		return nil, errors.New("db stmt init failed")
	}
	var gamePrice GamePrice
	err := stmtQueryGamePrice.QueryRow(id).Scan(&gamePrice.Id, &gamePrice.Name, &gamePrice.Cname, &gamePrice.Cover, &gamePrice.Region, &gamePrice.Price)
	if err != nil {
		return nil, err
	}
	return &gamePrice, nil
}

func QueryGameFullPrice(id string) (string, error) {
	d := getdb()
	if d == nil {
		return "", errors.New("db error")
	}
	if stmtQueryGameFullPrice == nil {
		return "", errors.New("db stmt init failed")
	}
	var price string
	err := stmtQueryGameFullPrice.QueryRow(id).Scan(&price)
	if err != nil {
		return "", err
	}
	return price, nil
}

func QuerySearchGamePrice(name string) (*[]GamePrice, error) {
	d := getdb()
	if d == nil {
		return nil, errors.New("db error")
	}
	if stmtQuerySearchGamePrice == nil {
		return nil, errors.New("db stmt init failed")
	}
	var gamePrices []GamePrice
	np := "%" + name + "%"
	rows, err := stmtQuerySearchGamePrice.Query(np, np)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p GamePrice
		err := rows.Scan(&p.Id, &p.Name, &p.Cname, &p.Cover, &p.Region, &p.Price)
		if err != nil {
			return nil, errors.New("scan error")
		}
		//fmt.Println(p)
		gamePrices = append(gamePrices, p)
	}
	return &gamePrices, nil
}

func QueryRecommendGames(limit int) (*[]GamePrice, error) {
	d := getdb()
	if d == nil {
		return nil, errors.New("db error")
	}
	if stmtQueryRecommend == nil {
		return nil, errors.New("db stmt init failed")
	}
	var gamePrices []GamePrice
	rows, err := stmtQueryRecommend.Query(limit)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var p GamePrice
		err := rows.Scan(&p.Id, &p.Name, &p.Cname, &p.Cover, &p.Region, &p.Price)
		if err != nil {
			return nil, errors.New("scan error")
		}
		//fmt.Println(p)
		gamePrices = append(gamePrices, p)
	}
	return &gamePrices, nil
}

func QueryPriceListByIds(ids []string) (*[]GamePrice, error) {
	d := getdb()
	if d == nil {
		return nil, errors.New("db error")
	}
	// todo optmz ...
	if stmtQueryMultiPrice == nil {
		return nil, errors.New("db stmt init failed")
	}
	var gamePrices []GamePrice
	for _, id := range ids {
		var p GamePrice
		err := stmtQueryMultiPrice.QueryRow(id).Scan(&p.Id, &p.Name, &p.Cname, &p.Cover, &p.Region, &p.Price)
		if err != nil {
			continue
		}
		gamePrices = append(gamePrices, p)
	}
	return &gamePrices, nil
}

func UpdateRegion(region Region) bool {
	d := getdb()
	if d == nil {
		return false
	}

	if stmtUpdateRegion == nil {
		return false
	}
	_, err := stmtUpdateRegion.Exec(region.Region_id, region.Name, region.Cname, region.Logo, region.Name)
	if err != nil {
		panic(err)
	}
	return nil == err
}

func UpdateGame(gameInfo GameInfo) bool {
	d := getdb()
	if d == nil {
		return false
	}

	if stmtUpdateGame == nil {
		return false
	}
	_, err := stmtUpdateGame.Exec(gameInfo.Id, gameInfo.Name, gameInfo.Cname, gameInfo.Ref, gameInfo.Desc, gameInfo.Language, gameInfo.Cover, gameInfo.ReleaseDate, gameInfo.Status, gameInfo.Name, gameInfo.Ref, gameInfo.Desc, gameInfo.ReleaseDate)
	if err != nil {
		panic(err)
	}
	return err == nil
}

func UpdatePrice(price Price) bool {
	d := getdb()
	if d == nil {
		return false
	}

	if stmtUpdatePrice == nil {
		return false
	}
	//fmt.Println("UpdatePrice stmtUpdatePrice")
	_, err := stmtUpdatePrice.Exec(price.Id, price.Price, price.Discount, price.LPrice, price.LRegion, price.HPrice, price.HRegion, price.Price, price.LPrice, price.LRegion, price.HPrice, price.HRegion)
	if err != nil {
		//panic(err)
		return false
	}

	return UpdateLowestPrice(price.Id, price.LPrice, price.LRegion)
}

func UpdateGameCname(id string, cname string) bool {
	d := getdb()
	if d == nil {
		return false
	}

	if stmtUpdateGameCname == nil {
		//fmt.Println("create stmtUpdateGameCname")
		stmt, err := d.Prepare(`
			UPDATE game SET cname=? where game_id=?
		`)
		if err != nil {
			//log.Fatal(err)
			return false
		}
		stmtUpdateGameCname = stmt
	}
	_, err := stmtUpdateGameCname.Exec(cname, id)
	//if err != nil {
	//	panic(err)
	//}
	return err == nil
}

func UpdateLowestPrice(id string, price string, region string) bool {
	d := getdb()
	if d == nil {
		return false
	}

	if stmtUpdateQueryLowestPrice == nil {
		return false
	}
	//fmt.Println("UpdateLowestPrice stmtUpdateQueryLowestPrice")
	var lowest string
	err := stmtUpdateQueryLowestPrice.QueryRow(id).Scan(&lowest)
	if err != nil {
		panic(err)
		return false
	}
	hasLowest := false
	if (len(lowest) > 0) {
		hasLowest = true
	}

	//fmt.Println("UpdateLowestPrice hasLowest:", hasLowest)

	if stmtUpdateLowestPrice == nil {
		return false
	}
	if stmtUpdateLowestFlag == nil {
		return false
	}
	if hasLowest {
		lowest_old, err := strconv.ParseFloat(lowest, 32)
		if err != nil {
			return false
		}
		lowest_new, err := strconv.ParseFloat(price, 32)
		if err != nil {
			return false
		}

		//fmt.Println("UpdateLowestPrice hasLowest, old, new", lowest_old, lowest_new)

		if lowest_new <= lowest_old {
			//fmt.Println("UpdateLowestPrice stmtUpdateLowestPrice")
			_, err := stmtUpdateLowestPrice.Exec(price, region, id)
			if err != nil {
				return false
			}
		} else {
			//fmt.Println("UpdateLowestPrice stmtUpdateLowestFlag")
			_, err := stmtUpdateLowestFlag.Exec(0, id)
			if err != nil {
				return false
			}
		}
	} else {
		//fmt.Println("UpdateLowestPrice stmtUpdateLowestPrice")
		_, err := stmtUpdateLowestPrice.Exec(price, region, id)
		if err != nil {
			return false
		}
	}
	return true
}