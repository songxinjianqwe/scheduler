package main

import (
	"fmt"
	"github.com/songxinjianqwe/scheduler/cli/command"
	"github.com/urfave/cli"
	"os"
)

const (
	appName    = "scheduler"
	appVersion = "1.0.0"
)

/**
	主goroutine从命令行读入类型、延迟时间与打印内容，并构造延迟任务或定时任务并使用timer来调度
 */
func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = appVersion
	app.Commands = []cli.Command{
		command.GetCommand,
		command.ListCommand,
		command.SubmitCommand,
		command.StopCommand,
		command.DeleteCommand,
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}
