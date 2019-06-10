package worker

import (
	"time"
)

import "github.com/EurasianMagpie/celadon/mon"
import "github.com/EurasianMagpie/celadon/ipc"
import "github.com/EurasianMagpie/celadon/db"


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