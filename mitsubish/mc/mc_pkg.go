package mc

import (
	"bufio"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/injoyai/base/bytes"
)

/*

允许外部设备读写PLC内部寄存器

//参考
https://blog.csdn.net/wy749929317/article/details/124144389

*/

// Block 数据块
type Block uint8

func (this Block) Byte() byte { return byte(this) }

const (

	// 位

	BlockSM Block = 0x91 //特殊继电器 0x00~0x1023
	BlockX  Block = 0x9C //输入		0x00~0x07ff
	BlockY  Block = 0x9D //输出		0x00~0x07ff
	BlockM  Block = 0x90 //内部继电器	0x00~0x8191
	BlockL  Block = 0x92 //锁存继电器	0x00~0x2027
	BlockF  Block = 0x93 //报警器	0x00~0x1023
	BlockV  Block = 0x94 //变址继电器	0x00~0x1023
	BlockB  Block = 0xA0 //链接继电器	0x00~0x07ff

	// 字节

	BlockSD Block = 0xA9 //特殊寄存器	0x00~0x1023
	BlockD  Block = 0xA8 //数据寄存器 0x00~0x011135
	BlockW  Block = 0xB4 //链接寄存器	0x00~0x07ff

)

type Order uint16

func (this Order) Len() int { return 2 }

func (this Order) Bytes() []byte {
	return []byte{byte(this), byte(this >> 8)}
}

const (
	OrderMoreReadByte  Order = 0x0406 //多个块批量读取
	OrderMoreWriteByte Order = 0x1406 //多个块批量写入
	OrderRead          Order = 0x0401 //批量读取位
	OrderWrite         Order = 0x1401 //批量写入
	OrderReadRand      Order = 0x1402 //随机读取
	OrderMonitorData   Order = 0x0801 //监视数据登录
	OrderMonitor       Order = 0x0802 //监视
)

var (
	pkgReqStart = []byte{0x50, 0x00}
	pkgResStart = []byte{0xD0, 0x00}
	pkgResEnd   = []byte{0x00, 0x00}
)

type Pkg struct {
	Addr  uint32       //地址3字节(软元件起始)
	Block Block        //数据块(软元件)
	Order Order        //指令 读,写...
	Point uint16       //点数,1点数等于2字节
	Data  bytes.Entity //写入的数据
}

func (this *Pkg) Len() int {
	return len(this.Data) + 12
}

func (this *Pkg) LenBytes() []byte {
	length := len(this.Data) + 12
	return []byte{byte(length), byte(length >> 8)}
}

func (this *Pkg) PointBytes() []byte {
	return []byte{byte(this.Point), byte(this.Point >> 8)}
}

func (this *Pkg) AddrBytes() []byte {
	return []byte{byte(this.Addr), byte(this.Addr >> 8), byte(this.Addr >> 16)}
}

/*
Bytes
指令为向软元件D7000写入值H000C
50 00 00 FF FF 03 00 0E 00 10 00 01 14 00 00 58 1B 00 A8 01 00 0C 00
读取软元件D7000开始的连续5个软元件的值，其中0C是上一次写入的数
50 00 00 FF FF 03 00 0C 00 10 00 01 04 00 00 58 1B 00 A8 05 00
*/
func (this *Pkg) Bytes() bytes.Entity {
	data := []byte(nil)
	data = append(data, pkgReqStart...)         //帧头
	data = append(data, 0x00)                   //网络编号
	data = append(data, 0xff)                   //可编程控制器网络编号
	data = append(data, 0xff, 0x03)             //请求目标模块I/O编号
	data = append(data, 0x00)                   //请求目标模块站号
	data = append(data, this.LenBytes()...)     //请求数据长度
	data = append(data, 0x10, 0x00)             //CPU监视定时器
	data = append(data, this.Order.Bytes()...)  //指令批量写入
	data = append(data, 0x00, 0x00)             //子指令
	data = append(data, this.AddrBytes()...)    //起始软元件
	data = append(data, this.Block.Byte())      //软元件代码
	data = append(data, this.PointBytes()...)   //软元件点数
	data = append(data, this.Data.Reverse()...) //软元件点数的数据
	return data
}

