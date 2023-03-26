package go645

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
)

type Decoder func(buffer *bytes.Buffer) (*InformationElement, error)

func Handler(control *Control, buffer *bytes.Buffer, size byte, ver ProtoVersion) (InformationElement, error) {
	//从站响应异常响应
	if control == nil {
		return nil, errors.New("未知错误")
	}
	if control.IsState(ver, SlaveErr) {
		return nil, DecodeException(buffer)
	}
	//从站读正确响应
	if control.IsState(ver, Read) {
		return DecodeRead(buffer, int(size), ver), nil
	}
	//佳和强制联机
	if control.Data == 0x8a {
		return DecodeNullData(buffer), nil
	}
	return nil, errors.New("未定义的数据类型")
}
func Decode(ver ProtoVersion, buffer *bytes.Buffer) (*Protocol, error) {
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(buffer, binary.LittleEndian, data)
	}
	p := new(Protocol)
	read(&p.Start)
	p.Address, err = DecodeAddress(buffer, 6)
	read(&p.Start2)
	p.Control, err = DecodeControl(buffer)
	read(&p.DataLength)
	p.Data, err = Handler(p.Control, buffer, p.DataLength, ver)
	read(&p.CS)
	read(&p.End)
	if err != nil {
		log.Print(err.Error())
	}
	return p, err
}
func DecodeAddress(buffer *bytes.Buffer, size int) (Address, error) {
	var a Address

	value := make([]byte, size)
	if err := binary.Read(buffer, binary.LittleEndian, &value); err != nil {
		return a, err
	}
	{
		a.value = value
		a.strValue = Bcd2Number(a.value)
	}
	return a, nil
}
func DecodeData(buffer *bytes.Buffer, size byte) (*ReadData, error) {
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(buffer, binary.LittleEndian, data)
	}
	data := new(ReadData)
	var dataType []byte
	dataValue := make([]byte, size-4)
	read(&dataType)
	//for i, j := 0, len(dataType)-1; i < j; i, j = i+1, j-1 {
	//	dataType[j], dataType[i] = dataType[i], dataType[j]
	//}
	for index, item := range dataType {
		dataType[index] = item - 0x33
	}
	read(&dataValue)
	for index, item := range dataValue {
		dataValue[index] = item - 0x33
	}
	for i, j := 0, len(dataValue)-1; i < j; i, j = i+1, j-1 {
		dataValue[i], dataValue[j] = dataValue[j], dataValue[i]
	}

	//瞬时功率及当前需量最高位表示方向，0正，1负，三相三线B相为0。取值范围：0.0000～79.9999。
	//表内温度 最高位0表示零上，1表示零下。取值范围：0.0～799.9
	//电流最高位表示方向，0正,1负,取值范围：0.000～799.999。功率因数最高位表示方向，0正，1负，取值范围
	data.bcdValue = Bcd2Number(dataValue)
	data.rawValue = dataValue
	data.dataType = dataType
	return data, nil
}
func DecoderData(buffer *bytes.Buffer, size int) (*bytes.Buffer, error) {
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(buffer, binary.LittleEndian, data)
	}
	var value = make([]byte, size)
	read(value)
	for i, j := 0, len(value)-1; i <= j; i, j = i+1, j-1 {
		value[i], value[j] = value[j]-0x33, value[i]-0x33
	}

	return bytes.NewBuffer(value), nil
}
func DecodeRead(buffer *bytes.Buffer, size int, ver ProtoVersion) InformationElement {
	df, _ := DecoderData(buffer, size)
	var err error
	read := func(data interface{}) {
		if err != nil {
			return
		}
		err = binary.Read(df, binary.LittleEndian, data)
	}
	data := new(ReadData)
	dataTypeLen := 0
	switch ver {
	case Ver2007:
		dataTypeLen = 4
	case Ver1997:
		dataTypeLen = 2
	}
	var dataType = make([]byte, dataTypeLen)
	dataValue := make([]byte, size-dataTypeLen)
	read(&dataValue)
	read(&dataType)
	for i, j := 0, len(dataType)-1; i < j; i, j = i+1, j-1 {
		dataType[j], dataType[i] = dataType[i], dataType[j]
	}
	//瞬时功率及当前需量最高位表示方向，0正，1负，三相三线B相为0。取值范围：0.0000～79.9999。
	//表内温度 最高位0表示零上，1表示零下。取值范围：0.0～799.9
	//电流最高位表示方向，0正,1负,取值范围：0.000～799.999。功率因数最高位表示方向，0正，1负，取值范围
	//判断最高位是否为0
	if ver == Ver2007 {
		if dataType[3] == 0x02 && dataType[2] >= 0x3 && dataType[2] <= 0x6 {
			if IsStateUin8(dataValue[0], 7) {
				data.Negative = true
				dataValue[0] = dataValue[0] << 1 >> 1
			}
		}
		if dataType[3] == 02 && dataType[2] == 0x80 && dataType[0] >= 0x4 && dataType[0] <= 0x7 {
			if IsStateUin8(dataValue[0], 7) {
				data.Negative = true
				dataValue[0] = dataValue[0] << 1 >> 1
			}
		}
	}

	data.bcdValue = Bcd2Number(dataValue)
	data.dataType = dataType
	return data
}
func DecodeException(buffer *bytes.Buffer) error {
	var data uint16
	err := binary.Read(buffer, binary.LittleEndian, &data)
	if err != nil {
		log.Print(err)
		return nil
	}
	return &Exception{data}
}
func DecodeNullData(*bytes.Buffer) InformationElement {
	return NullData{}
}
func IsStateUin8(des uint8, i int) bool {
	return (des & (0x01 << i)) != 0
}

func reverse(s []byte) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
