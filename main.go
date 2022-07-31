package main

import (
	"main.go/conf"
	"main.go/router"
	"main.go/service"
)

func main() {
	conf.Init()
	go service.Manager.Start()
	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)
}
