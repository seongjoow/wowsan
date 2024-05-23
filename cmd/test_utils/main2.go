package main

import (
	"fmt"
	"os"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"
	"github.com/shirou/gopsutil/v3/process"
	// "github.com/shirou/gopsutil/net"
)

type data struct {
	item []int
}

func getCPUInfo() {
	cpuInfo, _ := cpu.Percent(time.Second, false)
	fmt.Println("CPU Usage: ", cpuInfo[0])
}

func getMemoryInfo() {
	v, _ := mem.VirtualMemory()
	fmt.Printf("메모리 사용률: %.2f%%\n", v.UsedPercent)
}

func getNetInfo() {
	netInfo, _ := net.IOCounters(true)
	fmt.Println(netInfo)
}

func main() {
	pid := os.Getpid()
	// add data
	var d data
	go func() {
		for i := 0; i < 100000000; i++ {
			time.Sleep(1 * time.Millisecond)
			d.item = append(d.item, 1)
		}
	}()

	p, err := process.NewProcess(int32(pid))
	if err != nil {
		fmt.Println("프로세스 찾기 실패:", err)
		return
	}

	for i := 0; i < 100000000; i++ {
		time.Sleep(2 * time.Second)
		cpuPercent, _ := p.CPUPercent()
		memInfo, _ := p.MemoryInfo()
		fmt.Printf("CPU 사용률: %.2f%%\n", cpuPercent)
		fmt.Printf("메모리 사용량: %v bytes\n", memInfo.RSS)
		fmt.Println(pid)
	}

}
