package jsonx

import (
	"io"

	"github.com/zzzzer91/gopkg/langx"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

func Decode(r io.Reader, v any) error {
	err := sonic.ConfigDefault.NewDecoder(r).Decode(v)
	if err != nil {
		return errors.Wrap(err, "json Decode error")
	}
	return nil
}

func Marshal(obj any) ([]byte, error) {
	res, err := sonic.Marshal(obj)
	if err != nil {
		return nil, errors.Wrap(err, "json Marshal error")
	}
	return res, nil
}

func MarshalToString(obj any) (string, error) {
	res, err := Marshal(obj)
	if err != nil {
		return "", err
	}
	return langx.BytesToString(res), nil
}

func Unmarshal(data []byte, obj any) error {
	err := sonic.Unmarshal(data, obj)
	if err != nil {
		return errors.Wrap(err, "json Unmarshal error")
	}
	return nil
}

func UnmarshalFromString(data string, obj any) error {
	err := Unmarshal(langx.StringToBytes(data), obj)
	if err != nil {
		return err
	}
	return nil
}

func MustDecode(r io.Reader, v any) {
	err := Decode(r, v)
	if err != nil {
		panic(err)
	}
}

func MustMarshal(obj any) []byte {
	res, err := Marshal(obj)
	if err != nil {
		panic(err)
	}
	return res
}

func MustMarshalToString(obj any) string {
	res, err := MarshalToString(obj)
	if err != nil {
		panic(err)
	}
	return res
}

func MustUnmarshal(data []byte, obj any) {
	err := Unmarshal(data, obj)
	if err != nil {
		panic(err)
	}
}

func MustUnmarshal2[T any](data []byte) T {
	var obj T
	err := Unmarshal(data, &obj)
	if err != nil {
		panic(err)
	}
	return obj
}

func MustUnmarshalFromString(data string, obj any) {
	err := UnmarshalFromString(data, obj)
	if err != nil {
		panic(err)
	}
}

func MustConvert[T any](v any) T {
	var v2 T
	MustUnmarshal(MustMarshal(v), &v2)
	return v2
}
