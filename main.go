package main

import (
	"fmt"
	"os"

	"github.com/stream1080/godis/config"
	"github.com/stream1080/godis/lib/logger"
	"github.com/stream1080/godis/resp/handler"
	"github.com/stream1080/godis/tcp"
)

const configFileName string = "resp.conf"

// 默认配置
var defaultProperties = &config.ServerProperties{
	Bind: "0.0.0.0",
	Port: 6379,
}

func fileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	return err == nil && !info.IsDir()
}

var banner = `
   ______          ___
  / ____/___  ____/ (_)____
 / / __/ __ \/ __  / / ___/
/ /_/ / /_/ / /_/ / (__  )
\____/\____/\__,_/_/____/


`

func main() {

	fmt.Print(banner)

	// 默认配置
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})

	if fileExists(configFileName) {
		config.SetupConfig(configFileName)
	} else {
		config.Properties = defaultProperties
	}

	// 启动服务
	err := tcp.ListenAndServeWithSignal(&tcp.Config{
		Address: fmt.Sprintf("%s:%d", config.Properties.Bind, config.Properties.Port),
	}, handler.MakeRespHandler())

	if err != nil {
		logger.Error(err)
	}
}
