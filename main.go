package main

import (
	"fmt"

	"DStorage/config/server"
	"DStorage/router"
)

func main() {
	// gin framework
	router := router.Router()

	// 启动服务并监听端口
	err := router.Run(server.UploadServiceHost)
	if err != nil {
		fmt.Printf("Failed to start server, err:%s\n", err.Error())
	}
}
