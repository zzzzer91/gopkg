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

func MustMarshal(obj any) []byte {
	res, err := sonic.ConfigDefault.Marshal(obj)
	if err != nil {
		panic(err)
	}
	return res
}

func MustMarshalToString(obj any) string {
	res, err := sonic.ConfigDefault.MarshalToString(obj)
	if err != nil {
		panic(err)
	}
	return res
}

func MustUnmarshal(data []byte, obj any) {
	err := sonic.ConfigDefault.Unmarshal(data, obj)
	if err != nil {
		panic(err)
	}
}

func MustUnmarshalFromString(data string, obj any) {
	err := sonic.ConfigDefault.UnmarshalFromString(data, obj)
	if err != nil {
		panic(err)
	}
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
