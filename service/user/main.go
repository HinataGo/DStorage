package main

import (
	"log"
	"time"

	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"

	// k8s "github.com/micro/kubernetes/go/micro"

	"DStorage/common"
	dbproxy "DStorage/service/dbproxy/client"
	"DStorage/service/user/handler"
	proto "DStorage/service/user/proto"
)

func main() {
	service := micro.NewService(
		// service := k8s.NewService(
		micro.Name("go.micro.service.user"),
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)

	// 初始化service, 解析命令行参数等
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	proto.RegisterUserServiceHandler(service.Server(), new(handler.User))
	if err := service.Run(); err != nil {
		log.Println(err)
	}
}
