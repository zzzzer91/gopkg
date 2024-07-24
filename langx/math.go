package langx

import (
	"fmt"
	"strconv"

	"golang.org/x/exp/constraints"
)

func Round(value float64, n int) float64 {
	v, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(n)+"f", value), 64)
	return v
}

func CountDigits[T int | int64 | int32](n T) T {
	if n == 0 {
		return 1
	}
	count := T(0)
	for n != 0 {
		n /= 10
		count++
	}
	return count
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}
