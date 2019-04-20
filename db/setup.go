package db

import "errors"
import "fmt"
import "log"
import "database/sql"
import _ "database/sql/driver"
import _ "github.com/go-sql-driver/mysql"

import "github.com/EurasianMagpie/celadon/config"



func init() {
	setupDB()
}

func setupDB() {
	dbcfg := config.GetConfig().Db
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/", dbcfg.User, dbcfg.Pass, dbcfg.Host)
	fmt.Println("setupDB | DSN:", dsn)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	c := "CREATE DATABASE IF NOT EXISTS %s CHARACTER SET utf8 COLLATE utf8_general_ci;"
	s := fmt.Sprintf(c, dbcfg.Name)
	_, err = db.Exec(s)
	if err != nil {
		panic(err)
	}

	_, err = db.Exec("USE " + dbcfg.Name)
	if err != nil {
		panic(err)
	}

	c = `create table if not exists region (
        region_id VARCHAR(16) NOT NULL,
        name VARCHAR(255) NOT NULL DEFAULT '',
        cname VARCHAR(255) NOT NULL DEFAULT '',
        logo TEXT,
        PRIMARY KEY (region_id)
		);`
	_, err = db.Exec(c)
	if err != nil {
		panic(err)
	}

	c = `create table if not exists game (
        game_id INT NOT NULL,
        name VARCHAR(255) NOT NULL DEFAULT '',
        cname VARCHAR(255) NOT NULL DEFAULT '',
        ref TEXT,
        description TEXT,
        language TEXT,
        cover TEXT,
        release_date DATE,
        status INT NOT NULL DEFAULT 0,
        PRIMARY KEY (game_id)
        );`
	_, err = db.Exec(c)
	if err != nil {
		panic(err)
	}

	c = `create table if not exists price (
        game_id INT NOT NULL,
        price TEXT NOT NULL,
        discount TEXT,
        lprice VARCHAR(255) NOT NULL DEFAULT '',
        lregion VARCHAR(16) NOT NULL DEFAULT '',
        hprice VARCHAR(255) NOT NULL DEFAULT '',
        hregion VARCHAR(16) NOT NULL DEFAULT '',
        lowestprice VARCHAR(255) NOT NULL DEFAULT '',
        lowestregion VARCHAR(255) NOT NULL DEFAULT '',
        islowest INT NOT NULL DEFAULT 0,
        PRIMARY KEY (game_id)
		);`
	_, err = db.Exec(c)
	if err != nil {
		panic(err)
	}
	
	updateRegion(db)
}

var regions = [...]Region {
	Region{"Australia","澳大利亚","AUS",""},
    Region{"Austria","奥地利","AUT",""},
    Region{"Belgium","比利时","BEL",""},
    Region{"Bulgaria","保加利亚","BGR",""},
    Region{"Canada","加拿大","CAN",""},
    Region{"Croatia","克罗地亚","HRV",""},
    Region{"Cyprus","塞浦路斯","CYP",""},
    Region{"Czech Republic","捷克","CZE",""},
    Region{"Denmark","丹麦","DNK",""},
    Region{"Estonia","爱沙尼亚","EST",""},
    Region{"Finland","芬兰","FIN",""},
    Region{"France","法国","FRA",""},
    Region{"Germany","德国","DEU",""},
    Region{"Greece","希腊","GRC",""},
    Region{"Hungary","匈牙利","HUN",""},
    Region{"Ireland","爱尔兰","IRL",""},
    Region{"Italy","意大利","ITA",""},
    Region{"Japan","日本","JPN",""},
    Region{"Latvia","拉脱维亚","LVA",""},
    Region{"Lithuania","立陶宛","LTU",""},
    Region{"Luxembourg","卢森堡","LUX",""},
    Region{"Malta","马耳他","MLT",""},
    Region{"Mexico","墨西哥","MEX",""},
    Region{"Netherlands","荷兰","NLD",""},
    Region{"New Zealand","新西兰","NZL",""},
    Region{"Norway","挪威","NOR",""},
    Region{"Poland","波兰","POL",""},
    Region{"Portugal","葡萄牙","PRT",""},
    Region{"Romania","罗马尼亚","ROU",""},
    Region{"Russia","俄罗斯","RUS",""},
    Region{"Slovakia","斯洛伐克","SVK",""},
    Region{"Slovenia","斯洛文尼亚","SVN",""},
    Region{"South Africa","南非","ZAF",""},
    Region{"Spain","西班牙","ESP",""},
    Region{"Sweden","瑞典","SWE",""},
    Region{"Switzerland","瑞士","CHE",""},
    Region{"United Kingdom","英国","GBR",""},
    Region{"United States","美国","USA",""},
}

func updateRegion(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO region (name,cname,region_id,logo) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE name=?, cname=?, logo=?;")
	if err != nil {
		log.Fatal(err)
	}
	for _, r := range regions {
		_, err := stmt.Exec(r.Name, r.Cname, r.Region_id, r.Logo, r.Name, r.Cname, r.Logo)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func FindRegionByAbbr(abbr string) (*Region, error) {
    for _, r := range regions {
        if r.Region_id == abbr {
            return &r, nil
        }
    }
    return nil, errors.New("not found")
}

func FindRegionByIndex(index int) (*Region, error) {
    if index < len(regions) {
        return &regions[index], nil
    }
    return nil, errors.New("not found")
}