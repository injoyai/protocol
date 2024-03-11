package knx

import (
	"bufio"
	"fmt"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"github.com/injoyai/io/buf"
	"io"
	"net"
)

/*
KNX是建筑自动化和家庭控制系统的标准
参考
https://blog.csdn.net/sixroom/article/details/130172030
https://blog.csdn.net/qq_40648827/article/details/127654947
https://github.com/vapourismo/knx-go

*/

const (

	// Prefix 数据报文的起始帧 含义是 0x06是头部的长度(固定),0x10是协议版本(暂时固定)
	Prefix uint16 = 0x0610

	// DefaultPort 默认网关端口
	DefaultPort = 3671
)

type Protocol uint8

const (
	UDP Protocol = 0x01
	TCP Protocol = 0x02
)

/*
Frame
所有数据都是基于HAPI组成
Header 目前使用固定的
*/
type Frame struct {

	/*
		Header 6字节 这个一般不用填,是固定的值
		ip数据包的信息头，固定为6个字节，这个信息头是任何一种数据包都需要的，后面会经常使用到，6个字节的含义如下

		Byte0: 信息头的长度，目前的版本固定为0x06，即信息头为6个字节长度
		Byte1: knx ip协议版本号, 目前为0x10
		Byte2-3: 服务类型, 表示数据包的用途，这个服务类型并不多，大概10几种类型，对于简单的应用需要的更少，我们只需要了解用到的类型，其他暂时不管。在连接请求数据包中，这个值为0x0205。
		Byte4-5: 整个数据包长度，这个长度要包含信息头的长度
	*/
	//Header HAPI 这个也是HAPI,但是我们只需要Service就能生成,简化用户输入
	Service Service

	HAPIs []HAPI
}

func (this Frame) Bytes() g.Bytes {
	data := []byte(nil)
	length := uint16(6) //头部数据长度
	for _, h := range this.HAPIs {
		length += h.Len()
	}
	data = append(data, NewHeader(this.Service, length)...)
	for _, h := range this.HAPIs {
		data = append(data, h...)
	}
	return data
}

/*
HAPI
HAPI格式字节的含义如下:

Byte0: HAPI数据的长度
Byte1: HAPI的类型
Byte2-n: HAPI数据域

例如网络地址信息的HAPI:
长度(8字节)	类型(udp)	ip地址(4字节)	端口号(2字节)
0x08 		0x01 		192 168 0 1 	8080
*/
type HAPI []byte

func (this HAPI) Len() uint16 {
	return uint16(len(this))
}

func (this HAPI) Type() uint8 {
	if this.Len() < 2 {
		return 0
	}
	return this[1]
}

func (this HAPI) Data() []byte {
	if this.Len() < 2 {
		return []byte{}
	}
	return this[2:]
}

func NewHAPI(Type uint8, data []byte) HAPI {
	h := []byte{byte(len(data) + 2), Type}
	h = append(h, data...)
	return h
}

/*
	NewHAPIAddress
	网络地址 8字节 HAPI ip + port
	这两个部分的都是HAPI格式的，HAPI专门用于描述主机的信息，便于连接双方互相了解对方的一些关键信息。
	HAPI格式8个字节的含义如下:

	Byte0: HAPI数据的长度，固定为0x08
	Byte1: 使用的通讯协议, 目前为0x01，表示UDP
	Byte2-5: ip地址，4个字节
	Byte6-7: 端口号，2个字节
*/
func NewHAPIAddress(proto Protocol, ip net.IP, port uint16) HAPI {
	return NewHAPI(uint8(proto), append(ip.To4(), byte(port>>8), byte(port)))
}

func NewHeader(Service Service, length uint16) HAPI {
	// 0x10是协议版本,目前固定
	return NewHAPI(0x10, append(Service.Bytes(), byte(length>>8), byte(length)))
}

func NewCRI() HAPI {
	return NewHAPI(0x04, []byte{0x02, 0x00})
}

/*



 */

func Decode(bs []byte) (*Frame, error) {
	baseLength := 22
	if len(bs) < baseLength {
		return nil, fmt.Errorf("基础长度错误,预期(%d),得到(%d)", baseLength, len(bs))
	}
	if conv.Uint16(bs[:2]) != Prefix {
		return nil, fmt.Errorf("前缀错误,预期(%x),得到(%x)", Prefix, bs[:2])
	}
	totalLength := conv.Int(bs[4:6])
	if len(bs) != totalLength {
		return nil, fmt.Errorf("数据长度错误: 预期(%d),得到(%d)", len(bs), totalLength)
	}

	f := &Frame{
		Service: Service(conv.Uint16(bs[2:4])),
	}

	bs = bs[6:]
	for len(bs) > 0 {
		length := int(bs[0])
		if len(bs) >= length {
			f.HAPIs = append(f.HAPIs, bs[:length])
			bs = bs[length:]
			continue
		}
		break
	}
	return f, nil
}

func ReadFunc(r *bufio.Reader) ([]byte, error) {
	for {
		result := []byte(nil)
		if _, err := buf.ReadPrefix(r, conv.Bytes(Prefix)); err != nil {
			return nil, err
		}
		result = append(result, conv.Bytes(Prefix)...)

		buf := make([]byte, 4)
		if _, err := io.ReadAtLeast(r, buf, 4); err != nil {
			return nil, err
		}
		result = append(result, buf...)

		length := conv.Int(buf[2:])
		if length >= 22 {
			length -= 6
			buf = make([]byte, length)
			if _, err := io.ReadAtLeast(r, buf, length); err != nil {
				return nil, err
			}
			result = append(result, buf...)
			return result, nil
		}
	}
}
