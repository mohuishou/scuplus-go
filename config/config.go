package config

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

// Mysql 配置
type Mysql struct {
	Host     string
	User     string
	Password string
	DB       string
	Port     string
}

// Config 对应config.yml文件的位置
type Config struct {
	Mysql `toml:"mysql"`
}

var config Config

// GetConfig 获取config
func GetConfig(path string) Config {

	if config.Host == "" {
		// 默认配置文件在同级目录
		filepath := getPath(path)

		// 解析配置文件
		if _, err := toml.DecodeFile(filepath, &config); err != nil {
			log.Fatal("配置文件读取失败！", err)
		}
	}
	return config
}

func getPath(path string) string {
	if path != "" {
		return path
	}
	// 获取当前环境
	env := os.Getenv("SCUPLUS_ENV")
	if env == "" {
		env = "develop"
	}

	// 默认配置文件在同级目录
	filepath := "config.toml"

	// 根据环境变量获取配置文件目录
	switch env {
	case "test":
		filepath = os.Getenv("GOPATH") + "/src/github.com/mohuishou/scuplus-go/config/" + filepath
	}
	return filepath
}
