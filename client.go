package go645

import (
	"bytes"
	"sync"
)

var (
	// check implements Client interface.
	_ Client = (*client)(nil)
)

type client struct {
	ClientProvider
	mu sync.Mutex
}

func (c *client) Read(address Address, itemCode int32, ver ProtoVersion) (*ReadData, bool, error) {
	resp, err := c.ClientProvider.SendAndRead(ReadRequest(address, itemCode, ver))
	if err != nil {
		return nil, false, err
	}
	decode, err := Decode(bytes.NewBuffer(resp), ver)
	if err != nil {
		return nil, false, err
	}
	return decode.Data.(*ReadData), decode.Control.IsState(ver, HasNext), err
}

// Broadcast 设备广播 todo 版本号
func (c *client) Broadcast(p InformationElement, control Control, ver ProtoVersion) error {
	var err error
	bf := bytes.NewBuffer(make([]byte, 0))
	err = p.Encode(bf)
	if err != nil {
		return err
	}
	return c.Send(NewProtocol(NewAddress(BroadcastAddress, LittleEndian), p, &control))
}

func (c *client) ReadWithBlock(address Address, data ReadRequestData, ver ProtoVersion) (*Protocol, error) {
	resp, err := c.ClientProvider.SendAndRead(ReadRequestWithBlock(address, data, ver))
	if err != nil {
		return nil, err
	}
	return Decode(bytes.NewBuffer(resp), ver)
}

// Option custom option
type Option func(c *client)

func NewClient(p ClientProvider, opts ...Option) Client {
	c := &client{ClientProvider: p}
	for _, opt := range opts {
		opt(c)
	}
	return c
}
