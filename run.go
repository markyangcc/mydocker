package main

import (
	"github.com/markyangcc/mydocker/container"
	"github.com/sirupsen/logrus"
)

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}
	if err := parent.Wait(); err != nil {
		logrus.Error(err)
	}
}
