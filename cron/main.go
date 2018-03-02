package main

import (
	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/mohuishou/scuplus-go/job"
)

func main() {
	sign := &tasks.Signature{
		Name: "update_all",
		Args: []tasks.Arg{
			{
				Type:  "uint",
				Value: 1,
			},
		},
	}
	server, err := job.StartServer()
	if err != nil {
		panic(err)
	}
	_, err = server.SendTask(sign)
	if err != nil {
		panic(err)
	}
}
