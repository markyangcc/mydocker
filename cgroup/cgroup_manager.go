package cgroup

import "path/filepath"

const (
	CgroupV2Root = "/sys/fs/cgroup"

	// cgroup v2 暴露文件接口
	V2MemoryMax = "memory.max"
	V2Procs     = "cgroup.procs"
	V2CPUSets   = "cpuset.cpus"
	V2CPUShare  = "cpu.weight"
)

type Resource struct {
	CpuShare    string
	CpuSet      []string
	MemoryLimit string
}

type Subsystem interface {
	// cgroup 名称
	Name() string
	// 将 pid 加入 cgroup 中
	Apply(path string, pid string) error
	// 设置 cgroup 资源限制
	Set(path string, res *Resource) error
	// 删除 cgroup
	Remove(path string) error
}

// 返回 cgroupv2 资源限制路径
func GetV2ResourcePath(cgroupName string, resourceName string) string {
	return filepath.Join("/sys/fs/cgroup", cgroupName, resourceName)
}

// var Subsystems  = []Subsystem{
// 	&
// }
