package db


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