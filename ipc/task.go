package ipc


type TaskArg struct {
	Id []string
}

func NewTaskArg(id []string) TaskArg {
	arg := TaskArg{Id:id}
	return arg
}