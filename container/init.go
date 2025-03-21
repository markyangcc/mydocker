package container

import (
	"fmt"
	"os"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess(command []string, args []string) error {
	logrus.Infof("[init] command %s", command)

	// disable for now
	// defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	cmd := command[0]
	argv := command
	if err := syscall.Exec(cmd, argv, os.Environ()); err != nil {
		return fmt.Errorf("failed to call syscall.exec: %v", err)
	}
	return nil
}
