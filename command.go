package main

import (
	"fmt"

	"github.com/markyangcc/mydocker/container"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var runCommand = cli.Command{
	Name: "run",
	Usage: `Create a container with namespace and cgroups limit
			mydocker run -ti [command]`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "ti",
			Usage: "enable tty",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		entryCmd := context.Args().Get(0)
		tty := context.Bool("ti")
		Run(tty, entryCmd)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {
		cmd := context.Args().Get(0)
		logrus.Infof("init come on ,command %s", cmd)
		err := container.RunContainerInitProcess(cmd, nil)
		return err
	},
}
