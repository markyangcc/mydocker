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
	fullPath := GetV2ResourcePath(path, V2Procs)

	if err := os.WriteFile(fullPath, []byte(strconv.Itoa(pid)), 0655); err != nil {
		return fmt.Errorf("failed to add process %v to cgroup %v: %v", pid, fullPath, err)
	}
	return nil
}

func (s *CpuSubsystem) Set(path string, res *Resource) error {
	fullPath := GetV2ResourcePath(path, V2CPUMax)
	share := res.CpuShare
	period := 100000 // https://docs.kernel.org/admin-guide/cgroup-v2.html

	if err := os.WriteFile(fullPath, []byte(share), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}

	if err := os.WriteFile(fullPath,
		[]byte(fmt.Sprintf("%s %d", share, period)), 0644); err != nil {
	}
	return nil
}

func (s *CpuSubsystem) Remove(path string) error {
	fullPath := GetV2ResourcePath(path, V2CPUMax)
	nolimit := "max 100000"

	if err := os.WriteFile(fullPath, []byte(nolimit), 0655); err != nil {
		return fmt.Errorf("failed to set memory to cgroup %v: %v", fullPath, err)
	}
	return nil
}
