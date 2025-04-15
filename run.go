package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/markyangcc/mydocker/cgroup"
	"github.com/markyangcc/mydocker/cgroup/subsystem"
	"github.com/markyangcc/mydocker/container"
	"github.com/sirupsen/logrus"
)

func Run(tty bool, command []string, res *subsystem.Resource, volume string) error {
	logrus.Infof("[run] Run %q with tty %v res %v", command, tty, res)

	parent, writePipe, err := container.NewParentProcess(tty, volume)
	if err != nil {
		return fmt.Errorf("failed new parent process:%v", err)
	}

	if err := parent.Start(); err != nil {
		return fmt.Errorf("failed to start:%v", err)
	}

	// 添加 cgroup 资源限制
	logrus.Infof("[run] Apply cgroup limit")
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
	if err := cm.Apply(parent.Process.Pid); err != nil {
		return fmt.Errorf("failed to add pid to cgroup:%v", err)
	}

	// send command to child
	if err := sendInitCommand(command, writePipe); err != nil {
		return fmt.Errorf("failed to send command to child:%v", err)
	}

	// parent pending here
	if err := parent.Wait(); err != nil {
		return fmt.Errorf("failed to wait init:%v", err)
	}

	// umount volume
	if volume != "" {
		if err := container.UnmountVolume("busybox", volume); err != nil {
			return fmt.Errorf("failed to cleanup volume mount:%v", err)
		}
	}
	// do cleanup
	if err := container.CleanupContainerRootfs("busybox"); err != nil {
		return fmt.Errorf("failed to cleanup rootfs:%v", err)
	}

	return nil
}

func sendInitCommand(command []string, pipe *os.File) error {
	cmd := strings.Join(command, " ")
	if _, err := pipe.WriteString(cmd); err != nil {
		return err
	}

	// Close to trigger EOF for reader (unblocks io.ReadAll)
	pipe.Close()
	return nil
}
