package bacnet

import "github.com/injoyai/base/g"

/*

参考文档
https://www.xjx100.cn/news/650081.html?action=onClick

*/

/*
Pkg

+---------------------------+---------------+---------------+
|BACnet/IP					|BACnet网络层	|BACnet应用层	|
+-------+-----------+-------+---------------+---------------+
|Type	|Function	|Length	|NPDU			|APDU			|
+-------+-----------+-------+---------------+---------------+
|1字节	|1字节		|1字节	|				|				|
+-------+-----------+-------+---------------+---------------+

*/
type Pkg struct {
	Function Function
	Control  Control
	DNET     uint16 //最终目标网络号，2个字节。
	DLEN     uint8
	DADR     []byte
	SNET     uint16
	SLEN     uint8
	SADR     []byte
	Count    uint8
	NetType  uint8
	Vendor   uint16
}

func (this *Pkg) Bytes() g.Bytes {
	data := []byte{TypeIP}                                   //BACnet/IP
	data = append(data, uint8(this.Function))                //报文类型,点对点,广播
	data = append(data, 0x01)                                //版本
	data = append(data, this.Control.Byte())                 //控制
	data = append(data, byte(this.DNET>>8), byte(this.DNET)) //最终目标网络号，2个字节。
	data = append(data, this.DLEN)                           //最终目标的MAC层地址的长度，1个字节，0表示对目标网络的广播。
	data = append(data, this.DADR...)                        //最终目标的MAC层地址
	data = append(data, byte(this.SNET>>8), byte(this.SNET)) //初始源网络号，2个字节
	data = append(data, this.SLEN)                           //初始源的MAC层地址的长度，1个字节
	data = append(data, this.SADR...)                        //初始源的MAC层地址
	data = append(data, 0xFF)                                //递减计数器值 防止递归
	return data
}

type Type uint8

const (
	// TypeIP 代表BACnet/IP网络
	TypeIP Type = 0x81
)

type Function uint8

const (
	// FunctionP2P 点对点通讯
	FunctionP2P = 0x0A

	// FunctionBroadcast 广播通信
	FunctionBroadcast = 0x0B
)

type Control struct {

	// HasAPDU
	// 0表示NPDU传送的是一个包含BACnet APDU的数据报文，报文类型域不存在。
	// 1表示NPDU传送的是一个网络层报文，报文类型域存在。
	HasAPDU bool

	// HasCount
	// 0 DNET, DLEN, DADR, Hop Count 不存在
	// 1 DNET, DLEN, Hop Count 存在,DLEN = 0表示广播MAC，DADR不存在,DLEN > 0规定了DADR域的长度
	HasCount bool

	// HasSent
	// 0 SNET, SLEN, SADR 不存在
	// 1 SNET, SLEN, SADR 存在,SLEN = 0无效,SLEN > 0规定了SADR域的长度
	HasSent bool

	// NetResponse
	// 0 存在一个网络层报文期待的应答(除BACnet-Confirmed-Request-PDU，BACnet-ComlexACK-PDU段外)
	// 0 存在一个BACnet-Confirmed-Request-PDU，一个BACnet-ComlexACK-PDU段，或者一个网络层报文期待的应答。
	NetResponse bool

	// Priority 网络优先级
	// 11 关于楼宇安全性的报文。
	// 10 关于楼宇关键设备的报文。
	// 01 紧急报文。
	// 00 一般报文。
	Priority uint8
}

func (this *Control) Byte() (b uint8) {
	if this.HasAPDU {
		b += 1 << 7
	}
	if this.HasCount {
		b += 1 << 5
	}
	if this.HasSent {
		b += 1 << 3
	}
	if this.NetResponse {
		b += 1 << 2
	}
	b += this.Priority % 4
	return b
}
