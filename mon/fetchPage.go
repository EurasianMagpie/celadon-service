package mon

import "fmt"
import "io/ioutil"
import "os"
import "path/filepath"

import "net/http"

import "github.com/EurasianMagpie/celadon/config"





func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, serr := os.Stat(dirName); serr != nil {
	  merr := os.MkdirAll(dirName, os.ModePerm)
	  if merr != nil {
		  panic(merr)
	  }
	}
}

func FetchPage() {
	url := config.GetConfig().Mon.Url
	fmt.Println("FetchPage url:", url)
	resp, err := http.Get(url)
	if err != nil {

	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("FetchPage | read resp failed")
		return
	}

	dir, err := os.Getwd()
	fname := dir + "/mondata/page.html"
	fmt.Println(fname)
	ensureDir(fname)
	
	f, err := os.Create(fname)
	if err != nil {
		fmt.Println("FetchPage | create file failed")
		return
	}
	defer f.Close()

	n, err := f.Write(body)
	if err != nil {
		fmt.Println("FetchPage | write file failed")
		return
	}
	fmt.Printf("FetchPage | wrote %d bytes\n", n)
	f.Sync()

}