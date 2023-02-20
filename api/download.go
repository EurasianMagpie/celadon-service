package api

import (
	"celadon-service/util"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

var regVersionFilter *regexp.Regexp

func init() {
	regVersionFilter, _ = regexp.Compile(`([a-zA-Z0-9]+\-){3}([0-9]+).([0-9]+).([0-9]+).([0-9]+)[\w.]+`)
}

func RegisterDownloadRoutes(r *gin.Engine) {
	celadonSubdomain := r.Group("/celadon")
	celadonSubdomain.GET("/download/:what", downloadLatestApk)
}

func downloadLatestApk(c *gin.Context) {
	what := c.Param("what")
	if len(what) == 0 {
		c.JSON(404, formResult(301, string("invalid param ..."), gin.H{}))
		return
	}

	if what == "latest.apk" {
		fname, err := findLastestApkFileName()
		if err != nil {
			c.JSON(404, formResult(301, string("file not found"), gin.H{}))
			return
		}

		//fmt.Println("[Debug] Request:", c.Request)
		r := "http://" + c.Request.Host + "/celadon/download/" + fname
		c.Redirect(301, r)
	} else {
		dir, err := getDownloadApkDir()
		if err != nil {
			c.JSON(404, formResult(301, string("file not found - Y"), gin.H{}))
			return
		}
		path := dir + "/" + what
		fmt.Println("[Download] FilePath:", path)
		if !util.IsFileExist(path) {
			c.JSON(404, formResult(301, string("file not found - Y"), gin.H{}))
			return
		}
		c.File(path)
	}
}

func getDownloadApkDir() (string, error) {
	d, err := util.GetDownloadDir()
	if err != nil {
		return "", err
	}
	dir := d + "/apk"
	return dir, nil
}

type version struct {
	V1 int
	V2 int
	V3 int
	V4 int
}

func newVersion(v1 string, v2 string, v3 string, v4 string) version {
	x1, _ := strconv.Atoi(v1)
	x2, _ := strconv.Atoi(v2)
	x3, _ := strconv.Atoi(v3)
	x4, _ := strconv.Atoi(v4)
	ver := version{V1: x1, V2: x2, V3: x3, V4: x4}
	return ver
}

func compareVersion(left version, right version) int {
	if left.V1 > right.V1 {
		return left.V1 - right.V1
	}
	if left.V2 > right.V2 {
		return left.V2 - right.V2
	}
	if left.V3 > right.V3 {
		return left.V3 - right.V3
	}
	return left.V4 - right.V4
}

func findLastestApkFileName() (string, error) {
	dir, err := getDownloadApkDir()
	if err != nil {
		return "", err
	}

	// app-naiwen17-release-1.1.0.1009_2019-07-22_aligned_signed.apk
	if regVersionFilter == nil {
		return "", nil
	}

	var latestVersion = version{V1: 0, V2: 0, V3: 0, V4: 0}
	var latestApkFileName = ""
	filepath.Walk(dir,
		func(path string, info os.FileInfo, e error) error {
			if e != nil {
				return e
			}

			// check if it is a regular file (not dir)
			if info.Mode().IsRegular() {
				//fmt.Println("file name:", info.Name())
				//fmt.Println("file path:", path)

				params := regVersionFilter.FindStringSubmatch(info.Name())
				//fmt.Println("version:", params[2], params[3], params[4], params[5])
				nextVersion := newVersion(params[2], params[3], params[4], params[5])
				if compareVersion(nextVersion, latestVersion) > 0 {
					latestVersion = nextVersion
					latestApkFileName = info.Name()
				}
			}
			return nil
		})

	fmt.Println("[Latest Apk] Version:", latestVersion, "File:", latestApkFileName)
	return latestApkFileName, nil
}
