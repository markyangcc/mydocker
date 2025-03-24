package subsystem

import (
	"fmt"
	"os"

	"github.com/markyangcc/mydocker/cgroup"
)

type CpuSetSubsystem struct {
}

func (s *CpuSetSubsystem) Name() string {
	return "cpuset"
}

func (s *CpuSetSubsystem) Apply(path string, pid string) error {
	fullPath := cgroup.GetV2ResourcePath(path, cgroup.V2MemoryMax)

	if err := os.WriteFile(fullPath, []byte(pid), 0655); err != nil {
		return fmt.Errorf("failed to add process %v to cgroup %v: %v", pid, fullPath, err)
	}
	return nil
}

func (s *CpuSetSubsystem) Set(path string, res *cgroup.Resource) error {
	fullPath := cgroup.GetV2ResourcePath(path, cgroup.V2MemoryMax)
	cpuSet := res.CpuSet

	if err := os.WriteFile(fullPath, []byte(cpuSet), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}

func (s *CpuSetSubsystem) Remove(path string) error {
	fullPath := cgroup.GetV2ResourcePath(path, cgroup.V2CPUSets)
	nolimit := ""

	if err := os.WriteFile(fullPath, []byte(nolimit), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}
