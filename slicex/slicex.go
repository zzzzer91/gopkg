package slicex

import (
	"strconv"
)

func String2Int(a []string) ([]int, error) {
	res := make([]int, 0, len(a))
	for _, e := range a {
		n, err := strconv.Atoi(e)
		if err != nil {
			return nil, err
		}
		res = append(res, n)
	}
	return res, nil
}

func Int2String(a []int) []string {
	res := make([]string, 0, len(a))
	for _, e := range a {
		n := strconv.Itoa(e)
		res = append(res, n)
	}
	return res
}
