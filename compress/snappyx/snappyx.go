package snappyx

import (
	"github.com/golang/snappy"
	"github.com/pkg/errors"
	"github.com/zzzzer91/gopkg/stringx"
)

func Decode(src []byte) ([]byte, error) {
	res, err := snappy.Decode(nil, src)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return res, nil
}

func Encode(src []byte) []byte {
	return snappy.Encode(nil, src)
}

func DecodeToString(src []byte) (string, error) {
	res, err := Decode(src)
	if err != nil {
		return "", err
	}
	return stringx.BytesToString(res), nil
}
