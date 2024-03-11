package bacnet

import (
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"math"
	"net"
)

/*

参考文档
https://www.xjx100.cn/news/650081.html?action=onClick
https://www.docin.com/p-2030266256.html

*/

const (
	// TypeIP 代表BACnet/IP网络
	TypeIP      = 0x81
	DefaultPort = 0xBAC0
)

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
	BVLC
	NPDU
	APDU
}

func (this *Pkg) Bytes() g.Bytes {
	data := []byte(nil)
	npduBytes := this.NPDU.Bytes()
	apduBytes := this.APDU.Bytes()
	bvlcBytes := this.BVLC.Bytes(uint16(len(npduBytes) + len(apduBytes)))
	data = append(data, bvlcBytes...)
	data = append(data, npduBytes...)
	data = append(data, apduBytes...)
	return data
}

/*
BVLC
BVLC Type： 0x81，代表BACnet/IP网络

BVLC Function： 指定报文的类型，何种作用；本应用使用的报文类型分为两种：

                             0A：点对点通讯

                             0B：广播通信

Length： 指定报文的长度，包括BVLC Type、BVLC Function以及本身在内
*/
type BVLC struct {
	IP       net.IP
	Port     uint16
	Function Function
}

func (this *BVLC) Bytes(length uint16) g.Bytes {
	data := make([]byte, 10)
	copy(data[0:4], this.IP.To4())
	if this.Port == 0 {
		this.Port = DefaultPort
	}
	copy(data[4:6], conv.Bytes(this.Port))
	data[6] = TypeIP
	data[7] = uint8(this.Function)
	copy(data[8:10], conv.Bytes(length+10))
	return data
}

type Function uint8

const (

	// FunctionResult BVLC-Result
	// 这个报文提供一种机制,用来确认那些需要确认的BVLL服务请求的结果
	FunctionResult = 0x00

	// FunctionWriteBroadcastDistributionTable 写入广播分布表?
	FunctionWriteBroadcastDistributionTable = 0x01

	// FunctionReadBroadcastDistributionTable 读取广播分布表?
	FunctionReadBroadcastDistributionTable = 0x02

	// FunctionReadBroadcastDistributionTableAck 读取广播分布表确认?
	FunctionReadBroadcastDistributionTableAck = 0x03

	// FunctionForwardedNPDU 转发NPDU
	FunctionForwardedNPDU = 0x04

	// FunctionRegisterForeignDevice 注册外来设备
	FunctionRegisterForeignDevice = 0x05

	// FunctionReadForeignDeviceTable 读取外来设备路由表
	FunctionReadForeignDeviceTable = 0x06

	// FunctionReadForeignDeviceTableAck 读外来设备路由表确认
	FunctionReadForeignDeviceTableAck = 0x07

	// FunctionDeleteForeignDeviceTableEntry 删除外来设备
	FunctionDeleteForeignDeviceTableEntry = 0x08

	// FunctionDistributeBroadcastToNetwork 分布广播到网络?
	FunctionDistributeBroadcastToNetwork = 0x09

	// FunctionOriginalUnicast 单播
	FunctionOriginalUnicast = 0x0A

	// FunctionOriginalBroadcast 广播通信
	FunctionOriginalBroadcast = 0x0B
)

type NPDU struct {
	//Version uint8   //版本号,1
	Control Control //控制
	DNET    uint16  //最终目标网络号，2个字节。 广播0xFFFF
	DLEN    uint8   //最终目标的MAC层地址的长度，1个字节，0表示对目标网络的广播。
	DADR    []byte  //最终目标的MAC层地址。
	SNET    uint16  //初始源网络号，2个字节。
	SLEN    uint8   //初始源的MAC层地址的长度，1个字节。
	SADR    []byte  //初始源的MAC层地址。
	//HopCount   uint8   //递减计数器值，用来防止报文不被循环路由。1个字节，初始化为X‘FF’。当报文通过每个路由器时，其值被至少减一。如果路由器发现该值已为0，则丢弃此报文。
	NetType NetType //如果控制域中的比特7为1，这表示此报文是一个网络层报文，其报文类型域存在。这是个1字节的域，其内容表示报文携带的各种网络层的控制信息。
	Vendor  uint16  //如果控制域的比特7为1和报文类型域的值为X‘80’至X‘FF’时，Vendor ID域存在，生产商可以有2个字节来编码自己的专有网络层报文类型。
}

