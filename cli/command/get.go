package command

import (
	"errors"
	"github.com/songxinjianqwe/scheduler/cli/client"
	"github.com/urfave/cli"
)

var GetCommand = cli.Command{
	Name: "get",
	Aliases: []string{"g"},
	Usage: "get results of a task",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "w",
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
		watch := c.Bool("w")
		// 如果watch为true，则执行一个http long polling
		// 1. cli发起一个HTTP请求，watch=false，获取初始数据，带回版本v1
		// 2. 服务器返回最新数据
		// 3. cli再次发起一个HTTP请求，watch=true，v=v1
		// 4. 判断v是否为v1：
		//    - 如果不是，说明服务器数据有变更，则返回最新数据，v=v2；
		// 	  - 如果是，说明服务器数据没有变更，则在服务器注册一个Listener至该任务对象中，并阻塞；其他goroutine如果修改了任务对象，则v=v2，并唤醒该Listener。
		//	  Listener被唤醒后返回最新数据，v=v2
		// 5. 重复3~4
		// 如果2~3之间有数据变化，那么如何感知到呢？版本概念！
		task, err := schedulerClient.Get(taskId, watch)
		if err != nil {
			return err
		}
		task.PrintMe()
		return nil
	},
}
