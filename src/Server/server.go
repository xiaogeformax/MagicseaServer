package main

import (
	"flag"
	"log"
	
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
	log.Println("===Run===", conf)
}

