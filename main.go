package main

import "fmt"
import "flag"
import "strings"

import "github.com/gin-gonic/gin"

//import "github.com/EurasianMagpie/celadon-service/debug"
import "github.com/EurasianMagpie/celadon-service/api"
import "github.com/EurasianMagpie/celadon-service/mon"
import "github.com/EurasianMagpie/celadon-service/worker"

func main()  {
	//debug.Info()
	t := flag.String("type", "", "process type. api or mon")
	d := flag.Int("deep", 0, "1 > deep fetch game info")
	flag.Parse()

	if strings.Compare(*t, "api") == 0 {
		fmt.Println("main | api")
		r := gin.Default()
		api.RegisterApiRoutes(r)
		api.RegisterCoverRoutes(r)
		api.RegisterDownloadRoutes(r)
		api.RegisterResourceRoutes(r)
		api.PrepareToRun()
		r.Run("localhost:8080")
	} else if strings.Compare(*t, "mon") == 0 {
		fmt.Println("main | mon , deep:", *d==1)
		mon.RunMonTask(*d==1)
	} else if strings.Compare(*t, "worker") == 0 {
		fmt.Println("main | worker")
		worker.RunWorker()
	} else {
		fmt.Println("Please specify process type")
	}
}