package simulator

import (
	"math/rand"
	"time"
)

// func GetExponentialInterval(lambda float64) time.Duration {
// 	expRandom := rand.ExpFloat64() / lambda
// 	return time.Duration(expRandom * float64(time.Second))
// }

// 가우시안 분포(정규분포)를 따르는 서비스 시간을 반환
func GetGaussianFigure(mean float64, stdDev float64) time.Duration {
	GaussianRandom := rand.NormFloat64()*stdDev + mean
	return time.Duration(GaussianRandom * float64(time.Second))
}
