package main

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const usage = `mydocker is a simple container runtime implementation.
			   The purpose of this project is to learn how docker works and how to write a docker by ourselves
			   Enjoy it, just for fun.`

func main() {
	app := cli.NewApp()
	app.Name = "mydocker"
	app.Usage = usage

	app.Commands = []cli.Command{
		initCommand,
		runCommand,
	}

	app.Before = func(context *cli.Context) error {
		logrus.SetOutput(os.Stdout)
		return nil
	}

	// unmount /proc in RunContainerInitProcess()
	app.After = func(context *cli.Context) error {
		var err error
		if err = syscall.Unmount("/proc", syscall.MNT_DETACH); err != nil {
			logrus.Errorf("failed to umount /proc: %v", err)
		}
		return err
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}

}
