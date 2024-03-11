package protocol

import (
	"bufio"
	"github.com/injoyai/base/g"
)

type Interface interface {
	Bytes() g.Bytes
}

type (

	// ReadFunc 分包函数,从数据流中拆分一个个数据包
	ReadFunc func(r *bufio.Reader) ([]byte, error)
)

/*

这里协议只涉及到协议,不关心传输,不管modbus-tcp用串口/蓝牙/udp等等传输数据
而传输只需要遵循io.ReadWriteCloser
数据传输和协议2个组合就能实现正常的采集控制的操作



*/
