package jsonx

import (
	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

func Marshal(obj any) ([]byte, error) {
	res, err := sonic.ConfigDefault.Marshal(obj)
	if err != nil {
		return nil, errors.Wrap(err, "Marshal error")
	}
	return res, nil
}

func MarshalToString(obj any) (string, error) {
	res, err := sonic.ConfigDefault.MarshalToString(obj)
	if err != nil {
		return "", errors.Wrap(err, "MarshalToString error")
	}
	return res, nil
}

func Unmarshal(data []byte, obj any) error {
	err := sonic.ConfigDefault.Unmarshal(data, obj)
	if err != nil {
		return errors.Wrap(err, "Unmarshal error")
	}
	return nil
}

func UnmarshalFromString(data string, obj any) error {
	err := sonic.ConfigDefault.UnmarshalFromString(data, obj)
	if err != nil {
		return errors.Wrap(err, "UnmarshalFromString error")
	}
	return nil
}

func Marshal2(obj any) []byte {
	res, _ := sonic.ConfigDefault.Marshal(obj)
	return res
}

func MarshalToString2(obj any) string {
	res, _ := sonic.ConfigDefault.MarshalToString(obj)
	return res
}

func Unmarshal2(data []byte, obj any) {
	_ = sonic.ConfigDefault.Unmarshal(data, obj)
}

func UnmarshalFromString2(data string, obj any) {
	_ = sonic.ConfigDefault.UnmarshalFromString(data, obj)
}
