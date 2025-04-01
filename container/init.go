package container

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

func RunContainerInitProcess() error {

	// disable for now
	// defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	// syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	command, err := readCommandFromPipe()
	if err != nil {
		return fmt.Errorf("reading command from pipe failed: %w", err)
	}

	if len(command) == 0 {
		return errors.New("received empty command")
	}

	execPath, err := exec.LookPath(command[0])
	if err != nil {
		return fmt.Errorf("executable lookup failed: %w", err)
	}
	if err != nil {
		return fmt.Errorf("resolving executable path failed: %w", err)
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
