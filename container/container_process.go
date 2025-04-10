package container

import (
	"fmt"
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
	contaienrName := "busybox"
	workspace := filepath.Join(OverlayfsRoot, contaienrName)
	mergeddir := filepath.Join(OverlayfsRoot, contaienrName, "merged")
	if err := NewContainerRootfs(contaienrName); err != nil {
		if err1 := os.RemoveAll(workspace); err1 != nil {
			return nil, nil, fmt.Errorf("failed to cleanup workspace:%v:%v", err, err1)
		}
		return nil, nil, err
	}
	cmd.Dir = mergeddir

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd, writePipe, nil
}
