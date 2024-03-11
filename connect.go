package protocol

import (
	"github.com/injoyai/io"
	"github.com/injoyai/io/dial"
)

func NewConnect(dial io.DialFunc, options ...io.OptionClient) (*io.Client, error) {
	return io.NewDial(dial, options...)
}

func NewTCPConnect(addr string, options ...io.OptionClient) (*io.Client, error) {
	return dial.NewTCP(addr, options...)
}

func NewUDPConnect(addr string, options ...io.OptionClient) (*io.Client, error) {
	return dial.NewUDP(addr, options...)
}

func NewSerialConnect(cfg *dial.SerialConfig, options ...io.OptionClient) (*io.Client, error) {
	return dial.NewSerial(cfg, options...)
}
