package mon

import "github.com/EurasianMagpie/celadon/db"


func UpdateResult(result *ParseResult) {
	// fetch gameinfo.cover image

	// update table region, game, price
	// caution: symbols to escape in sql
	for _, region := range result.Regions {
		db.UpdateRegion(region)
	}

	for _, game := range result.Games {
		db.UpdateGame(game)

		FetchGameCoverIfNeeded(game.Id, game.CoverUrl, game.CoverType)
	}//*/

	for _, price := range result.Prices {
		db.UpdatePrice(price)
	}//*/
}