package langx

import (
	"fmt"
	"strconv"
)

func Round(value float64, n int) float64 {
	v, _ := strconv.ParseFloat(fmt.Sprintf("%."+strconv.Itoa(n)+"f", value), 64)
	return v
}
