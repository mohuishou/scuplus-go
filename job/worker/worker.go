package main

import (
	"github.com/mohuishou/scuplus-go/job"
)

func main() {
	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := job.Server.NewWorker("machinery_worker", 0)

	worker.Launch()
}
