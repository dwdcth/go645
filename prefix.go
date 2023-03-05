package go645

import (
	"bytes"
	"io"
)

var _ PrefixHandler = (*DefaultPrefix)(nil)

type PrefixHandler interface {
	EncodePrefix(buffer *bytes.Buffer) error
	// DecodePrefix 前缀内容，是否是前缀，错误
	DecodePrefix(reader io.Reader) ([]byte, bool, error)
}

type DefaultPrefix struct {
}

func (d DefaultPrefix) EncodePrefix(buffer *bytes.Buffer) error {

	return nil
}

func (d DefaultPrefix) DecodePrefix(reader io.Reader) ([]byte, bool, error) {
	return nil, true, nil
}
