package go645

import (
	"github.com/goburrow/serial"
	"time"
)

type ProtoVersion int

const (
	Ver1997 ProtoVersion = iota + 1
	Ver2007
)

type LogSaver interface {
	Write(url string, station string, data []byte)
}

type ClientProvider interface {
	// Connect try to connect the remote server
	Connect() error
	// IsConnected returns a bool signifying whether
	// the client is connected or not.
	IsConnected() bool
	LogMode(enable bool)
	// Close disconnect the remote server
	Close() error
	setSerialConfig(config serial.Config)
	setPrefixHandler(handler PrefixHandler)
	// setTCPTimeout set tcp connect & read timeout
	setTCPTimeout(t time.Duration)
	setLogProvider(p LogProvider)
	setLogSaver(l LogSaver)
	SendAndRead(*Protocol) (aduResponse []byte, err error)
	SendRawFrameAndRead(station string, aduRequest []byte) (aduResponse []byte, err error)
	SendRawFrame(station string, aduRequest []byte) (err error)
	ReadRawFrame() (aduResponse []byte, err error)
	Send(*Protocol) (err error)
}

// LogProvider  log message levels only Debug and Error
type LogProvider interface {
	Errorf(format string, v ...interface{})
	Debugf(format string, v ...interface{})
}
