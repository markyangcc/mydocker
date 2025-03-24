package subsystem

import (
	"fmt"
	"os"

	"github.com/markyangcc/mydocker/cgroup"
)

type CpuSubsystem struct {
}

func (s *CpuSubsystem) Name() string {
	return "cpu"
}

func (s *CpuSubsystem) Apply(path string, pid string) error {
	fullPath := cgroup.GetV2ResourcePath(path, cgroup.V2CPUShare)

	if err := os.WriteFile(fullPath, []byte(pid), 0655); err != nil {
		return fmt.Errorf("failed to add process %v to cgroup %v: %v", pid, fullPath, err)
	}
	return nil
}

func (s *CpuSubsystem) Set(path string, res *cgroup.Resource) error {
	fullPath := cgroup.GetV2ResourcePath(path, cgroup.V2CPUShare)
	share := res.CpuShare

	if err := os.WriteFile(fullPath, []byte(share), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}

func (s *CpuSubsystem) Remove(path string) error {
	fullPath := cgroup.GetV2ResourcePath(path, cgroup.V2CPUShare)
	nolimit := "100" // default value 100

	if err := os.WriteFile(fullPath, []byte(nolimit), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}
