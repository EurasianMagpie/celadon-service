package api

import "errors"
//import "net/http"
import "os"
import "fmt"
//import "bufio"

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/config"

/*
func RegisterStaticRoutes(r *gin.Engine) {
	dir, err := gameCoverDir()
	if err == nil {
		staticsubdomain := r.Group("/img")
		staticsubdomain.StaticFS("/cover", http.Dir(dir))
	}
}//*/

func RegisterCoverRoutes(r *gin.Engine) {
	imgSubdomain := r.Group("/celadon")
	imgSubdomain.GET("/cover/:name", getCover)
}

func getCover(c *gin.Context) {
	name := c.Param("name")
	if len(name) == 0 {
		c.JSON(200, formResult(301, string("invalid param ..."), gin.H{}))
		return
	}
	

	d, err := gameCoverLocalDir()
	if err != nil {
		c.JSON(200, formResult(301, string("cover not found 1"), gin.H{}))
		return
	}
	s := "/" + name
	path := d + s

	fmt.Println(path)
	/*f, err := os.Open(path)
	if err != nil {
		c.JSON(200, formResult(301, string("cover not found 2"), gin.H{}))
		return
	}

	reader := bufio.NewReader(f)
	contentLength := int64(reader.Size())
	contentType := "image/webp"
	extraHeaders := map[string]string {
		"x-i": id,
	}
	c.DataFromReader(http.StatusOK, contentLength, contentType, reader, extraHeaders)//*/
	c.File(path)//*/
	return
}

func gameCoverLocalDir() (string, error) {
	d, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := d + config.GetMonDataDir() + "/cover"
	return dir, nil
}

func isFileExist(fname string) bool {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return false
	}
	return true
}

// relative path
func getGameCoverRefPath(id string) (string, error) {
	d, err := gameCoverLocalDir()
	if err != nil {
		return "", err
	}
	s := "/" + id + ".webp"
	p := d + s
	if isFileExist(p) {
		return "/celadon/cover" + s, nil
	}
	return "", errors.New("file not found")
}

