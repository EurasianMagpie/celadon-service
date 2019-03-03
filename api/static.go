package api

import "errors"
import "net/http"
import "os"

import "github.com/gin-gonic/gin"

import "github.com/EurasianMagpie/celadon/config"


func RegisterStaticRoutes(r *gin.Engine) {
	dir, err := gameCoverDir()
	if err == nil {
		staticsubdomain := r.Group("/img")
		staticsubdomain.StaticFS("/cover", http.Dir(dir))
	}
}

func gameCoverDir() (string, error) {
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
func getGameCoverFilePath(id string) (string, error) {
	d, err := gameCoverDir()
	if err != nil {
		return "", err
	}
	s := "/" + id + ".webp"
	p := d + s
	if isFileExist(p) {
		return "/img/cover" + s, nil
	}
	return "", errors.New("file not found")
}

