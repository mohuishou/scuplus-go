package main

import (
	"github.com/mohuishou/scuplus-go/job"
)

func main() {
	server, err := job.StartServer()
	if err != nil {
		panic(err)
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker("machinery_worker", 0)

	worker.Launch()
}
