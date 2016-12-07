package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli"
)

var execCommand = cli.Command{
	Name:  "exec",
	Usage: "exec command in a container in a pod",
	ArgsUsage: `<container-id> <container command> [command options]
	
Where "<container-id>" is the id for the instance of the container and
"<container command>" is the command to be executed in the container.`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "tty, t",
			Usage: "exec a command in a TTY",
		},
		cli.BoolFlag{
			Name:  "stdin, i",
			Usage: "stream stdin while exec a command",
		},
	},
	Action: func(context *cli.Context) error {
		if os.Geteuid() != 0 {
			return fmt.Errorf("kpod should be run as root")
		}
		status, err := execCommand(context)
		if err != nil {
			return fmt.Errorf("exec command failed: %v", err)
		}
		os.Exit(status)
	},
}

func execCommand(context *cli.Context) (int, error) {
	container, err := getContainer(context)
	if err != nil {
		return -1, err
	}
}