func (this *NPDU) Bytes() g.Bytes {
	data := []byte{0x01} //版本号，1
	data = append(data, this.Control.Byte())
	if this.Control.HasDNet {
		data = append(data, conv.Bytes(this.DNET)...)
		data = append(data, this.DLEN)
		if this.DLEN > 0 {
			data = append(data, this.DADR...)
		}
	}
	if this.Control.HasSNet {
		data = append(data, conv.Bytes(this.SNET)...)
		data = append(data, this.SLEN)
		if this.SLEN > 0 {
			data = append(data, this.SADR...)
		}
	}
	data = append(data, 0xFF)
	if this.Control.HasAPDU {
		data = append(data, this.NetType.Byte())
	}
	if this.Control.HasAPDU && (this.NetType >= 0x80 && this.NetType <= 0xFF) {
		data = append(data, conv.Bytes(this.Vendor)...)
	}
	return data
}

// Control
// 网络层协议控制信息
type Control struct {

	// HasAPDU
	// 0表示NPDU传送的是一个包含BACnet APDU的数据报文，报文类型域不存在。
	// 1表示NPDU传送的是一个网络层报文，报文类型域存在。
	HasAPDU bool

	// HasCount
	// 0 DNET, DLEN, DADR, Hop Count 不存在
	// 1 DNET, DLEN, Hop Count 存在,DLEN = 0表示广播MAC，DADR不存在,DLEN > 0规定了DADR域的长度
	HasDNet bool

	// HasSent
	// 0 SNET, SLEN, SADR 不存在
	// 1 SNET, SLEN, SADR 存在,SLEN = 0无效,SLEN > 0规定了SADR域的长度
	HasSNet bool

	// NetResponse
	// 0 存在一个网络层报文期待的应答(除BACnet-Confirmed-Request-PDU，BACnet-ComlexACK-PDU段外)
	// 0 存在一个BACnet-Confirmed-Request-PDU，一个BACnet-ComlexACK-PDU段，或者一个网络层报文期待的应答。
	NetResponse bool

	// Priority 网络优先级
	// 0b11 关于楼宇安全性的报文。
	// 0b10 关于楼宇关键设备的报文。
	// 0b01 紧急报文。
	// 0b00 一般报文。
	Priority uint8
}

func (this *Control) Byte() (b uint8) {
	if this.HasAPDU {
		b += 1 << 7
	}
	if this.HasDNet {
		b += 1 << 5
	}
	if this.HasSNet {
		b += 1 << 3
	}
	if this.NetResponse {
		b += 1 << 2
	}
	b += this.Priority & 0x0F
	return b
}

type APDU struct {
	Type     APDUType //4位 7~4 0-0x0F
	Flag     Flag     //3位 3(SEG)(数据是否分片),2(MOR)(是否有后续分片数据),1(SA)(是否分段确认)
	MaxSeg   uint8    //3位 6~4
	MaxResp  uint8    //4位 3~0
	SeqNum   uint8    //SEG有效时生效,分段的序列号
	PropSize uint8    //SEG有效时生效,表示发送方设备的最大发送窗口

	InvokeID       uint8 //作用是讲证实服务的请求与响应联系起来
	ServiceChoice  uint8 //此处表明次报文的作用，详见BACnetConfirmedServiceChoice,作用是对所有BACnet证实服务的类别进行编码,用以标识产生该APDU的证实服务种类
	ServiceRequest []byte
}

func (this *APDU) Bytes() g.Bytes {
	data := []byte(nil)

	//data = append(data, this.Type.Byte()|this.Flag.Byte())
	//data = append(data, ((this.MaxSeg&0x70)<<4)|(this.MaxResp&0x0F))

	switch this.Type {
	case ConfirmedReq:

		data = append(data, this.Type.Byte()|this.Flag.Byte())
		data = append(data, ((this.MaxSeg&0x70)<<4)|(this.MaxResp&0x0F))
		data = append(data, this.InvokeID)
		if this.Flag.SEG {
			data = append(data, this.SeqNum)
			data = append(data, this.PropSize)
		}
		data = append(data, this.ServiceChoice)
		data = append(data, this.ServiceRequest...)

	case UnConfirmedReq:

		data = append(data, this.Type.Byte())
		data = append(data, this.ServiceChoice)
		data = append(data, this.ServiceRequest...)

	case SimpleAck:

		data = append(data, this.Type.Byte())
		data = append(data, this.InvokeID)
		data = append(data, this.ServiceChoice)

	case ComPlexAck:

		data = append(data, this.Type.Byte()|this.Flag.Byte())
		data = append(data, this.InvokeID)
		if this.Flag.SEG {
			data = append(data, this.SeqNum)
			data = append(data, this.PropSize)
		}
		data = append(data, this.ServiceChoice)
		data = append(data, this.ServiceRequest...)

	case SegmentAck:

		data = append(data, this.Type.Byte()|this.Flag.Byte())
		data = append(data, this.InvokeID)
		if this.Flag.SEG {
			data = append(data, this.SeqNum)
			data = append(data, this.PropSize)
		}

	case Error:

		data = append(data, this.Type.Byte())
		data = append(data, this.InvokeID)
		data = append(data, this.ServiceChoice)
		data = append(data, this.ServiceRequest...)

	case Reject:

		data = append(data, this.Type.Byte())
		data = append(data, this.InvokeID)
		data = append(data, this.ServiceRequest...)

	case Abort:

		data = append(data, this.Type.Byte()|this.Flag.Byte())
		data = append(data, this.InvokeID)
		data = append(data, this.ServiceRequest...)

	}

	return data
}

