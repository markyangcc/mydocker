package subsystem

import (
	"fmt"
	"os"
	"strconv"
)

type CpuSetSubsystem struct {
}

func (s *CpuSetSubsystem) Name() string {
	return "cpuset"
}

func (s *CpuSetSubsystem) Apply(path string, pid int) error {
	fullPath := GetV2ResourcePath(path, V2Procs)

	if err := os.WriteFile(fullPath, []byte(strconv.Itoa(pid)), 0655); err != nil {
		return fmt.Errorf("failed to add process %v to cgroup %v: %v", pid, fullPath, err)
	}
	return nil
}

func (s *CpuSetSubsystem) Set(path string, res *Resource) error {
	fullPath := GetV2ResourcePath(path, V2CPUSets)
	cpuSet := res.CpuSet

	if err := os.WriteFile(fullPath, []byte(cpuSet), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}

func (s *CpuSetSubsystem) Remove(path string) error {
	fullPath := GetV2ResourcePath(path, V2CPUSets)
	nolimit := ""

	if err := os.WriteFile(fullPath, []byte(nolimit), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}
