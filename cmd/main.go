package main

import (
	"flag"
	"fmt"
	"go-cache/config"
	"go-cache/server"
	"go-cache/server/handler"
)

func main() {
	fileName := flag.String("--config", "", "path of the config file")
	// 加载配置文件
	if *fileName != "" {
		config.SetupConfig(*fileName)
	}

	fmt.Println(fmt.Sprintf("starting... please dial %s:%d", config.GlobalConfig.Address, config.GlobalConfig.Port))

	server.Start(handler.NewDeal())

}
