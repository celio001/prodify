package validator

import (
	"math"
)

func GetEntropy(password string) float64 {
	return getEntropy(password)
}

func getEntropy(password string) float64 {
	base := getBase(password)
	length := getLength(password)

	return logPow(float64(base), length, 2)
}

func logX(base, n float64) float64 {
	if base == 0 {
		return 0
	}
	return math.Log2(n) / math.Log2(base)
}

func logPow(expBase float64, pow int, logBase float64) float64 {
	// logb (MN) = logb M + logb N
	total := 0.0
	for i := 0; i < pow; i++ {
		total += logX(logBase, expBase)
	}
	return total
}
