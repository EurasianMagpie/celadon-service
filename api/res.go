package api

import (
	"fmt"
)

import "github.com/gin-gonic/gin"
import "github.com/EurasianMagpie/celadon-service/util"


func RegisterResourceRoutes(r *gin.Engine) {
	celadonSubdomain := r.Group("/celadon")
	celadonSubdomain.GET("/res/:what", fetchResource)
}

func fetchResource(c *gin.Context) {
	what := c.Param("what")
	if len(what) == 0 {
		c.JSON(404, formResult(301, string("invalid param ..."), gin.H{}))
		return
	}

	dir, err := util.GetResDir()
	if err != nil {
		c.JSON(404, formResult(301, string("file not found - Y"), gin.H{}))
		return
	}
	path := dir + "/" + what
	fmt.Println("[Res] FilePath:", path)
	if !util.IsFileExist(path) {
		c.JSON(404, formResult(301, string("file not found - Y"), gin.H{}))
		return
	}
	c.File(path)
}
