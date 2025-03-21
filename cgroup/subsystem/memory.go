package subsystem

import (
	"fmt"
	"os"
	"strconv"
)

type MemorySubsystem struct {
}

func (s *MemorySubsystem) Name() string {
	return "memory"
}

func (s *MemorySubsystem) Apply(path string, pid int) error {
	fullPath := GetV2ResourcePath(path, V2Procs)

	if err := os.WriteFile(fullPath, []byte(strconv.Itoa(pid)), 0655); err != nil {
		return fmt.Errorf("failed to add process %v to cgroup %v: %v", pid, fullPath, err)
	}
	return nil
}

func (s *MemorySubsystem) Set(path string, res *Resource) error {
	maxMemory := res.MemoryLimit
	fullPath := GetV2ResourcePath(path, V2MemoryMax)

	if err := os.WriteFile(fullPath, []byte(maxMemory), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}

func (s *MemorySubsystem) Remove(path string) error {
	fullPath := GetV2ResourcePath(path, V2MemoryMax)
	nolimit := "max"

	if err := os.WriteFile(fullPath, []byte(nolimit), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}
