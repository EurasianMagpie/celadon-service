package db

import "fmt"
import "log"

import "database/sql"
import _ "database/sql/driver"
import _ "github.com/go-sql-driver/mysql"

import "github.com/EurasianMagpie/celadon/config"

/*
func Getdb() {
	dbcfg := config.GetConfig().Db
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbcfg.User, dbcfg.Pass, dbcfg.Host, dbcfg.Name)
	fmt.Println("Getdb | DSN:", dsn)
	d, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	defer d.Close()

	sel, err := d.Query("select region_id, name, cname, logo from region")
	if err != nil {
		panic(err.Error())
	}
	for sel.Next() {
		var region Region
		err = sel.Scan(&region.region_id, &region.name, &region.cname, &region.logo)
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%s %s %s\n", region.region_id, region.name, region.cname)
	}

	defer sel.Close()

}//*/

var edb *sql.DB
var stmtQueryRegion *sql.Stmt

func getdb() *sql.DB {
	if edb != nil {
		return edb
	}
	dbcfg := config.GetConfig().Db
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbcfg.User, dbcfg.Pass, dbcfg.Host, dbcfg.Name)
	fmt.Println("Getdb | DSN:", dsn)
	edb, err := sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}
	return edb
}

func QueryRegionInfo(id string) *Region {
	d := getdb()
	if d == nil {
		return nil
	}

	if stmtQueryRegion == nil {
		stmt, err := d.Prepare("select region_id, name, cname, logo from region where region_id = ?")
		if err != nil {
			log.Fatal(err)
		}
		stmtQueryRegion = stmt
	}
	var region Region
	err := stmtQueryRegion.QueryRow(id).Scan(&region.Region_id, &region.Name, &region.Cname, &region.Logo)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s %s %s\n", region.Region_id, region.Name, region.Cname)
	return &region
}