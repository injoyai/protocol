package ads

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"io"
)

/*

参考
https://weibo.com/ttarticle/p/show?id=2313501000014301876463058653

*/

/*
Pkg
小端

读数据示例
头部信息 00 00 2c 00 00 00		44字节
目的地址 05 16 a7 68 01 01		5.22.167.104.1.1
目的端口 21 03					8451
出发地址 a9 fe a4 a3 01 01		169.254.164.163.1.1
出发端口 76 80					30336
命令类型 02 00					读取数据
状态标识 04 00					ADSCommand
命令长度 0c 00 00 00				12
错误编码 00 00 00 00				00
命令序号 01 00 00 00				序号
索引分组 20 40 00 00				分组地址16416
索引偏移 00 00 00 00				偏移量0
读取长度 01 00 00 00				1字节
*/
type Pkg struct {
	Header
	AMS

	/*
		Data
		The ADS data range contains the parameter of the single ADS
		commands. The structure of the data array depends on the ADS
		command. Some ADS commands require no additional data.
	*/
	Data []byte
}

func (this *Pkg) Bytes() g.Bytes {
	data := []byte(nil)
	dataBytes := this.Data
	amsBytes := this.AMS.Bytes(uint32(len(dataBytes)))
	headerBytes := this.Header.Bytes(uint32(len(amsBytes) + len(dataBytes)))
	data = append(data, headerBytes...)
	data = append(data, amsBytes...)
	data = append(data, dataBytes...)
	return data
}

const (
	Prefix uint16 = 0x0000
)

/*
Header 6字节
Contains the length of the data packet.
*/
type Header struct {
	Reserved uint16
}

func (this Header) Bytes(length uint32) g.Bytes {
	data := []byte(nil)
	data = append(data, g.Bytes(conv.Bytes(this.Reserved)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(length)).Reverse()...)
	return data
}

/*
AMS 32字节
此处包含了通讯的发送方和接收方地址，以及ADS错误代码、ADS命令代码和其他一些信息
The AMS/TCP-Header contains the addresses of the transmitter and
receiver. In addition the AMS error code , the ADS command Id and
some other information.
*/
type AMS struct {
	Target     [6]byte
	TargetPort uint16
	Source     [6]byte
	SourcePort uint16
	CommandID  CommandID
	Flag       Flag
	ErrorCode  uint32
	InvokeID   uint32
}

func (this *AMS) Bytes(length uint32) g.Bytes {
	data := []byte(nil)
	data = append(data, this.Target[:]...)
	data = append(data, conv.Bytes(this.TargetPort)...)
	data = append(data, this.Source[:]...)
	data = append(data, conv.Bytes(this.SourcePort)...)
	data = append(data, this.CommandID.Bytes().Reverse()...)
	data = append(data, this.Flag.Bytes().Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(length)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.ErrorCode)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.InvokeID)).Reverse()...)
	return data
}

type CommandID uint16

func (this CommandID) Bytes() g.Bytes {
	return []byte{byte(this >> 8), byte(this)}
}

func (this CommandID) Name() string {
	switch this {
	case ReadInfo:
		return "读取名称和版本号"
	case Read:
		return "读数据"
	case Write:
		return "写数据"
	case ReadState:
		return "读取状态"
	case WriteControl:
		return "改变状态"
	case AddNotice:
		return "添加通知"
	case DelNotice:
		return "删除通知"
	case GetNotice:
		return "获取通知"
	case ReadWrite:
		return "同时读写"
	}
	return "未知"
}

const (
	ReadInfo     CommandID = 0x0001 //读取ADS设备的名称和版本号
	Read         CommandID = 0x0002 //读取ADS设备的数据
	Write        CommandID = 0x0003 //写入ADS设备的数据
	ReadState    CommandID = 0x0004 //读取ADS设备的状态
	WriteControl CommandID = 0x0005 //改变ADS设备的状态3
	AddNotice    CommandID = 0x0006 //创建一个通知
	DelNotice    CommandID = 0x0007 //删除一个通知
	GetNotice    CommandID = 0x0008 //获取一个通知,不知道是否理解对,数据将单独的从ADS设备发送到ADS客户端
	ReadWrite    CommandID = 0x0009 //同时读写
)

