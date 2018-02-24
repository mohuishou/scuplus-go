package job

import (
	"os"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

// StartServer 初始化server
func StartServer() (*machinery.Server, error) {
	cnf, err := config.NewFromYaml(os.Getenv("SCUPLUS_JOB_CONF"), true)
	if err != nil {
		return nil, err
	}

	// Create server instance
	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}

	// Register tasks
	tasks := map[string]interface{}{}

	return server, server.RegisterTasks(tasks)
}
