package go645

import (
	"time"

	"github.com/goburrow/serial"
)

// ClientProviderOption client provider option for user.
type ClientProviderOption func(ClientProvider)

// WithLogProvider set logger provider.
func WithLogProvider(provider LogProvider) ClientProviderOption {
	return func(p ClientProvider) {
		p.setLogProvider(provider)
	}
}

func WithLogSaver(l LogSaver) ClientProviderOption {
	return func(p ClientProvider) {
		p.setLogSaver(l)
	}
}

// WithEnableLogger enable log output when you has set logger.
func WithEnableLogger() ClientProviderOption {
	return func(p ClientProvider) {
		p.LogMode(true)
	}
}

// WithSerialConfig set serial config, only valid on serial.
func WithSerialConfig(config serial.Config) ClientProviderOption {
	return func(p ClientProvider) {
		p.setSerialConfig(config)
	}
}

func WithPrefixHandler(prefixHandler PrefixHandler) ClientProviderOption {
	return func(p ClientProvider) {
		p.setPrefixHandler(prefixHandler)
	}
}

// WithTCPTimeout set tcp Connect & Read timeout, only valid on TCP.
func WithTCPTimeout(t time.Duration) ClientProviderOption {
	return func(p ClientProvider) {
		p.setTCPTimeout(t)
	}
}
