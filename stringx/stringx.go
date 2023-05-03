package stringx

import "github.com/zzzzer91/gopkg/conv"

func Reverse(s string) string {
	bs := []byte(s)
	l, r := 0, len(s)-1
	for l < r {
		bs[l], bs[r] = bs[r], bs[l]
		l++
		r--
	}
	return conv.Bytes2string(bs)
}
