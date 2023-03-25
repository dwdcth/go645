package go645

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"testing"
)

// TestRead 测试上报请求
func TestDecode(t *testing.T) {
	str := "681401003182216891083333343339333333f116"
	decodeString, err := hex.DecodeString(str)
	if err != nil {
		return
	}
	p, _ := Decode(Ver2007, bytes.NewBuffer(decodeString))
	if p.Address.strValue != "140100318221" {
		t.Errorf("地址解析错误")
	}
	if p.Address.GetLen() != 6 {
		t.Errorf("长度错误")
	}
	if !p.Control.IsState(Ver2007, IsSlave) || !p.Control.IsState(Ver2007, Read) {
		t.Errorf("状态解析错误")
	}
	fmt.Printf("%f \n", p.Data.(*ReadData).GetFloat64Value())
	if p.GetLen() != 0x08 {
		t.Errorf("长度错误")
	}

}

// TestRead 测试解析读请求
func TestRead(t *testing.T) {
	str := "6881040007213a68910733353535333333ba16"
	data := make([]byte, 0)
	c := NewControl()
	c.SetState(Ver2007, Read)
	r := ReadRequest(Ver2007, NewAddress("610100000000", BigEndian), 0x00_01_00_00)
	bf := bytes.NewBuffer(data)
	_ = r.Encode(bf)
	decodeString, _ := hex.DecodeString(str)
	p2, _ := Decode(Ver2007, bytes.NewBuffer(decodeString))
	print(p2.Data.(*ReadData).GetValue() + "\n")
	print(p2.Data.(*ReadData).GetDataTypeStr() + "\n")
	fmt.Printf("%f", p2.Data.(*ReadData).GetFloat64Value())
}

// TestSend 测试发送
func TestSend(t *testing.T) {
	str := "68610100000000681104333334331416"
	data := make([]byte, 0)
	c := NewControl()
	c.SetState(Ver2007, Read)
	r := ReadRequest(Ver2007, NewAddress("610100000000", BigEndian), 0x00_01_00_00)
	bf := bytes.NewBuffer(data)
	_ = r.Encode(bf)
	p, _ := Decode(Ver2007, bf)
	decodeString, _ := hex.DecodeString(str)
	p2, _ := Decode(Ver2007, bytes.NewBuffer(decodeString))
	print(p.Data.(*ReadData).GetValue())
	AssertEquest("地址错误", p2.Address.strValue, p.Address.strValue, t)
	AssertEquest("校验码错误", p.CS, p2.CS, t)

}
func TestLEnd(t *testing.T) {
	str := "68610100000000681104333334331416"
	data := make([]byte, 0)
	c := NewControl()
	c.SetState(Ver2007, Read)
	r := ReadRequest(Ver2007, NewAddress("610100000000", LittleEndian), 0x00_01_00_00)
	bf := bytes.NewBuffer(data)
	_ = r.Encode(bf)
	p, _ := Decode(Ver2007, bf)
	decodeString, _ := hex.DecodeString(str)
	p2, _ := Decode(Ver2007, bytes.NewBuffer(decodeString))
	AssertEquest("地址错误", p2.Address.GetStrAddress(LittleEndian), "000000000161", t)
	AssertEquest("地址错误", p2.Address.GetStrAddress(BigEndian), "610100000000", t)
	AssertEquest("校验码错误", p.CS, p2.CS, t)
}
func Assert(msg string, assert func() bool, t *testing.T) {
	if !assert() {
		t.Errorf(msg)
	}
}
func AssertEquest(msg string, exp interface{}, act interface{}, t *testing.T) {
	Assert(msg, func() bool { return exp == act }, t)
}

func AssertState(assert func() bool, t *testing.T) {
	Assert("状态解析错误", assert, t)
}
func TestControl_IsState(t *testing.T) {
	d := NewControlValue(0x0A)
	println(d.IsStates(Ver2007, IsSlave))
	println(d.IsStates(Ver2007, Broadcast))
	println(d.IsStates(Ver2007, Read))

}
func TestControl(t *testing.T) {
	c := &Control{}
	if c.getLen() != 1 {
		t.Errorf("长度错误")
	}
	c.SetState(Ver2007, SlaveErr)
	if !c.IsState(Ver2007, SlaveErr) {
		t.Errorf("设置错误")
	}
	c.Reset()
	if c.IsState(Ver2007, SlaveErr) {
		t.Errorf("复归错误")
	}
	c.SetStates(Ver2007, SlaveErr, IsSlave, HasNext, Retain, Broadcast, ReadNext, ReadAddress)
	if !c.IsStates(Ver2007, SlaveErr, IsSlave) {
		t.Errorf("设置错误")
	}
}
func (c *Control) TestErr(buffer *bytes.Buffer) error {

	var bf *bytes.Buffer
	r := ReadRequest(Ver2007, NewAddress("610100000000", LittleEndian), 0x00_01_00_00)
	_ = r.Encode(bf)

	if err := binary.Write(buffer, binary.BigEndian, c.Data); err != nil {
		s := fmt.Sprintf("Control , %v", err)
		return errors.New(s)
	}
	return nil
}
func TestReadResponse(t *testing.T) {
	rp := ReadResponse(Ver2007, NewAddress("610100000000", LittleEndian), 0x00_01_00_00, NewControl(), "200")

	if rp.Address.GetStrAddress(LittleEndian) != "610100000000" {
		t.Errorf("地址错误")
	}
}
