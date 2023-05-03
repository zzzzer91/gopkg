package jsonx

import "github.com/bytedance/sonic"

func Marshal(obj any) []byte {
	res, _ := sonic.ConfigDefault.Marshal(obj)
	return res
}

func MarshalToString(obj any) string {
	res, _ := sonic.ConfigDefault.MarshalToString(obj)
	return res
}
