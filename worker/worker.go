package worker

import (
	"time"
)

import "github.com/EurasianMagpie/celadon-service/mon"
import "github.com/EurasianMagpie/celadon-service/ipc"
import "github.com/EurasianMagpie/celadon-service/db"


func RunWorker() {

	go idleProc()

	ipc.RunServer()
}

func idleProc() {
	for {
		time.Sleep(time.Duration(5)*time.Minute)

		if mon.IsCacheValid() {
			ipc.GenerateIdleTask()
		} else {
			mon.RunMonTask(false)

			db.ReCheckGameDetail()
		}
	}
}