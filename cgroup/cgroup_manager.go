package cgroup

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/markyangcc/mydocker/cgroup/subsystem"
)

type CgroupManager struct {
	// 当前 cgroup 路径，但不包括RootPath (/sys/fs/cgroup)
	Path string
	// 资源限制
	Resource *subsystem.Resource
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

func (cm *CgroupManager) Apply(pid int) error {
	for _, sys := range subsystem.Subsystems {
		err := sys.Apply(cm.Path, pid)
		if err != nil {
			return fmt.Errorf("failed to apply pid %v to cgroup %v: %v", pid, cm.Path, err)
		}
	}
	return nil
}

func (cm *CgroupManager) Set(res *subsystem.Resource) error {

	for _, sys := range subsystem.Subsystems {
		err := sys.Set(cm.Path, res)
		if err != nil {
			return fmt.Errorf("failed to apply resource %v to cgroup %v: %v", res, cm.Path, err)
		}
	}
	return nil
}

func (cm *CgroupManager) Destory(path string) error {
	fullPath := filepath.Join(subsystem.CgroupV2Root, path)
	if err := os.RemoveAll(fullPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to ddestory cgroup %v: %v", fullPath, err)
	}
	return nil
}
