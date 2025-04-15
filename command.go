package main

import (
	"fmt"

	"github.com/markyangcc/mydocker/cgroup/subsystem"
	"github.com/markyangcc/mydocker/container"
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
			Name:  "v",
			Usage: "volume",
		},
		cli.StringFlag{
			Name:  "m",
			Usage: "memory limit. For example, -m=100m to limit memory usage",
		},
		cli.StringFlag{
			Name:  "cpushare",
			Usage: "cpushare limit, 100000 means 100ms/100%. For example, -cpushare 80000 to limit cpu to 80%",
		},
		cli.StringFlag{
			Name:  "cpuset",
			Usage: "cpuset limit,The CPU numbers are comma-separated numbers or ranges. For example:, -cpuset 0-4,6,8-10",
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
		volume := context.String("v")

		resConf := &subsystem.Resource{
			MemoryLimit: context.String("m"),
			CpuSet:      context.String("cpuset"),
			CpuShare:    context.String("cpushare"),
		}

		if err := Run(tty, cmdArray, resConf, volume); err != nil {
			return err
		}
		return nil
	},
}

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Init container process run user's process in container. Do not call it outside",
	Action: func(context *cli.Context) error {

		err := container.RunContainerInitProcess()
		return err
	},
}
