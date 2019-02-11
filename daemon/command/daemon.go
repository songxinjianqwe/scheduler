package command

import (
	"github.com/songxinjianqwe/scheduler/daemon/server"
	"github.com/urfave/cli"
)

var DaemonCommand = cli.Command{
	Name:    "daemon",
	Aliases: []string{"d"},
	Usage:   "run daemon server",
	Action: func(c *cli.Context) error {
		server.Run()
		return nil
	},
}
