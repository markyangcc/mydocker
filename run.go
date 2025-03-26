package main

import (
	"github.com/markyangcc/mydocker/cgroup"
	"github.com/markyangcc/mydocker/cgroup/subsystem"
	"github.com/markyangcc/mydocker/container"
	"github.com/sirupsen/logrus"
)

func Run(tty bool, command []string, res *subsystem.Resource) {
	logrus.Infof("[run] Run %q with tty %v res %v", command, tty, res)
	parent := container.NewParentProcess(tty, command)
	if err := parent.Start(); err != nil {
		logrus.Error(err)
	}

	// 添加 cgroup 资源限制
	cm := cgroup.NewCgroupManager("mydocker-cgroup")
	cm.Set(res)
	cm.Apply(parent.Process.Pid)

	if err := parent.Wait(); err != nil {
		logrus.Error(err)
	}
}
