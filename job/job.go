package job

import (
	"os"

	"github.com/mohuishou/scuplus-go/job/tasks"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

// StartServer 初始化server
func StartServer() (*machinery.Server, error) {
	// 根据环境变量读取配置文件
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
	t := map[string]interface{}{
		"update_all": tasks.UpdateAll,
	}

	return server, server.RegisterTasks(t)
}
