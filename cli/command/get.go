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
		task, err := schedulerClient.Get(taskId, false)
		if err != nil {
			return err
		}
		fmt.Printf("Id: %s\n", task.Id)
		fmt.Printf("TaskType: %s\n", task.TaskType)
		fmt.Printf("Time: %s\n", task.Time)
		fmt.Printf("Script: %s\n", task.Script)
		fmt.Println("Results: ")
		for index, result := range task.Results {
			fmt.Printf("[%d]%s\n", index, result.Timestamp)
			fmt.Printf("%s\n", result.Result)
		}
		return nil
	},
}
