package container

import (
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// Retune init's cmd & writePipe
func NewParentProcess(tty bool) (*exec.Cmd, *os.File, error) {
	readPipe, writePipe, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}

	cmd := exec.Command("/proc/self/exe", "init")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS |
			syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}
	cmd.ExtraFiles = []*os.File{readPipe}
	// Hardcoding for now
	rootPath := "/tmp"
	mntPath := filepath.Join(rootPath, "mnt")
	if err := NewWorkSpace(rootPath); err != nil {
		return nil, nil, err
	}
	cmd.Dir = mntPath

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd, writePipe, nil
}
