package main

import (
	"flag"
	"log"
	"MagicseaServer/Server/config"
	"MagicseaServer/GAServer/app"
	"MagicseaServer/Server/center"
	"MagicseaServer/Server/login"
	"MagicseaServer/Server/db"
	"MagicseaServer/Server/cluster"
	"MagicseaServer/Server/gate"
	"MagicseaServer/Server/session"
	"MagicseaServer/Server/game"
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
	app.RegisterService(session.Type(), session.Service)
	app.RegisterService(login.Type(), login.Service)
	app.RegisterService(gate.Type(), gate.Service)
	app.RegisterService(game.Type(), game.Service)

	log.Println("===Run===", conf)
	app.Run(&conf.Base, cluster.New(), db.NewRedisMgr())
}

