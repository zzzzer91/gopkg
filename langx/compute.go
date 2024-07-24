package langx

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
