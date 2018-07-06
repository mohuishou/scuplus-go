package main

import (
	"log"

	"github.com/mohuishou/scuplus-go/config"
	"github.com/mohuishou/scuplus-go/job"
)

func main() {
	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	log.Println("最大并行数目：", config.Get().JobWorkers)
	worker := job.Server.NewWorker("machinery_worker", config.Get().JobWorkers)

	worker.Launch()
}
