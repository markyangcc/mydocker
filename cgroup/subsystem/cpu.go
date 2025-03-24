package subsystem

import (
	"fmt"
	"os"
	"strconv"
)

type CpuSubsystem struct {
}

func (s *CpuSubsystem) Name() string {
	return "cpu"
}

func (s *CpuSubsystem) Apply(path string, pid int) error {
	fullPath := GetV2ResourcePath(path, V2CPUShare)

	if err := os.WriteFile(fullPath, []byte(strconv.Itoa(pid)), 0655); err != nil {
		return fmt.Errorf("failed to add process %v to cgroup %v: %v", pid, fullPath, err)
	}
	return nil
}

func (s *CpuSubsystem) Set(path string, res *Resource) error {
	fullPath := GetV2ResourcePath(path, V2CPUShare)
	share := res.CpuShare

	if err := os.WriteFile(fullPath, []byte(share), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}

func (s *CpuSubsystem) Remove(path string) error {
	fullPath := GetV2ResourcePath(path, V2CPUShare)
	nolimit := "100" // default value 100

	if err := os.WriteFile(fullPath, []byte(nolimit), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}
