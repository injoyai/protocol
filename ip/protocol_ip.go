package ip

import (
	"errors"
	"fmt"
)

/*
69 0 0 66 87 70 0 0 30 6 175 197 192 168 10 66 192 168 10 24 0 102 192 11 0 103 56 250 24 158 164 201 80 24 32 0 15 230 0 0 3 0 0 26 2 240 128 50 3 0 0 166 156 0 2 0 5 0 0 4 1 255 3 0 1 1
*/

// IP IP协议组成
type IP struct {
	Version  uint8  //版本(第一字节前4bit) 4表示IPv4
	HeadLen  uint16 //首部长度(第一字节后4bit) 需要乘以4,大于等于5 可以判断是否有选项
	Priority uint8  //优先权(已启用)
	TOS      uint8  //TOS字段(最小延迟,最大吞吐量,最高可靠性,最小成本,4选一)
	TotalLen uint16 //报文总长度(头部+数据域) 2字节

	ID          uint16 //报文标识,如果IP层进行了分片,则所有分片的标识相同
	ShardForbid bool   //禁止分片 1位
	ShardNext   bool   //是否有后续分片 1位
	ShardOffset uint16 //分片偏移 1位

	TTL        uint8  //生存时间 1字节 一般64,经过一次路由-1,如果等于0,则丢弃,防止循环
	UpProtocol uint8  //上层协议 1字节 例如tcp udp
	HeadCRC    uint16 //头部的CRC校验
	SrcIP      uint32 //源IP地址
	DstIP      uint32 //目的IP地址
	Data       []byte //数据域
}

func (this *IP) String() string {
	return fmt.Sprintf(`协议:%s  方向:%s>>>%s  总长度:%d(字节)  有效数据:%v`,
		this.GetVersion(), this.GetSrcIP(), this.GetDstIP(), this.TotalLen, this.Payload())
}

func (this *IP) GetVersion() string {
	if this.Version == 4 {
		return "IPv4"
	}
	return "未知"
}

func (this *IP) GetSrcIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", uint8(this.SrcIP>>24), uint8(this.SrcIP>>16), uint8(this.SrcIP>>8), uint8(this.SrcIP))
}

func (this *IP) GetDstIP() string {
	return fmt.Sprintf("%d.%d.%d.%d", uint8(this.DstIP>>24), uint8(this.DstIP>>16), uint8(this.DstIP>>8), uint8(this.DstIP))
}

func (this *IP) Type() string {
	return this.GetVersion()
}

func (this *IP) Payload() []byte {
	return this.Data
}

// DecodeIP 解析IP协议
func DecodeIP(bs []byte) (*IP, error) {
	if len(bs) < 20 {
		return nil, errors.New("无效协议数据")
	}
	ip := &IP{}

	// 第0字节
	// 版本 4位 4代表IPv4
	ip.Version = bs[0] >> 4
	// 首部长度(不算数据域) 4位 结果需要乘以4
	ip.HeadLen = uint16(bs[0]<<4>>4) << 2

	// 第1字节
	// 优先权(弃用) 3位
	ip.Priority = bs[1] >> 5
	// TOS字段(最小延迟,最大吞吐量,最高可靠性,最小成本,4选一) 4位
	ip.TOS = bs[1] << 3 >> 6
	// 保留字段 1位 必选是0
	_ = bs[1] % 2

	// 第2.3字节
	// 报文总长度(头部+数据域) 2字节
	ip.TotalLen = uint16(bs[2])<<8 + uint16(bs[3])

	// 第4.5字节
	// 报文标识,分片的报文标识相同
	ip.ID = uint16(bs[4])<<8 + uint16(bs[5])

	// 第6.7字节
	// 保留字段 1位
	_ = bs[6] >> 7
	// 禁止分片 1位 如果报文过长,则丢弃
	ip.ShardForbid = bs[6]<<1>>7 == 1
	// 是否有后续分片 1位
	ip.ShardNext = bs[6]<<2>>7 == 1
	// 分片偏移 13位 结果需要乘以8 = 2^8 =65535
	ip.ShardOffset = uint16(bs[6]<<3>>3)<<8 + uint16(bs[7])

	// 第8字节
	// 生存时间 1字节 一般64,经过一次路由-1,如果等于0,则丢弃,防止循环
	ip.TTL = bs[8]

	// 第9字节
	// 上层协议 1字节 例如tcp udp
	ip.UpProtocol = bs[9]

	// 第10.11字节
	// 头部的CRC校验
	ip.HeadCRC = uint16(bs[10])<<8 + uint16(bs[11])

	// 第12.13.14.15字节
	// 源IP地址
	ip.SrcIP = uint32(bs[12])<<24 + uint32(bs[13])<<16 + uint32(bs[14])<<8 + uint32(bs[15])

	// 第16.17.18.19字节
	// 目的IP地址
	ip.DstIP = uint32(bs[16])<<24 + uint32(bs[17])<<16 + uint32(bs[18])<<8 + uint32(bs[19])

	if int(ip.TotalLen) != len(bs) {
		return nil, errors.New("数据长度错误")
	}
	//数据域
	ip.Data = bs[20:]

	return ip, nil
}
