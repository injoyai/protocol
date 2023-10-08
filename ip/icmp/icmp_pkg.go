package icmp

import (
	"github.com/injoyai/base/g"
	"net"
	"time"
)

func Ping(ip string, timeout time.Duration) (time.Duration, error) {
	conn, err := net.DialTimeout("ip:icmp", ip, timeout)
	if err != nil {
		return 0, err
	}
	defer conn.Close()
	t := time.Now()
	if err = conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return 0, err
	}
	_, err = conn.Write([]byte{
		8, 0, 247, 253, 0, 1, 0, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0})
	if err != nil {
		return 0, err
	}
	buf := make([]byte, 65535)
	_, err = conn.Read(buf)
	return time.Since(t), err
}

/*
Pkg
参考
http://www.360doc.com/content/22/1022/11/41428029_1052717848.shtml
Type	Code	描述
0		0		回显应答(Ping应答)
3		0		网络不可达
3		1		主机不可达
3		2		协议不可达
3		3		端口不可达
3		4		需要分片但设置不分片
3		5		源站选路失败
3		6		目的网络未知
3		7		目的主机未知
3		8		源主机被隔离(作废)
3		9		目的网络被强制禁止
3		10		目的主机被强制禁止
3		11		由于服务TOS,网络不可达
3		12		由于服务TOS,主机不可达
3		13		由于过滤,通信被强制禁止
3		14		主机越权
3		15		优先终止生效
4		0		源端被关闭(基本流控制)
5		0		网络重定向
5		1		主机重定向
5		2		TOS和网络重定向
5		3		TOS和主机重定向
8		0		回显请求(Ping请求)
9		0		路由器通告
10		0		路由器请求
11		0		传输生存时间为0
11		1		数据组装时间为0
12		0		错误的IP首部
12		1		缺少必须选项
17		0		地址掩码请求
18		0		地址掩码应答



*/
type Pkg struct {
	Type Type  //类型
	Code uint8 //代码
}

func (this *Pkg) Bytes() g.Bytes {
	data := []byte(nil)
	data = append(data, uint8(this.Type))
	data = append(data, this.Code)
	data = append(data, 0, 0)
	return data
}

type Type uint8

func (this Type) string() string {
	switch this {
	case 0:
		return "回显应答"
	case 3:
		return "目标不可达"
	case 4:
		return "源端中止"
	case 5:
		return "重定向"
	case 8:
		return "回显请求"
	case 11:
		return "超时"
	}
	return "未知"
}
