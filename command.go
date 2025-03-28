package main

import (
	"fmt"

	"github.com/markyangcc/mydocker/cgroup/subsystem"
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
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit",
		},
	},
	Action: func(context *cli.Context) error {
		if len(context.Args()) < 1 {
			return fmt.Errorf("missing container command")
		}
		// entryCmd := context.Args().Get(0)
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}
		tty := context.Bool("ti")

		resConf := &subsystem.Resource{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}

		Run(tty, cmdArray, resConf)
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {

		// cmd := context.Args().Get(0)
		var cmdArray []string
		for _, arg := range context.Args() {
			cmdArray = append(cmdArray, arg)
		}

		logrus.Infof("[init] user command: %s", cmdArray)
		err := container.RunContainerInitProcess(cmdArray, nil)
		return err
	},
}
