package db

import "strings"
import "database/sql"
import "fmt"

import "github.com/EurasianMagpie/celadon/util"

var stmtUpdateGameDesc *sql.Stmt

func init() {
	d := getdb()
	if d == nil {
		return
	}
	if stmtUpdateGameDesc == nil {
		stmt, err := d.Prepare("UPDATE game SET description=? where game_id=?")
		if err == nil {
			stmtUpdateGameDesc = stmt
		}
	}
}

func fixGameDesc(id string, desc string) {
	d := getdb()
	if d == nil {
		return
	}
	if stmtUpdateGameDesc == nil {
		return
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
		stmtUpdateGameDesc.Exec(desc_new, id)
	}
}

func fixUnEscapedDesc(id string, desc string) {
	bDesc, rDesc := util.UnEscape(desc)
	if bDesc {
		stmtUpdateGameDesc.Exec(rDesc, id)
		//fmt.Println(id, rDesc)
	}
}