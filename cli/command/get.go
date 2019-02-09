package command

import (
	"errors"
	"fmt"
	"github.com/songxinjianqwe/scheduler/cli/client"
	"github.com/urfave/cli"
)

var GetCommand = cli.Command{
	Name: "get",
	Aliases: []string{"g"},
	Usage: "get results of a task",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "watch",
			Usage: "watch task results",
		},
	},
	Action: func(c *cli.Context) error {
		// 第一个参数默认为id，检查id是否存在，且非空
		if c.NArg() == 0 {
			return errors.New("must specify task id")
		}
		taskId := c.Args().Get(0)
		schedulerClient, err := client.NewClient()
		if err != nil {
			return err
		}
		// 找不到则默认为false
		watch := c.Bool("watch")
		// 如果watch为true，则执行一个http long polling
		// 1. cli发起一个HTTP请求，watch=false，获取初始数据，带回版本v1
		// 2. 服务器返回最新数据
		// 3. cli再次发起一个HTTP请求，watch=true，v=v1
		// 4. 判断v是否为v1：
		//    - 如果不是，说明服务器数据有变更，则返回最新数据，v=v2；
		// 	  - 如果是，说明服务器数据没有变更，则在服务器注册一个Listener至该任务对象中，并阻塞；其他goroutine如果修改了任务对象，则v=v2，并唤醒该Listener。
		//	  Listener被唤醒后返回最新数据，v=v2
		// 5. 重复3~4

		task, err := schedulerClient.Get(taskId, false, 0)
		if err != nil {
			return err
		}
		task.PrintMe()
		// 如果是需要监听，则进入一个无线循环
		if watch {
			for {
				fmt.Println("------------------------------------")
				version := task.Version
				task, err = schedulerClient.Get(taskId, watch, version)
				if err != nil {
					return err
				}
				task.PrintMe()
			}
		}
		return nil
	},
}
