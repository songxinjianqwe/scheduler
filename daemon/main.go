package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/songxinjianqwe/scheduler/daemon/server"
	"os"
)

func init() {
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	log.SetFormatter(&log.TextFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	log.SetOutput(os.Stdout)
	//设置最低loglevel
	log.SetLevel(log.InfoLevel)
}

func main() {
	server.Run()
}
