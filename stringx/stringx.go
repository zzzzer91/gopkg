package stringx

import (
	"reflect"
	"unsafe"
)

func String2bytes(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(&bh))
}

func Bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Reverse(s string) string {
	bs := []byte(s)
	l, r := 0, len(s)-1
	for l < r {
		bs[l], bs[r] = bs[r], bs[l]
		l++
		r--
	}
	return Bytes2string(bs)
}
