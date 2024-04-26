package simulator

import (
	"math"
	"runtime"
	"sync"
)

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

// simulateHighMemoryUsage 함수는 메모리 사용량을 인위적으로 증가시킵니다.
// sizeInMB 파라미터는 메가바이트 단위로 메모리 할당량을 결정합니다.
func IncreaseMemoryUsage(sizeInMB int) {
	// 메가바이트 당 약 1,000,000 바이트를 할당합니다 (실제로는 조금 더 정확하게 1024*1024를 사용할 수 있습니다).
	slice := make([]byte, sizeInMB*1000000)
	for i := range slice {
		slice[i] = byte(i % 256) // 메모리 할당을 확실히 사용하도록 초기화
	}

	// GC가 슬라이스를 회수하지 않도록 슬라이스를 유지합니다.
	// 메모리 사용량을 보려면 이 함수가 반환된 후에도 slice를 유지하세요.
	runtime.GC() // GC를 호출하여 메모리 사용량을 정확히 확인할 수 있습니다.
}
