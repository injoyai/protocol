package ip

import (
	"errors"
	"fmt"
	"hash/crc32"
)

/*
	以太网封装(RFG 894)
	目的MAC地址(6字节)+源地址(6字节)+类型(2字节)+数据域(46-15000字节)+CRC(4字节)
	所以数据链路层最小一个包的长度是 64字节 最大包是15018字节(过大会在IP层分包)

	源地址和目的地址是指网卡的硬件地址（也叫MAC地址），长度是48位(6字节)，是在网卡出厂时固化的。
	帧协议类型字段有三种值，分别对应
		IP协议(0x0800),
		ARP协议(0x0806),地址解析协议（Address Resolution Protocol），是根据IP地址获取MAC地址的一个TCP/IP协议。
		RARP协议(0x8035),反向地址转换协议（Reverse Address Resolution Protocol），是根据MAC地址获取IP地址的一个TCP/IP协议。
	帧末尾是CRC校验码。

	当主机A将该MAC帧发送到局域网当中后，局域网当中的所有主机都可以收到这个MAC帧，包括主机A自己。
	主机A收到该MAC帧后，可以对收到的MAC帧进行CRC校验，如果校验失败则说明数据发送过程中产生了碰撞，此时主机A就会执行碰撞避免算法，后续进行MAC帧重发。
	主机B收到该MAC帧后，提取出MAC帧当中的目的地址，发现该目的地址与自己的MAC地址相同，于是在CRC校验成功后就会将有效载荷交付给上层IP层进行进一步处理。
	局域网中的其他主机收到该MAC帧后，也会提取出MAC帧当中的目的地址，但发现该目的地址与自己的MAC地址不匹配，于是就会直接将这个MAC帧丢弃掉。
	也就是说，当底层收到一个MAC帧后，会根据MAC帧当中的目的地址来判断该MAC帧是否是发给自己的，如果是发送给自己的则会再对其进行CRC校验，如果校验成功则会根据该MAC帧的帧协议类型，将该MAC交付给对应的上层协议进行处理。

*/

type Ethernet struct {
	SrcMAC   uint64 //源MAC地址
	DstMAC   uint64 //目的MAC地址
	DataType uint16 //类型(IP协议,ARP协议,RARP协议)
	Data     []byte //数据域
	CRC      uint32 //CRC校验
}

func (this *Ethernet) Bytes() []byte {
	data := make([]byte, 0)
	data = append(data, byte(this.SrcMAC>>40), byte(this.SrcMAC>>32), byte(this.SrcMAC>>24),
		byte(this.SrcMAC>>16), byte(this.SrcMAC>>8), byte(this.SrcMAC))
	data = append(data, byte(this.DstMAC>>40), byte(this.DstMAC>>32), byte(this.DstMAC>>24),
		byte(this.DstMAC>>16), byte(this.DstMAC>>8), byte(this.DstMAC))
	data = append(data, byte(this.DataType>>8), byte(this.DataType))
	data = append(data, this.Data...)
	crc := crc32.ChecksumIEEE(data)
	data = append(data, byte(crc>>24), byte(crc>>16), byte(crc>>8), byte(crc))
	return data
}

func (this *Ethernet) String() string {
	return fmt.Sprintf(`协议:Ethernet  方向:%s>>>%s  总长度:%d  有效数据:%v`, this.GetSrcMAC(), this.GetDstMAC(), len(this.Data)+22, this.Payload())
}

func (this *Ethernet) GetSrcMAC() string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x", uint8(this.SrcMAC>>40), uint8(this.SrcMAC>>32), uint8(this.SrcMAC>>24), uint8(this.SrcMAC>>16), uint8(this.SrcMAC>>8), uint8(this.SrcMAC))
}

func (this *Ethernet) GetDstMAC() string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x", uint8(this.DstMAC>>40), uint8(this.DstMAC>>32), uint8(this.DstMAC>>24), uint8(this.DstMAC>>16), uint8(this.DstMAC>>8), uint8(this.DstMAC))
}

func (this *Ethernet) Type() string {
	return "Ethernet"
}

func (this *Ethernet) Payload() []byte {
	return this.Data
}

func DecodeMac(bs []byte) (*Ethernet, error) {
	if len(bs) < 18 {
		return nil, errors.New("数据长度错误")
	}
	eth := &Ethernet{}
	eth.SrcMAC = uint64(bs[0])<<40 + uint64(bs[1])<<32 + uint64(bs[2])<<24 + uint64(bs[3])<<16 + uint64(bs[4])<<8 + uint64(bs[5])
	eth.DstMAC = uint64(bs[6])<<40 + uint64(bs[7])<<32 + uint64(bs[8])<<24 + uint64(bs[9])<<16 + uint64(bs[10])<<8 + uint64(bs[11])
	eth.DataType = uint16(bs[12])<<8 + uint16(bs[13])
	length := len(bs)
	eth.Data = bs[14 : length-4]
	eth.CRC = uint32(bs[length-4])<<24 + uint32(bs[length-3])<<16 + uint32(bs[length-2])<<8 + uint32(bs[length-1])
	return eth, nil
}
