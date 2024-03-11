package knx

import (
	"bufio"
	"errors"
	"io"
	"net"
	"strconv"
	"strings"
)

// Register 注册,对应服务ConnReqService,
// 该协议需要向网关进行注册操作,
// 注册成功网关会返回channel,是一个一字节的唯一标识
// 需要定时向网关发送数据,以保持连接,类似心跳
func Register(clientAddress, gatewayAddress string, r io.ReadWriter) (uint8, error) {

	//解析网关地址
	ipGateway, portGateway, err := ParseV4(gatewayAddress)
	if err != nil {
		return 0, err
	}

	//解析客户端地址
	ipClient, portClient, err := ParseV4(clientAddress)
	if err != nil {
		return 0, err
	}

	f := &ConnReq{
		Client:  NewHAPIAddress(UDP, ipClient, portClient),
		Gateway: NewHAPIAddress(UDP, ipGateway, portGateway),
		Type:    0x02,
	}

	//写入请求数据
	if _, err := r.Write(f.Bytes()); err != nil {
		return 0, err
	}

	//读取响应数据
	readBytes, err := ReadFunc(bufio.NewReader(r))
	if err != nil {
		return 0, err
	}

	//解析响应数据
	connRes, err := DecodeConnRes(readBytes)
	if err != nil {
		return 0, err
	}

	return connRes.Channel, nil

}

func ReadReq(clientAddress, gatewayAddress string, r io.ReadWriter) *TunnelReq {
	return &TunnelReq{
		Channel:   0,
		SeqNumber: 0,
		//Cemi:      NewHAPI(),
	}
}

// ParseV4 解析ipv4网络地址
func ParseV4(s string) (net.IP, uint16, error) {
	host, portStr, err := net.SplitHostPort(s)
	if err != nil {
		if strings.HasSuffix(err.Error(), "missing port in address") {
			host, portStr, err = net.SplitHostPort(s + ":80")
		}
	}
	if err != nil {
		return nil, 0, err
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return nil, 0, errors.New("无效地址: " + s)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, 0, err
	}
	return ip, uint16(port), nil
}
