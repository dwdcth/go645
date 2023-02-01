package go645

import (
	"bytes"
	"io"
	"log"
	"time"
	"strings"
	"fmt"
)

var _ ClientProvider = (*RTUClientProvider)(nil)

type RTUClientProvider struct {
	serialPort
	logger
	PrefixHandler
}

func (sf *RTUClientProvider) setPrefixHandler(handler PrefixHandler) {
	sf.PrefixHandler = handler
}

//SendAndRead 发送数据并读取返回值
func (sf *RTUClientProvider) SendAndRead(p *Protocol) (aduResponse []byte, err error) {
	bf := bytes.NewBuffer(make([]byte, 0))
	err = sf.EncodePrefix(bf)
	if err != nil {
		return nil, err
	}

	err = p.Encode(bf)
	if err != nil {
		return nil, err
	}
	return sf.SendRawFrameAndRead(bf.Bytes())
}
func (sf *RTUClientProvider) Send(p *Protocol) (err error) {
	bf := bytes.NewBuffer(make([]byte, 0))
	err = sf.EncodePrefix(bf)
	if err != nil {
		return err
	}
	err = p.Encode(bf)
	if err != nil {
		return err
	}
	return sf.SendRawFrame(bf.Bytes())
}

//ReadRawFrame 读取返回数据
func (sf *RTUClientProvider) ReadRawFrame() (aduResponse []byte, err error) {
	fe, err := sf.DecodePrefix(sf.port)
	if err != nil {
		log.Printf(err.Error())
		return nil, err
	}
	//head := make([]byte, 14)
	//size, err := io.ReadAtLeast(sf.port, head, 14)
	//if err != nil || size != 10 {
	//	return nil, err
	//}
	////数据域+2
	//expLen := head[9] + 2
	//playLoad := make([]byte, expLen)
	//if _, err := io.ReadAtLeast(sf.port, playLoad, int(expLen)); err != nil {
	//	return nil, err
	//}
	////拆包器重新实现
	//content := append(head, playLoad...)


	// 延时在 20ms <= Td <= 500ms
	// 读取报文
	// 1. 先读取 14 个字节 	fefefefe 68 190002031122 68 91 84
	var data [1024]byte
	var n int
	var n1 int

	n, err = ReadAtLeast(sf.port, data[:], 14, 500*time.Millisecond)
	if err != nil {
		return
	}
	// 帧起始符长度
	frontLen := 0
	for _, b := range data {
		if b == 0xfe {
			frontLen++
		} else {
			break
		}
	}
	L := int(data[frontLen+9]) // 数据域长度
	// 总字节数
	bytesToRead := frontLen + 1 + 6 + 1 + 1 + 1 + L + 1 + 1
	// 读取剩余字节
	if n < bytesToRead {
		if bytesToRead > n {
			n1, err = io.ReadFull(sf.port, data[n:bytesToRead])
			n += n1
		}
	}
	a := fmt.Sprintf("%x", data[:4])
	if strings.ToUpper(a) == "FEFEFEFE"{

		aduResponse = data[4:n]
	}else {
		aduResponse = data[:n]
	}
	sf.Debugf("rec <==[% x]", append(fe, aduResponse...))
	return aduResponse, nil
}
func (sf *RTUClientProvider) SendRawFrameAndRead(aduRequest []byte) (aduResponse []byte, err error) {
	sf.mu.Lock()
	defer sf.mu.Unlock()
	if err = sf.connect(); err != nil {
		return
	}
	err = sf.SendRawFrame(aduRequest)
	if err != nil {
		log.Printf(err.Error())
		_ = sf.close()
		return
	}
	return sf.ReadRawFrame()
}
func (sf *RTUClientProvider) SendRawFrame(aduRequest []byte) (err error) {
	if err = sf.connect(); err != nil {
		return
	}
	// Send the request
	sf.Debugf("sending ==> [% x]", aduRequest)
	//发送数据
	_, err = sf.port.Write(aduRequest)
	return err
}

// NewRTUClientProvider allocates and initializes a RTUClientProvider.
// it will use default /dev/ttyS0 19200 8 1 N and timeout 1000
func NewRTUClientProvider(opts ...ClientProviderOption) *RTUClientProvider {
	p := &RTUClientProvider{
		logger:        newLogger("645RTUMaster => "),
		PrefixHandler: &DefaultPrefix{},
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

// calculateDelay roughly calculates time needed for the next frame.
// See MODBUS over Serial Line - Specification and Implementation Guide (page 13).
func (sf *RTUClientProvider) calculateDelay(chars int) time.Duration {
	var characterDelay, frameDelay int // us

	if sf.BaudRate <= 0 || sf.BaudRate > 19200 {
		characterDelay = 750
		frameDelay = 1750
	} else {
		characterDelay = 15000000 / sf.BaudRate
		frameDelay = 35000000 / sf.BaudRate
	}
	return time.Duration(characterDelay*chars+frameDelay) * time.Microsecond
}
