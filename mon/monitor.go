package mon

import "fmt"

import "github.com/EurasianMagpie/celadon/db"

func RunMonTask(deep bool) {
	d, err := FetchPage()
	if err != nil {
		return
	}
	result, err := Parse(d, deep)
	if err != nil {
		return
	}
	UpdateResult(result, deep)
}

func DeepFetchGame(id string) {
	// 1. check is not detailed
	// 2. fetch detail page and parse
	// 3. fetch cover img
	// 4. update to db game table
	// 5. mark is detailed
	fmt.Println("DeepFetchGame > ", id)

	if db.IsGameDetialed(id) {
		return
	}
	ginfo, err := db.QueryGameInfo(id)
	if err != nil {
		return
	}

	ok := DeepParseSingleGame(ginfo)
	if !ok {
		fmt.Println("DeepFetchGame > ", id, " DeepParseSingleGame failed")
		return
	}

	ok = db.UpdateGameFull(*ginfo)
	if !ok {
		fmt.Println("DeepFetchGame > ", id, " db.UpdateGame failed")
	}
	FetchGameCoverIfNeeded(ginfo.Id, ginfo.CoverUrl, ginfo.CoverType)

	if ok {
		db.MarkGameDetailed(id)
	}
}