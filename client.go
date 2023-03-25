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

func (c *client) Read(ver ProtoVersion, address Address, itemCode int32) (*ReadData, bool, error) {
	resp, err := c.ClientProvider.SendAndRead(ReadRequest(ver, address, itemCode))
	if err != nil {
		return nil, false, err
	}
	decode, err := Decode(ver, bytes.NewBuffer(resp))
	if err != nil {
		return nil, false, err
	}
	return decode.Data.(*ReadData), decode.Control.IsState(ver, HasNext), err
}

// Broadcast 设备广播 todo 版本号
func (c *client) Broadcast(ver ProtoVersion, p InformationElement, control Control) error {
	var err error
	bf := bytes.NewBuffer(make([]byte, 0))
	err = p.Encode(bf)
	if err != nil {
		return err
	}
	return c.Send(NewProtocol(NewAddress(BroadcastAddress, LittleEndian), p, &control))
}

func (c *client) ReadWithBlock(ver ProtoVersion, address Address, data ReadRequestData) (*Protocol, error) {
	resp, err := c.ClientProvider.SendAndRead(ReadRequestWithBlock(ver, address, data))
	if err != nil {
		return nil, err
	}
	return Decode(ver, bytes.NewBuffer(resp))
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
