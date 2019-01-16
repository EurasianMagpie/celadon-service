package db

import "fmt"
import "database/sql"
import _ "database/sql/driver"
import _ "github.com/go-sql-driver/mysql"

import "github.com/EurasianMagpie/celadon/config"


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

}

var edb *sql.DB
