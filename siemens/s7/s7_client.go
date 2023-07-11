package s7

import (
	"github.com/injoyai/conv"
	"math/rand"
	"time"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

// WriteCPUStart 将CPU设置成运行状态
func WriteCPUStart() []byte {
	return []byte{0x03, 0x00, 0x00, 0x25, 0x02, 0xf0, 0x80,
		0x32, 0x01, 0x00, 0x00, 0x00, 0x1b, 0x00, 0x14, 0x00, 0x00,
		byte(PLCStart), 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xfd, 0x00, 0x00, 0x09, 0x50, 0x5f, 0x50,
		0x52, 0x4f, 0x47, 0x52, 0x41, 0x4d}
}

// WriteCPUStop 将CPU设置成停止状态
func WriteCPUStop() []byte {
	return []byte{0x03, 0x00, 0x00, 0x21, 0x02, 0xf0, 0x80,
		0x32, 0x01, 0x00, 0x00, 0x02, 0x73, 0x00, 0x10, 0x00, 0x00,
		byte(PLCStop), 0x00, 0x00, 0x00, 0x00, 0x00, 0x09, 0x50, 0x5f, 0x50,
		0x52, 0x4f, 0x47, 0x52, 0x41, 0x4d}
}

// ReadCPUStatus 读取CPU状态
func ReadCPUStatus() *Request {
	return NewRead(0x000f18, 1, Byte, CPU)
}

// ReadMBit 读取M区位
func ReadMBit(addr uint32) *Request {
	return NewRead(addr, 1, Bit, Marker)
}

// ReadMByte 读取M区字节
func ReadMByte(addr uint32, size uint16) *Request {
	return NewRead(addr, size, Byte, Marker)
}

// ReadIByte 读取输入字节
func ReadIByte(addr uint32, size uint16) *Request {
	return NewRead(addr, size, Byte, Input)
}

// ReadDBBit 读取DB位
func ReadDBBit(addr uint32) *Request {
	return NewRead(addr, 1, Byte, DataBlock)
}

// ReadDBByte 读DB字节
func ReadDBByte(addr uint32, size uint16) *Request {
	return NewRead(addr, size, Byte, DataBlock)
}

// WriteMBit 写M位
func WriteMBit(addr uint32, value []byte) *Request {
	return NewWrite(addr, value, Bit, Marker)
}

// WriteMByte 写M字节
func WriteMByte(addr uint32, value []byte) *Request {
	return NewWrite(addr, value, Byte, Marker)
}

// WriteIBit 写I位
func WriteIBit(addr uint32, value []byte) *Request {
	return NewWrite(addr, value, Bit, Input)
}

// WriteIByte 写I字节
func WriteIByte(addr uint32, value []byte) *Request {
	return NewWrite(addr, value, Byte, Input)
}

// WriteDBBit 写DB位
func WriteDBBit(addr uint32, value bool) *Request {
	return NewWrite(addr, []byte{conv.SelectUint8(value, 1, 0)}, Byte, DataBlock)
}

// WriteDBByte 写DB字节
func WriteDBByte(addr uint32, value []byte) *Request {
	return NewWrite(addr, value, Byte, DataBlock)
}

// NewRead 读数据
func NewRead(addr uint32, size uint16, dataType DataType, block Block) *Request {
	return &Request{
		MsgID: uint16(r.Uint32()),
		Param: Param{
			OrderType: ReadVar,
			Area: Area{
				DataType: dataType,
				Block:    block,
				Addr:     addr,
				Size:     size,
			},
		},
	}
}

// NewWrite 写数据
func NewWrite(addr uint32, value []byte, dataType DataType, block Block) *Request {
	return &Request{
		MsgID: uint16(r.Uint32()),
		Param: Param{
			OrderType: WriteVar,
			Area: Area{
				DataType: dataType,
				Block:    block,
				Addr:     addr,
				Size:     uint16(len(value)),
			},
		},
		Write: Write{Value: value},
	}
}