type Flag struct {
	Broadcast      bool //是否是广播
	InitCommand    bool
	UDPCommand     bool //是否是UDP协议,否则是TCP协议
	TimestampAdded bool
	HighPriority   bool
	SysCommand     bool
	ADSCommand     bool
	NoReturn       bool
	Response       bool //是否是响应,否则是请求
}

func (this Flag) Bytes() g.Bytes {
	b := uint16(0)
	if this.Broadcast {
		b |= 0x8000
	}
	if this.InitCommand {
		b |= 0x0080
	}
	if this.UDPCommand {
		b |= 0x0040
	}
	if this.TimestampAdded {
		b |= 0x0020
	}
	if this.HighPriority {
		b |= 0x0010
	}
	if this.SysCommand {
		b |= 0x0008
	}
	if this.ADSCommand {
		b |= 0x0004
	}
	if this.NoReturn {
		b |= 0x0002
	}
	if this.Response {
		b |= 0x0001
	}
	return conv.Bytes(b)
}

/*




 */

func Decode(bs []byte) (*Pkg, error) {
	if len(bs) < 38 {
		return nil, errors.New("基础数据长度小于38字节")
	}
	if conv.Uint16(bs[:2]) != Prefix {
		return nil, fmt.Errorf("帧头错误,应为%x", Prefix)
	}
	totalLength := conv.Uint32(bs[2:6])
	if uint32(len(bs)) != totalLength+6 {
		return nil, fmt.Errorf("帧长度错误,预期(%d),得到(%d)", totalLength+6, len(bs))
	}
	p := &Pkg{
		Header: Header{
			Reserved: conv.Uint16(bs[:2]),
		},
		AMS: AMS{
			Target:     [6]byte{bs[6], bs[7], bs[8], bs[9], bs[10], bs[11]},
			TargetPort: conv.Uint16(bs[12:14]),
			Source:     [6]byte{bs[14], bs[15], bs[16], bs[17], bs[18], bs[19]},
			SourcePort: conv.Uint16(bs[20:22]),
			CommandID:  CommandID(conv.Uint16(bs[22:24])),
			Flag: Flag{
				Broadcast:      bs[25]&0x80 == 1,
				InitCommand:    bs[24]&0x80 == 1,
				UDPCommand:     bs[24]&0x40 == 1,
				TimestampAdded: bs[24]&0x20 == 1,
				HighPriority:   bs[24]&0x10 == 1,
				SysCommand:     bs[24]&0x08 == 1,
				ADSCommand:     bs[24]&0x04 == 1,
				NoReturn:       bs[24]&0x02 == 1,
				Response:       bs[24]&0x01 == 1,
			},
			ErrorCode: conv.Uint32(bs[26:30]),
			InvokeID:  conv.Uint32(bs[30:34]),
		},
		Data: bs[34:],
	}

	return p, nil
}

func ReadFunc(r *bufio.Reader) (result []byte, err error) {
	for {
		buf := make([]byte, 2)
		if _, err := io.ReadAtLeast(r, buf, 2); err != nil {
			return nil, err
		}
		if conv.Uint16(buf) == Prefix {
			result = append(result, buf...)
			buf = make([]byte, 4)
			if _, err := io.ReadAtLeast(r, buf, 4); err != nil {
				return nil, err
			}
			length := int(conv.Uint32(buf)) //精度丢失
			result = append(result, buf...)
			buf = make([]byte, length)
			if _, err := io.ReadAtLeast(r, buf, length); err != nil {
				return nil, err
			}
			result = append(result, buf...)
			return
		}
	}
}
