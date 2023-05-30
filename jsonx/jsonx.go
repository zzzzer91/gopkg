package jsonx

import (
	"io"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
)

func Marshal(obj any) ([]byte, error) {
	res, err := sonic.ConfigDefault.Marshal(obj)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func MarshalToString(obj any) (string, error) {
	res, err := sonic.ConfigDefault.MarshalToString(obj)
	if err != nil {
		return "", errors.WithStack(err)
	}
	return res, nil
}

func Unmarshal(data []byte, obj any) error {
	err := sonic.ConfigDefault.Unmarshal(data, obj)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func UnmarshalFromString(data string, obj any) error {
	err := sonic.ConfigDefault.UnmarshalFromString(data, obj)
	if err != nil {
		return errors.WithStack(err)
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
	sonic.ConfigDefault.Unmarshal(data, obj)
}

func UnmarshalFromString2(data string, obj any) {
	sonic.ConfigDefault.UnmarshalFromString(data, obj)
}

func MarshalToWriter(writer io.Writer, obj any) error {
	err := sonic.ConfigDefault.NewEncoder(writer).Encode(obj)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func UnmarshalFromReader(reader io.Reader, obj any) error {
	err := sonic.ConfigDefault.NewDecoder(reader).Decode(obj)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
