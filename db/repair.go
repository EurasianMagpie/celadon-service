package db

import "strings"
import "database/sql"
import "fmt"


var stmtFixGameDesc *sql.Stmt

func fixGameDesc(id string, desc string) {
	d := getdb()
	if d == nil {
		return
	}
	if stmtFixGameDesc == nil {
		stmt, err := d.Prepare("UPDATE game SET description=? where game_id=?")
		if err == nil {
			stmtFixGameDesc = stmt
		}
	}

	l := len(desc)
	if l == 0 {
		return
	}
	a := strings.Index(desc, "title=")
	if a == -1 {
		return
	}
	b := strings.Index(desc, "\">")
	if b !=-1 && b < l {
		desc_new := strings.Trim(desc[b+2:], " \n")
		fmt.Println("update desc : ", id, desc_new)
		stmtFixGameDesc.Exec(desc_new, id)
	}
}