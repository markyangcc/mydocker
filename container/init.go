package container

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {

	command, err := readCommandFromPipe()
	if err != nil {
		return fmt.Errorf("reading command from pipe failed: %w", err)
	}

	if len(command) == 0 {
		return errors.New("received empty command")
	}

	err = setupMount()
	if err != nil {
		return fmt.Errorf("setup mount failed: %w", err)
	}

	execPath, err := exec.LookPath(command[0])
	if err != nil {
		return fmt.Errorf("executable lookup failed: %w", err)
	}

	logrus.WithFields(logrus.Fields{
		"pid":  os.Getpid(),
		"path": execPath,
		"args": command,
	}).Info("Attempting to execute container command")

	argv := append([]string{execPath}, command[1:]...)
	env := os.Environ()

	if err := syscall.Exec(execPath, argv, env); err != nil {
		return fmt.Errorf("syscall.Exec failed [path=%s, args=%v]: %w",
			execPath, argv, err)
	}

	// Normally never reach here
	return nil
}

func readCommandFromPipe() ([]string, error) {
	const pipeFD = 3

	pipe := os.NewFile(pipeFD, "readpipe")
	defer pipe.Close()

	data, err := io.ReadAll(pipe)
	if err != nil {
		return nil, fmt.Errorf("pipe read failed: %w", err)
	}

	if len(data) == 0 {
		return nil, errors.New("empty data from pipe")
	}

	return strings.Fields(string(data)), nil
}

func pivotRoot(root string) error {
	if err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("mounting failed: %w", err)
	}
	pivotDir := filepath.Join(root, ".pivot_root")
	if err := os.MkdirAll(pivotDir, 0755); err != nil {
		return fmt.Errorf("creating pivot dir failed: %w", err)
	}

	if err := syscall.PivotRoot(root, pivotDir); err != nil {
		return fmt.Errorf("pivot root failed: %w", err)
	}

	if err := syscall.Chdir("/"); err != nil {
		return fmt.Errorf("chdir to / failed: %w", err)
	}

	// umount the old root
	pivotDir = filepath.Join("/", ".pivot_root")
	if err := syscall.Unmount(pivotDir, syscall.MNT_DETACH); err != nil {
		return fmt.Errorf("unmounting pivot dir failed: %w", err)
	}

	if err := os.Remove(pivotDir); err != nil {
		return fmt.Errorf("removing pivot dir failed: %w", err)
	}

	return nil
}

func setupMount() error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %v", err)

	}
	logrus.Infof("current working directory: %s", pwd)

	// 设置挂载传播，使用 private，这样我们当前 mountNS 与 host mountNS脱钩，两边挂载不会互相传播
	// 注意：在子进程任何 mount 操作之前调用
	if err := syscall.Mount("", "/", "", syscall.MS_PRIVATE|syscall.MS_REC, ""); err != nil {
		return fmt.Errorf("setting mount propagation to private failed: %v", err)
	}

	pivotRoot(pwd)

	// mount porc/dev
	if err := os.MkdirAll("/proc", 0755); err != nil {
		return fmt.Errorf("creating pivot dir failed: %w", err)
	}
	if err := os.MkdirAll("/dev", 0755); err != nil {
		return fmt.Errorf("creating pivot dir failed: %w", err)
	}
	err = syscall.Mount("proc", "/proc", "proc", syscall.MS_NOEXEC|syscall.MS_NOSUID|syscall.MS_NODEV, "")
	if err != nil {
		return fmt.Errorf("failed to mount /proc: %v", err)
	}
	err = syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")
	if err != nil {
		return fmt.Errorf("failed to mount /dev: %v", err)
	}

	return nil
}
