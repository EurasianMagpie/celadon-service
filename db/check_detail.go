package db


var mapCheckGameDetail map[string]bool

func ReCheckGameDetail() bool {
	d := getdb()
	if d == nil {
		return false
	}

	var m map[string]bool
	m = make(map[string]bool)
	rows, err := d.Query("select game_id, description from game")
	if err != nil {
		return false
	}
	defer rows.Close()
	for rows.Next() {
		var id, desc string
		err := rows.Scan(&id, &desc)
		if err != nil {
			return false
		}
		m[id] = len(desc)>0
		//fixGameDesc(id, desc)
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