package go645

import (
	"bytes"
	"io"
)

var _ PrefixHandler = (*DefaultPrefix)(nil)
var _prefix = []byte{0xfe, 0xfe, 0xfe, 0xfe}

type PrefixHandler interface {
	EncodePrefix(buffer *bytes.Buffer) error
	// DecodePrefix 前缀内容，是否是前缀，错误
	DecodePrefix(reader io.Reader) ([]byte, bool, error)
}

type DefaultPrefix struct {
}

func (d DefaultPrefix) EncodePrefix(buffer *bytes.Buffer) error {
	// 写入引导词
	buffer.Write(_prefix)
	return nil
}

func (d DefaultPrefix) DecodePrefix(reader io.Reader) ([]byte, bool, error) {
	fe := make([]byte, 4)
	_, err := io.ReadAtLeast(reader, fe, 4)
	if err != nil {
		return nil, false, err
	}
	isPrefix := true
	for i := 0; i < len(fe); i++ {
		if fe[i] != _prefix[i] {
			isPrefix = false
			break
		}
	}
	return fe, isPrefix, nil
}
