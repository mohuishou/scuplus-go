package job

import (
	"os"

	"github.com/mohuishou/scuplus-go/job/tasks"

	"github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
)

var Server *machinery.Server

// 初始化server
func init() {
	var err error
	// 根据环境变量读取配置文件
	cnf, err := config.NewFromYaml(os.Getenv("SCUPLUS_JOB_CONF"), true)
	if err != nil {
		panic(err)
	}

	// Create server instance
	Server, err = machinery.NewServer(cnf)
	if err != nil {
		panic(err)
	}

	// Register tasks
	t := map[string]interface{}{
		"update_all":      tasks.UpdateAll,
		"notify_grade":    tasks.NotifyGrade,
		"notify_book":     tasks.NotifyBook,
		"notify_exam":     tasks.NotifyExam,
		"notify_feedback": tasks.NotifyFeedback,
	}

	if err = Server.RegisterTasks(t); err != nil {
		panic(err)
	}
}