// GetMaxSeg 获取的分片数量
func (this *APDU) GetMaxSeg() uint8 {
	this.MaxSeg &= 0x08
	size := math.Pow(2, float64(this.MaxSeg))
	return uint8(size)
}

func (this *APDU) GetMaxResp() uint16 {
	this.MaxResp &= 0x0F
	switch this.MaxResp {
	case 0:
		return 50
	case 1:
		return 128
	case 2:
		return 206
	case 3:
		return 480
	case 4:
		return 1024
	case 5:
		return 1476
	}
	return 50
}

type APDUType uint8

func (this APDUType) Byte() byte {
	return uint8(this) << 4
}

const (
	ConfirmedReq   APDUType = iota //证实服务种类
	UnConfirmedReq                 //是所有BACnet非证实服务的请求原语产生的APDU数据类型
	SimpleAck                      //类型是对成功请求原语进行响应的APDU类型,这种类型APDU不含有用户数据部分
	ComPlexAck                     //是对成功请求原语进行复杂响应的APDU类型,其用户数据部分包含与确认有关的详细信息
	SegmentAck                     //是在报文分段传输过程中对接收到一个或多个报文进行确认的APDU类型,其作用是通知分段发送方发送下一个或多个分段
	Error                          //是对证实服务中请求原语"负响应Result(-)"的APDU类型,用于传送请求失败的原因
	Reject                         //是对证实服务请求进行拒绝的APDU类型,拒绝请求的原因可能是请求APDU组成错误,也可能是协议规程执行错误,该类型APDU只用于证实服务之中
	Abort                          //是用于终止设备间事务的APUD类型
)

type Flag struct {
	SEG bool //表示证实服务请求所产生的数据是否以分片的方式传输
	MOR bool //当SEG为true时有效,表示传输数据是否有后续分片
	SA  bool //表示发送分段请求服务的设备是否需要分段复杂确认,todo 意思是每个分片都要确认?
	NAK bool
	SRV bool
}

func (this *Flag) Byte() byte {
	var b byte
	if this.SEG {
		b |= 0x08
		if this.MOR {
			b |= 0x04
		}
		if this.SA {
			b |= 0x02
		}
	}
	//当Type是SegmentAck时候,这2个参数生效
	if this.NAK {
		b |= 0x02
	}
	if this.SRV {
		b |= 0x01
	}
	return b
}

type NetType uint8

func (this NetType) Byte() byte {
	return uint8(this)
}

const (
	// TypeWhoIsRouterToNetwork
	// 节点用来确定通达某目标网络的下一个路由器
	// 询问: 谁是到该网络的路由器
	TypeWhoIsRouterToNetwork NetType = iota

	// TypeIAmRouterToNetwork
	// 广播: 我是到该网络的路由器
	TypeIAmRouterToNetwork

	// TypeICouldBeRouterToNetwork
	// 广播: 我是到该网络的路由器
	TypeICouldBeRouterToNetwork

	// TypeRejectMessageToNetwork
	// 拒绝报文发向该网络
	TypeRejectMessageToNetwork

	// TypeBusyToNetwork
	// 到该网络的路由器繁忙
	TypeBusyToNetwork

	// TypeRouterAvailableToNetwork
	// 到该网络的路由器可用
	TypeRouterAvailableToNetwork

	// TypeInitializeRoutingTable
	// 初始化路由表
	TypeInitializeRoutingTable

	// TypeInitializeRoutingTableAck
	// 初始化路由表确认
	TypeInitializeRoutingTableAck

	// TypeEstablishConnectionToNetwork
	// 建立到该网络的连接
	TypeEstablishConnectionToNetwork

	// TypeDisConnectionToNetwork
	// 释放到该网络的连接
	TypeDisConnectionToNetwork

	/*
		0x0A-0x7F: 由ASHRAE保留使用
		0x08-0xFF: 可用于生产商专有报文
	*/

)
