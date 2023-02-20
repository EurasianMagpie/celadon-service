package mon

import (
	"bufio"
	"celadon-service/config"
	"celadon-service/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func monCacheFilePath() (string, error) {
	dir, err := util.GetMonDataDir()
	if err != nil {
		return "", err
	}
	fname := dir + "/page.html"
	return fname, nil
}

func monCacheCfgPath() (string, error) {
	dir, err := util.GetMonDataDir()
	if err != nil {
		return "", err
	}
	fname := dir + "/cfg"
	return fname, nil
}

func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
		merr := os.MkdirAll(dirName, os.ModePerm)
		if merr != nil {
			panic(merr)
		}
	}
}

func FetchHtmlFromUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("FetchPage | read resp failed")
		return "", err
	}
	return string(body), nil
}

func saveFileFromUrl(url string, fname string) bool {
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}

	f, err := os.Create(fname)
	if err != nil {
		return false
	}
	defer f.Close()
	n, err := f.Write(body)
	if err != nil {
		return false
	}
	f.Sync()
	return n > 0
}

func fetchPageNet() (string, error) {
	url := config.GetConfig().Mon.Url
	fmt.Println("FetchPage url:", url)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("FetchPage | read resp failed")
		return "", err
	}

	fname, err := monCacheFilePath()
	if err != nil {
		fmt.Println("FetchPage | get monCacheFilePath error")
		return "", err
	}
	fmt.Println(fname)
	ensureDir(fname)

	f, err := os.Create(fname)
	if err != nil {
		fmt.Println("FetchPage | create file failed")
		return "", err
	}
	defer f.Close()

	n, err := f.Write(body)
	if err != nil {
		fmt.Println("FetchPage | write file failed")
		return "", err
	}
	fmt.Printf("FetchPage | wrote %d bytes\n", n)
	f.Sync()

	cfg, err := monCacheCfgPath()
	if err != nil {
		return "", err
	}
	fcfg, err := os.Create(cfg)
	if err != nil {
		return "", err
	}
	now := time.Now()
	w := bufio.NewWriter(fcfg)
	w.WriteString(now.Format("2006-01-02 15:04:05 MST"))
	w.Flush()

	return string(body), nil
}

func lastFetchTime() time.Time {
	r, _ := time.Parse("2006-01-02", "2000-01-01")

	fname, err := monCacheCfgPath()
	if err != nil {
		return r
	}
	tm, err := ioutil.ReadFile(fname)
	if err != nil {
		return r
	}
	strTime := string(tm)
	strTime = strings.Trim(strTime, " \n")

	layout := "2006-01-02 15:04:05 MST"
	lastTime, err := time.Parse(layout, strTime)
	if err != nil {
		return r
	}
	return lastTime
}

func isFileExist(fname string) bool {
	if _, err := os.Stat(fname); os.IsNotExist(err) {
		return false
	}
	return true
}

func IsCacheValid() bool {
	lastTime := lastFetchTime()
	duration := time.Since(lastTime)
	fmt.Println("IsCacheValid | duration(Hours):", duration.Hours())
	if duration.Hours() < 6 {
		return true
	}
	return false
}

func fetchPageLocal() (string, error) {
	fmt.Println("FetchPage | from local cache")
	fpath, err := monCacheFilePath()
	if err != nil {
		return "", err
	}
	t, err := ioutil.ReadFile(fpath)
	if err != nil {
		return "", err
	}
	return string(t), nil
}

func FetchPage() (string, error) {
	if IsCacheValid() {
		r1, err := fetchPageLocal()
		if err == nil {
			return r1, nil
		}
	}
	return fetchPageNet()
}

func gameCoverDir() (string, error) {
	d, err := util.GetMonDataDir()
	if err != nil {
		return "", err
	}
	dir := d + "/cover"
	return dir, nil
}

func getUrlFileName(_url string) (string, error) {
	u, err := url.Parse(_url)
	if err != nil {
		return "", err
	}
	return filepath.Base(u.Path), nil
}

func FetchGameCoverIfNeeded(id string, url string, _type string) {
	if len(id) == 0 || len(url) == 0 {
		return
	}
	//fmt.Println("FetchGameCoverIfNeeded", id, url)
	ext := _type
	if _type != "webp" {
		ext = "jpg"
	}

	//fmt.Println("FetchGameCoverIfNeeded", ext)
	dir, err := gameCoverDir()
	if err != nil {
		return
	}

	coverImgPath := fmt.Sprintf("%s/%s.%s", dir, id, ext)
	if isFileExist(coverImgPath) {
		//fmt.Println("FetchGameCoverIfNeeded file allready exist")
		return
	}

	ensureDir(coverImgPath)
	saveFileFromUrl(url, coverImgPath)
	//fmt.Println("FetchGameCoverIfNeeded save cover image")
}
