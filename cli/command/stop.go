package command

import (
	"errors"
	"github.com/songxinjianqwe/scheduler/cli/client"
	"github.com/urfave/cli"
)

var StopCommand = cli.Command{
	Name: "stop",
	Aliases: []string{"s"},
	Usage: "stop a task",
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
		err = schedulerClient.Stop(taskId)
		return err
	},
}