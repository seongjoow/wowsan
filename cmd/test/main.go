package main

import (
	"fmt"
	"math"
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	// "github.com/shirou/gopsutil/v3/cpu"
	// "github.com/shirou/gopsutil/v3/mem"
	// "github.com/shirou/gopsutil/v3/net"
	// "github.com/shirou/gopsutil/net"
)

// func getCPUInfo() {
// 	cpuInfo, _ := cpu.Percent(time.Second, false)
// 	fmt.Println("CPU Usage: ", cpuInfo[0])
// }

// func getMemoryInfo() {
// 	v, _ := mem.VirtualMemory()
// 	fmt.Printf("메모리 사용률: %.2f%%\n", v.UsedPercent)
// }

// func getNetInfo() {
// 	netInfo, _ := net.IOCounters(true)
// 	fmt.Println(netInfo)
// }

func main() {
	pid := os.Getpid()

	IncreaseCpuUsage()

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
	}
}

func IncreaseCpuUsage() {
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ { // 동시에 여러 고루틴을 실행하여 CPU 사용량 증가
		wg.Add(1)
		go func() {
			defer wg.Done()
			result := 0.0
			for j := 0.0; j < 1000000; j++ {
				result += math.Sqrt(j) // 계산 비용이 높은 연산
			}
		}()
	}
	wg.Wait()
}
