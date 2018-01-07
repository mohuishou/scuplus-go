package main

import (
	"flag"

	"github.com/mohuishou/scuplus-go/task/course"

	"github.com/mohuishou/scuplus-go/config"
)

func main() {
	confPath := flag.String("c", "", "配置文件目录")
	t := flag.String("task", "", "需要执行的任务, 现有任务: course:update")
	flag.Parse()

	// 获取配置文件
	conf := config.GetConfig(*confPath)

	// 执行任务
	switch *t {
	case "course:update":
		course.Task(conf)
	}
}
