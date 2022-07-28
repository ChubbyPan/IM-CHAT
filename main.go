package main

import (
	"main.go/conf"
	"main.go/router"
)

func main() {
	conf.Init()
	r := router.NewRouter()
	_ = r.Run(conf.HttpPort)
}
