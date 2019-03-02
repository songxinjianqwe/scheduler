package command

import (
	"errors"
	"github.com/songxinjianqwe/scheduler/cli/client"
	"github.com/songxinjianqwe/scheduler/common"
	"github.com/urfave/cli"
)

var SubmitCommand = cli.Command{
	Name:    "submit",
	Aliases: []string{"s"},
	Usage:   "submit a task",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "type",
			Usage: "task type: delay or cron",
		},
		cli.DurationFlag{
			Name:  "time",
			Usage: "time",
		},
		cli.StringFlag{
			Name:  "script",
			Usage: "shell script",
		},
	},
	Action: func(c *cli.Context) error {
		// 第一个参数默认为id，检查id是否存在，且非空
		if c.NArg() == 0 {
			return errors.New("must specify task id")
		}
		taskId := c.Args().Get(0)
		if taskId == "" {
			return errors.New("task id can not be blank")
		}
		// 获取任务类型，值必须为delay或cron
		taskType := c.String("type")
		if taskType != "delay" && taskType != "cron" {
			return errors.New("must specify task type with -type=delay/cron")
		}
		// 获取任务执行时间
		time := c.Duration("time")
		// 获取脚本内容
		script := c.String("script")
		if script == "" {
			return errors.New("script can not be blank")
		}
		// 启动客户端
		schedulerClient, err := client.NewClient()
		if err != nil {
			return err
		}
		task := common.NewTask(taskId, taskType, time, script)
		// 提交任务
		err = schedulerClient.Submit(task)
		return err
	},
}
