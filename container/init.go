package container

import (
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command []string, args []string) error {
	logrus.Infof("command %s", command)

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	cmd := command[0]
	argv := command
	if err := syscall.Exec(cmd, argv, os.Environ()); err != nil {
		logrus.Error(err.Error())
	}
	return nil
}
