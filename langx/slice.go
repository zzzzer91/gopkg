package langx

import (
	"strconv"

	"github.com/pkg/errors"
)

func UniqueSlice[T comparable](slice []T) []T {
	seen := make(map[T]struct{})
	var res []T
	for _, v := range slice {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			res = append(res, v)
		}
	}
	return res
}

// pageNum starts from 0
func Page[T any](slice []T, pageNum, pageSize, maxPageSize int) ([]T, bool) {
	if pageSize == 0 || pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	offset := pageNum * pageSize
	if len(slice) == 0 || offset >= len(slice) {
		return make([]T, 0), false
	}
	limit := offset + int(pageSize)
	hasNext := true
	if limit >= len(slice) {
		limit = len(slice)
		hasNext = false
	}
	return slice[offset:limit], hasNext
}

func String2Int(a []string) ([]int, error) {
	res := make([]int, 0, len(a))
	for _, e := range a {
		n, err := strconv.Atoi(e)
		if err != nil {
			return nil, errors.WithStack(err)
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
