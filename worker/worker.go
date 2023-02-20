package worker

import (
	"celadon-service/db"
	"celadon-service/ipc"
	"celadon-service/mon"
	"time"
)

func RunWorker() {

	go idleProc()

	ipc.RunServer()
}

func idleProc() {
	for {
		time.Sleep(time.Duration(5) * time.Minute)

		if mon.IsCacheValid() {
			ipc.GenerateIdleTask()
		} else {
			mon.RunMonTask(false)

			db.ReCheckGameDetail()
		}
	}
}
