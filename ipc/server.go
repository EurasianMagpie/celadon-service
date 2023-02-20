package ipc

import (
	"celadon-service/db"
	"celadon-service/mon"
	"container/list"
	"errors"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"sync"
)

type CcList struct {
	data *list.List
	lock *sync.Mutex
}

func NewCcList() *CcList {
	l := new(CcList)
	l.data = list.New()
	l.lock = new(sync.Mutex)
	return l
}

func (l *CcList) push(v interface{}) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.data.PushBack(v)
}

func (l *CcList) pop() (interface{}, error) {
	l.lock.Lock()
	defer l.lock.Unlock()
	if l.data.Len() > 0 {
		e := l.data.Front()
		l.data.Remove(e)
		return e.Value, nil
	}
	return nil, errors.New("empty data")
}

func (l *CcList) empty() bool {
	l.lock.Lock()
	defer l.lock.Unlock()
	return l.data.Len() == 0
}

func (t *Server) AddTask(arg *TaskArg, reply *int) error {
	if arg != nil {
		*reply = len(arg.Id)
		taskList.push(arg)

	}
	return nil
}

func GenerateIdleTask() {
	if !taskList.empty() {
		return
	}
	arg := NewTaskArg(db.FindAnyUnDetailedGames(2))
	taskList.push(&arg)
}

func taskProc() {
	for {
		if taskList.empty() {
			continue
		}
		e, err := taskList.pop()
		if err == nil {
			var arg *TaskArg = e.(*TaskArg)
			for i := 0; i < len(arg.Id); i++ {
				mon.DeepFetchGame(arg.Id[i])
			}
		}
	}
}

var taskList *CcList

type Server int

func RunServer() {
	db.ReCheckGameDetail()

	taskList = NewCcList()

	go taskProc()

	server := new(Server)
	rpc.Register(server)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":24693")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
