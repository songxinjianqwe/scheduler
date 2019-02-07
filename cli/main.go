package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const (
	appName    = "scheduler"
	appVersion = "1.0.0"
)

/**
全局有效
*/
func init() {
	//设置输出样式，自带的只有两种样式logrus.JSONFormatter{}和logrus.TextFormatter{}
	log.SetFormatter(&log.TextFormatter{})
	//设置output,默认为stderr,可以为任何io.Writer，比如文件*os.File
	log.SetOutput(os.Stdout)
	//设置最低loglevel
	log.SetLevel(log.InfoLevel)
}

/**
	主goroutine从命令行读入类型、延迟时间与打印内容，并构造延迟任务或定时任务并使用timer来调度
 */
func main() {
	app := cli.NewApp()
	app.Name = appName
	app.Version = appVersion
	app.Commands = []cli.Command{
		{
			Name:    "submit",
			Aliases: []string{"s"},
			Usage:   "submit a task",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "type",
					Usage: "task type: delay or cron",
				},
				cli.DurationFlag{
					Name:  "hour",
					Usage: "hour",
				},
				cli.DurationFlag{
					Name:  "minute",
					Usage: "minute",
				},
				cli.DurationFlag{
					Name:  "second",
					Usage: "second",
				},
				cli.StringFlag{
					Name:  "scriptPath",
					Usage: "shell script",
				},
			},
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list all tasks",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
		{
			Name: "watch",
			Aliases: []string{"w"},
			Usage: "watch results of a task",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
