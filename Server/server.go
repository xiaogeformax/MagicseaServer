package main

import (
	"flag"
	"log"
	"MagicseaServer/Server/config"
	"MagicseaServer/GAServer/app"
	"MagicseaServer/Server/center"
	"MagicseaServer/Server/login"
)

var (
	confPath = flag.String("config", "config.json", "配置文件")
)

func main(){
	flag.Parse()

	conf, err := config.LoadConfig(*confPath)
	if err != nil {
		log.Println("load config err:", err)
		return
	}

	app.RegisterService(center.Type(), center.Service)
	app.RegisterService(login.Type(), login.Service)
	log.Println("===Run===", conf)
}

