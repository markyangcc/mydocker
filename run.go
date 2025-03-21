package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/markyangcc/mydocker/cgroup"
	"github.com/markyangcc/mydocker/cgroup/subsystem"
	"github.com/markyangcc/mydocker/container"
	"github.com/sirupsen/logrus"
)

func Run(tty bool, command []string, res *subsystem.Resource) error {
	logrus.Infof("[run] Run %q with tty %v res %v", command, tty, res)
	init := container.NewInitProcess(tty, command)
	if err := init.Start(); err != nil {
		return fmt.Errorf("failed to start:%v", err)
	}

	logrus.Infof("[run] cgroup limit")

	// 添加 cgroup 资源限制
	cm := cgroup.NewCgroupManager("mydocker-cgroup")
	cgPath := filepath.Join("/sys/fs/cgroup", "mydocker-cgroup")
	if _, err := os.Stat(cgPath); err != nil && os.IsNotExist(err) {
		if err1 := os.Mkdir(cgPath, 0644); err1 != nil {
			return fmt.Errorf("failed to create cgroup:%v", err)
		}
	}

	if err := cm.Set(res); err != nil {
		return fmt.Errorf("failed to apply resource limit:%v", err)
	}
	if err := cm.Apply(init.Process.Pid); err != nil {
		return fmt.Errorf("failed to add pid to cgroup:%v", err)
	}

	if err := init.Wait(); err != nil {
		return fmt.Errorf("failed to wait init:%v", err)
	}

	return nil
}
