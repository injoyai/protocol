package knx

import (
	"errors"
	"fmt"
)

const (

	// Succ indicates a successful operation. 成功
	Succ Code = 0x00

	// ErrHostProtocolType indicates an unsupported host protocol.
	ErrHostProtocolType Code = 0x01

	// ErrVersionNotSupported indicates an unsupported KNXnet/IP protocol version.
	ErrVersionNotSupported Code = 0x02

	// ErrSequenceNumber indicates that an out-of-order sequence number has been received.
	ErrSequenceNumber Code = 0x04

	// ErrConnectionID indicates that there is no active data connection with given ID.
	ErrConnectionID Code = 0x21

	// ErrConnectionType indicates an unsupported connection type.
	ErrConnectionType Code = 0x22

	// ErrConnectionOption indicates an unsupported connection option.
	ErrConnectionOption Code = 0x23

	// ErrNoMoreConnections is returned by a Tunnelling Server when it cannot accept more
	// connections.
	// 网关连接到达上限
	ErrNoMoreConnections Code = 0x24

	// ErrNoMoreUniqueConnections is returned by a Tunnelling Server when it has no free Individual
	// Address available that could be used by the connection.
	ErrNoMoreUniqueConnections Code = 0x25

	// ErrDataConnection indicates an error with a data connection.
	ErrDataConnection Code = 0x26

	// ErrKNXConnection indicates an error with a KNX connection.
	ErrKNXConnection Code = 0x27

	// ErrTunnellingLayer indicates an unsupported tunnelling layer.
	ErrTunnellingLayer Code = 0x29
)

// Code 状态码
type Code uint8

// Err 错误信息,翻译成中文,可能不是很精准
func (this Code) Err() error {
	switch this {
	case Succ:
		return nil
	case ErrHostProtocolType:
		return errors.New("不支持的主机协议")
	case ErrVersionNotSupported:
		return errors.New("不支持的协议版本")
	case ErrSequenceNumber:
		return errors.New("无效的序列号")
	case ErrConnectionID:
		return errors.New("无效的连接ID")
	case ErrConnectionType:
		return errors.New("不支持的连接类型")
	case ErrConnectionOption:
		return errors.New("不支持的连接选项")
	case ErrNoMoreConnections:
		return errors.New("连接已达上限")
	case ErrNoMoreUniqueConnections:
		return errors.New("没有空闲的唯一地址")
	case ErrDataConnection:
		return errors.New("数据连接错误")
	case ErrKNXConnection:
		return errors.New("KNX连接错误")
	case ErrTunnellingLayer:
		return errors.New("不支持的隧道层")
	default:
		return fmt.Errorf("未知错误,错误码:%x", this)
	}
}
