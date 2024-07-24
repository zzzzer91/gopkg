package langx

import "unsafe"

func StringToBytes(str string) []byte {
	if str == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(str), len(str))
}

func BytesToString(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(bs), len(bs))
}

func Reverse(s string) string {
	bs := []byte(s)
	l, r := 0, len(s)-1
	for l < r {
		bs[l], bs[r] = bs[r], bs[l]
		l++
		r--
	}
	return BytesToString(bs)
}
