package go645

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

var (
	controlVersionMap = make(map[ProtoVersion]map[ControlKind]ControlType)
)

type ControlKind int

const (
	IsSlave ControlKind = iota + 1
	SlaveErr
	HasNext
	//Retain 保留
	Retain
	//Broadcast 广播校时
	Broadcast
	// ReadNext 读后续10001
	ReadNext
	//ReadAddress 读通讯地址
	ReadAddress
	//Write 写数据
	Write
	//WriteAddress 读通讯地址
	WriteAddress
	//ToChangeCommunicationRate 更改通讯速率
	ToChangeCommunicationRate
	Freeze
	//PassWord 修改密码
	PassWord
	ResetMaxDemand
	//ResetEM 电表清零
	ResetEM
	ResetEvent
	//Read 读
	Read
)

type ControlType byte

type Control struct {
	Data ControlType
}

func DecodeControl(buffer *bytes.Buffer) (*Control, error) {
	c := new(Control)
	if err := binary.Read(buffer, binary.LittleEndian, &c.Data); err != nil {
		return nil, err
	}
	return c, nil
}
func NewControl() *Control {
	return &Control{Data: 0}
}
func NewControlValue(data ControlType) *Control {
	return &Control{Data: data}
}

func (c *Control) SetState(ver ProtoVersion, state ControlKind) {
	c.Data = c.Data | controlVersionMap[ver][state]
}

// SetStates 批量设置状态
func (c *Control) SetStates(ver ProtoVersion, state ...ControlKind) {
	for _, s := range state {
		c.Data = c.Data | controlVersionMap[ver][s]
	}
}
func (c *Control) IsState(ver ProtoVersion, state ControlKind) bool {
	return (c.Data & controlVersionMap[ver][state]) == controlVersionMap[ver][state]
}

// IsStates 判断控制域
func (c *Control) IsStates(ver ProtoVersion, state ...ControlKind) bool {
	for _, s := range state {
		if !c.IsState(ver, s) {
			return false
		}
	}
	return true
}

func (c *Control) Reset() {
	c.Data = 0
}
func (c *Control) getLen() uint16 {
	return 1
}

func (c *Control) Encode(buffer *bytes.Buffer) error {
	if err := binary.Write(buffer, binary.BigEndian, c.Data); err != nil {
		s := fmt.Sprintf("Control , %v", err)
		return errors.New(s)
	}
	return nil
}
