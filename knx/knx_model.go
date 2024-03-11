package knx

import (
	"fmt"
	"github.com/injoyai/base/g"
)

type ConnReq struct {
	Client  HAPI
	Gateway HAPI

	/*
		Type  连接类型

		有3种情况
		0x02: 	establishes a data-link layer tunnel. Send and receive L_Data.* messages.
				建立数据链路层隧道。发送和接收L_Data.*消息。

		0x04: 	establishes a raw tunnel. Send and receive L_Raw.* messages.
				建立一个原始隧道。发送和接收L_Raw.*消息。

		0x80: 	establishes a bus monitor tunnel. Receive L_Busmon.ind messages.
				建立总线监视器隧道。接收L_Busmon.ind消息。

	*/
	Type uint8

	/*
	   CRI
	   连接请求信息,4个字节，这个不做过多解释，先简单使用固定值。0x04, 0x04, 0x02, 0x00

	   连接请求数据包一共26个字节，其中信息头，HAPI在后面还会用到，是通用的，CRI是只出现在请求数据包内的。
	   2个HAPI是可以为空的（都为0），至少我在以前的一个测试中是这样的，可以正常通信，并不受影响。如果有问题请填充正确的IP及端口信息。
	   用UDP将这个数据包发送到网关的3671端口，向网关请求连接，如果通信正常，网关返回请求响应数据包。
	*/
	//CRI HAPI
}

func (this ConnReq) Bytes() g.Bytes {
	return Frame{
		Service: ConnReqService,
		HAPIs:   []HAPI{this.Client, this.Gateway, NewCRI()},
	}.Bytes()
}

// ConnRes 对应服务ConnResService 连接响应服务
type ConnRes struct {
	Channel uint8
	Status  Code
}

func DecodeConnRes(data []byte) (*ConnRes, error) {
	length := 8
	if len(data) != length {
		return nil, fmt.Errorf("响应数据长度错误: 预期(%d),得到(%d)", length, len(data))
	}
	code := Code(data[7])
	return &ConnRes{
		Channel: data[6],
		Status:  code,
	}, code.Err()
}

type ConnStateReq struct {
	Channel uint8
	Client  HAPI
}

func (this ConnStateReq) Bytes() g.Bytes {
	return Frame{
		Service: ConnStateReqService,
		HAPIs:   []HAPI{{this.Channel, 0}, this.Client},
	}.Bytes()
}

type ConnStateRes struct {
	Channel uint8
	Status  Code
}

func DecodeConnStateRes(data []byte) (*ConnStateRes, error) {
	length := 2
	if len(data) != length {
		return nil, fmt.Errorf("响应数据长度错误: 预期(%d),得到(%d)", length, len(data))
	}
	code := Code(data[1])
	return &ConnStateRes{
		Channel: data[0],
		Status:  code,
	}, code.Err()
}

type DiscReq struct {
	Channel uint8
	Client  HAPI
}

func (this DiscReq) Bytes() g.Bytes {
	return Frame{
		Service: DiscReqService,
		HAPIs:   []HAPI{{this.Channel, 0}, this.Client},
	}.Bytes()
}

type DiscRes struct {
	Channel uint8
	Status  Code
}

func DecodeDiscRes(data []byte) (*DiscRes, error) {
	length := 2
	if len(data) != length {
		return nil, fmt.Errorf("响应数据长度错误: 预期(%d),得到(%d)", length, len(data))
	}
	code := Code(data[1])
	return &DiscRes{
		Channel: data[0],
		Status:  code,
	}, code.Err()
}

type TunnelReq struct {
	Channel uint8

	// Sequential number, used to track acknowledgements
	// 序列号，用于跟踪确认
	SeqNumber uint8

	// Data to be tunneled
	Cemi HAPI
}

func (this TunnelReq) Bytes() g.Bytes {
	return Frame{
		Service: TunnelReqService,
		HAPIs:   []HAPI{NewHAPI(this.Channel, []byte{this.SeqNumber, 0}), this.Cemi},
	}.Bytes()
}
