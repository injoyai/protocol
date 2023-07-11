package ip

import (
	"errors"
	"fmt"
)

// TCP
// |源端口号	|目的端口号 	|序号      		|确认序号 		  |首部长度/标志	|窗口大小	|校验和 	|紧急指针|
// |252 111 |1 187   	|167 206 24 192 |67 218 185 221   |80 24    	|2 0   		|45 48  | 0 0
type TCP struct {
	SrcPort   uint16 //源端口号 16位
	DstPort   uint16 //目的端口号 16位 443={1,487}
	Serial    uint32 //序号
	SerialAck uint32 //确认序号
	HeadLen   uint8  //4位 首部长度 需要乘以4(字节)

	Sign uint8 //6位 标志位
	URG  bool  //紧急指针是否有效
	ASK  bool  //确认序号是否有效
	PSH  bool  //提示接收端应用程序立刻将TCP接收缓冲区当中的数据读走。
	RST  bool  //表示要求对方重新建立连接。我们把携带RST标识的报文称为复位报文段。
	SYN  bool  //表示请求与对方建立连接。我们把携带SYN标识的报文称为同步报文段。
	FIN  bool  //通知对方，本端要关闭了。我们把携带FIN标识的报文称为结束报文段。

	WindowSize uint16 //16位 窗口大小,缓冲区剩余大小,如果为0,则发过来的包会被丢弃,发送端需要减缓发送
	CRC        uint16 //16位校验和
	Urgent     uint16 //16位紧急指针 需要配合URG使用
	Option     []byte //选项,最多40字节
	Data       []byte //数据
}

func (this *TCP) String() string {
	return fmt.Sprintf("协议:TCP  方向:%d>>%d  总长度:%d(字节)  有效数据:%v",
		this.SrcPort, this.DstPort, len(this.Data)+int(this.HeadLen), this.Payload())
}

func (this *TCP) GetSrcPort() string {
	return fmt.Sprintf(":%d", this.SrcPort)
}

func (this *TCP) GetDstPort() string {
	return fmt.Sprintf(":%d", this.DstPort)
}

func (this *TCP) Type() string {
	return "TCP"
}

func (this *TCP) Payload() []byte {
	return this.Data
}

func DecodeTCP(bs []byte) (*TCP, error) {
	if len(bs) < 20 {
		return nil, errors.New("无效协议数据")
	}

	t := &TCP{}
	//第0.1字节 源端口号
	t.SrcPort = uint16(bs[0])<<8 + uint16(bs[1])
	//第2.3字节,目的端口
	t.DstPort = uint16(bs[2])<<8 + uint16(bs[3])
	//第4.5.6.7字节,序号
	t.Serial = uint32(bs[4])<<24 + uint32(bs[5])<<16 + uint32(bs[6])<<8 + uint32(bs[7])
	//第8.9.10.11字节,确认序号
	t.Serial = uint32(bs[8])<<24 + uint32(bs[9])<<16 + uint32(bs[10])<<8 + uint32(bs[11])
	//第12字节,	//首部长度4位,需乘以4(字节)
	t.HeadLen = (bs[12] >> 4) * 4
	//第13字节,标识6位
	t.Sign = bs[13] << 2 >> 2
	t.URG = (bs[13]>>5)%2 == 1
	t.ASK = (bs[13]>>4)%2 == 1
	t.PSH = (bs[13]>>3)%2 == 1
	t.RST = (bs[13]>>2)%2 == 1
	t.SYN = (bs[13]>>1)%2 == 1
	t.FIN = bs[13]%2 == 1
	//第14.15字节,窗口大小
	t.WindowSize = uint16(bs[14])<<8 + uint16(bs[15])
	//第16.17字节,CRC校验和
	t.CRC = uint16(bs[16])<<8 + uint16(bs[17])
	//第18.19字节,紧急指针
	t.Urgent = uint16(bs[18])<<8 + uint16(bs[19])

	if t.HeadLen < 20 {
		return nil, errors.New("首部长度错误")
	}
	//选项,最多40字节
	t.Option = bs[20:t.HeadLen]
	//数据
	t.Data = bs[t.HeadLen:]

	return t, nil
}
