package subsystem

import (
	"fmt"
	"os"

	"github.com/markyangcc/mydocker/cgroup"
)

type MemorySubsystem struct {
}

func (s *MemorySubsystem) Name() string {
	return "memory"
}

func (s *MemorySubsystem) Apply(path string, pid string) error {
	fullPath := cgroup.GetV2ResourcePath(path, cgroup.V2MemoryMax)

	if err := os.WriteFile(fullPath, []byte(pid), 0655); err != nil {
		return fmt.Errorf("failed to add process %v to cgroup %v: %v", pid, fullPath, err)
	}
	return nil
}

func (s *MemorySubsystem) Set(path string, res *cgroup.Resource) error {
	maxMemory := res.MemoryLimit
	fullPath := cgroup.GetV2ResourcePath(path, cgroup.V2MemoryMax)

	if err := os.WriteFile(fullPath, []byte(maxMemory), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}

func (s *MemorySubsystem) Remove(path string) error {
	fullPath := cgroup.GetV2ResourcePath(path, cgroup.V2MemoryMax)
	nolimit := "max"

	if err := os.WriteFile(fullPath, []byte(nolimit), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}
