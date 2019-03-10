package ipc

import (
	"net/rpc"
	"fmt"

)

var client *rpc.Client

func getClient() (*rpc.Client, error) {
	if client != nil {
		return client, nil
	}
	c, err := rpc.DialHTTP("tcp", "127.0.0.1:24693")
	if err != nil {
		return nil, err
	}
	client = c
	return client, nil
}

func invalidClient() {
	if client != nil {
		client = nil
	}
}

func AddTask(id []string) (error) {
	c, err := getClient()
	if err != nil {
		return err
	}

	arg := NewTaskArg(id)
	var reply int
	addTaskCall := c.Go("Worker.AddTask", arg, &reply, nil)
	replyCall := <-addTaskCall.Done
	if replyCall.Error != nil {
		fmt.Println("Worker error:", replyCall.Error)
		defer c.Close()
		invalidClient()
		return replyCall.Error
	}
	return nil
}