/*
Decode
指令为向软元件D7000写入值H000C响应
D0 00 00 FF FF 03 00 02 00 00 00
读取软元件D7000开始的连续5个软元件的值，其中0C是上一次写入的数
D0 00 00 FF FF 03 00 0C 00 00 00 0C 00 00 00 00 00 00 00 00 00
*/
func Decode(bs []byte) (*Pkg, error) {

	//校验基础长度
	if len(bs) < 11 {
		return nil, errors.New("基础长度错误:小于8")
	}

	//校验帧头
	if hex.EncodeToString(bs[:2]) != hex.EncodeToString(pkgResStart) {
		return nil, fmt.Errorf("帧头错误,预期(%X),得到(%X)", pkgResStart, bs[:2])
	}

	//网络编号啥的
	_ = bs[2:7]

	//后续数据长度
	length := int(bs[8])<<8 + int(bs[7])

	//校验长度
	if len(bs) != 9+length {
		return nil, fmt.Errorf("数据长度错误,预期(%d),得到(%d)", 9+length, len(bs))
	}

	//校验结束帧
	if bs[9] != 0x00 || bs[10] != 0x00 { //bs[len(bs)-1] != 0x00 || bs[len(bs)-2] != 0x00 {
		return nil, fmt.Errorf("结束帧错误,预期(%X),得到(%X)", pkgResEnd, bs[9:11])
	}

	p := &Pkg{}

	//转成大端
	for i := 11; i < 11+length-2; i += 2 {
		p.Data = append(p.Data, bs[i+1], bs[i])
	}

	return p, nil
}

type Result uint16

func (this Result) Error() string {
	switch true {
	case this >= 0x4000 && this < 0x4fff:
		return "CPU模块检测出的错误"
	case this >= 0xC051 && this < 0xC054:
		return "超出点数允许范围"
	default:
		switch this {
		case 0xC050:
			return "不允许写入"
		case 0xC055:
			return "无法转为二进制"
		case 0xC056:
			return "超出最大地址范围"
		case 0xC058:
			return "ASCII>二进制转换之后,字符部分数据数不一致"
		case 0xC059:
			return "指令错误,无法识别"
		case 0xC05B:
			return "CPU无法对指定软元件进行写入及读取"
		case 0xC05C:
			return "请求内容有错误"
		case 0xC05D:
			return "未进行监视登录"
		case 0xC05F:
			return "无法对对象CPU模块执行请求"
		case 0xC060:
			return "请求内容中有错误"
		case 0xC061:
			return "请求数据部分长度与字符部分长度不一致"
		case 0xC06F:
			return "通讯数据代码与实际不一致"
		case 0xC070:
			return "无法对对象站点进行软元件存储器的拓展指定"
		case 0xC0B5:
			return "指定了CPU模块中无法处理的数据"
		case 0xC200:
			return "远程口令中有错误"
		case 0xC201:
			return "通讯中使用的端口处于远程口令的锁定状态"
		case 0xC204:
			return "与镜像了远程口令解锁处理请求的对象设备不相符"
		}
	}
	return "未知错误"
}

func ReadFunc(c *bufio.Reader) (result []byte, err error) {
	//D00000FF03
	start := []byte{0xD0, 0x00, 0x00, 0xff, 0xff, 0x03}
	startLen := len(start)
	idx := 0
	for {
		b, err := c.ReadByte()
		if err != nil {
			return nil, err
		}
		if len(start) > idx && b == start[idx] {
			result = append(result, b)
			idx++
			if len(start) == idx {
				//头部信息读取完成

				for i := 0; i < 3; i++ {
					b, err := c.ReadByte()
					if err != nil {
						return nil, err
					}
					result = append(result, b)
				}

				length := int(result[startLen+2])<<8 + int(result[startLen+1])
				for i := 0; i < length; i++ {
					b, err := c.ReadByte()
					if err != nil {
						return nil, err
					}
					result = append(result, b)
				}
				return result, nil
			}
			continue
		}
		//重新
		idx = 0
		result = []byte{}
	}
}
