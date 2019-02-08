package command

import (
	"fmt"
	"github.com/songxinjianqwe/scheduler/cli/client"
	"github.com/urfave/cli"
	"os"
	"text/tabwriter"
)

var ListCommand = cli.Command{
	Name:         "list",
	Aliases:      []string{"l"},
	Usage:        "list all tasks",
	Action: func(c *cli.Context) error {
		schedulerClient, err := client.NewClient()
		if err != nil {
			return err
		}
		tasks, err := schedulerClient.List()
		if err != nil {
			return err
		}
		if len(tasks) == 0 {
			fmt.Println("no tasks")
			return nil
		}
		// 表格式打印
		w := tabwriter.NewWriter(os.Stdout, 15, 1, 3, ' ', 0)
		fmt.Fprint(w, "Id\tTaskType\tTime\tScript\tStatus\n")
		for _, item := range tasks {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				item.Id,
				item.TaskType,
				item.Time,
				item.Script,
				item.Status.String())
		}
		if err := w.Flush(); err != nil {
			return err
		}
		return nil
	},
}