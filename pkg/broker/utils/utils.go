package utils

import (
	"os"

	"github.com/shirou/gopsutil/v3/process"
)

func Utilization() (float64, uint64) {
	pid := os.Getpid()
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return 0, 0
	}
	cpuPercent, _ := p.CPUPercent()
	memInfo, _ := p.MemoryInfo()
	return cpuPercent, memInfo.RSS
}
