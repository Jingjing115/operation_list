package main

import (
	"flag"
	"parse_op/srv"
	"log"
	"os"
	"syscall"
)

var configFile string
var signalChan chan os.Signal
var hookableSignals []os.Signal

func init() {
	flag.StringVar(&configFile, "config", "config.yml", "The path of config file.")

	signalChan = make(chan os.Signal, 1)

	hookableSignals = []os.Signal{
		// 用于热更新可执行文件
		syscall.SIGHUP,
		// 终止运行, 可能会丢数据
		//syscall.SIGINT,
		// 终止运行, 会等到所有数据都处理完毕后再退出
		//syscall.SIGTERM,
	}
}

func main() {
	flag.Parse()

	server := srv.NewSrv(configFile)

	err := server.Start()
	if err != nil {
		log.Println(err)
	}

	go server.Process()
	for {
		sig := <-signalChan
		switch sig {
		case syscall.SIGHUP:
		case syscall.SIGINT:
		case syscall.SIGTERM:
		}
	}
}
