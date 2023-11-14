package fins

import (
	"errors"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
)

func DecodeICF(b byte) ICF {
	return ICF{
		IsResponse: (b>>6)%2 == 1,
		NoResponse: b%2 == 1,
	}
}

type ICF struct {
	IsResponse bool //是否命令,否则是响应
	NoResponse bool //是否需要响应
}

func (this ICF) Byte() (b byte) {
	b += byte(1 << 7)
	if this.IsResponse {
		//响应内容
		b += byte(1 << 6)
	}
	if this.NoResponse {
		//不需要响应
		b += byte(1 << 0)
	}
	return
}

type Code uint16

const (
	MemoryRead       Code = 0x0101 //读内存区
	MemoryWrite      Code = 0x0102 //写内存区
	MemoryPadding    Code = 0x0103 //填充内容区
	MemoryWriteMulti Code = 0x0104 //写多个内存区
	MemoryCopy       Code = 0x0105 //传输内容区

	ParamRead    Code = 0x0201 //读参数
	ParamWrite   Code = 0x0202 //写参数
	ParamPadding Code = 0x0203 //填充参数

	CPURun  Code = 0x0401 //运行CPU
	CPUStop Code = 0x0402 //停止CPU
	CPURead Code = 0x0601 //读取CPU状态

	ErrClear    Code = 0x2101 //清除错误
	ErrLogRead  Code = 0x2102 //读取错误日志
	ErrLogClear Code = 0x2103 //清除错误日志
)

type Frame struct {
	ICF ICF   //显示帧信息 1a00000b a(0发送,1响应) b(0需要响应,无需响应)
	RSV uint8 //系统保留,固定0
	GCT uint8 //允许的网关数量,固定2
	DNA uint8 //目标网络地址
	DA1 uint8 //目标节点地址,
	DA2 uint8 //目标单位地址
	SNA uint8 //源网络地址
	SA1 uint8 //源节点地址
	SA2 uint8 //源单元地址
	SID uint8 //服务ID
	//MRC uint8 //主请求代码
	//SRC uint8 //子请求代码
	RC Code //请求类型

	Data []byte //请求类型对应的数据
}

func (this *Frame) Bytes() g.Bytes {
	if this.GCT == 0 {
		this.GCT = 2
	}
	data := []byte(nil)
	data = append(data, this.ICF.Byte())                //ICF
	data = append(data, this.RSV)                       //RSV
	data = append(data, this.GCT)                       //GCT
	data = append(data, this.DNA)                       //DNA
	data = append(data, this.DA1)                       //DA1
	data = append(data, this.DA2)                       //DA2
	data = append(data, this.SNA)                       //SNA
	data = append(data, this.SA1)                       //SA1
	data = append(data, this.SA2)                       //SA2
	data = append(data, this.SID)                       //SID
	data = append(data, conv.Bytes(uint16(this.RC))...) //MRC
	data = append(data, this.Data...)                   //请求类型对的数据
	return data
}

func (this *Frame) DecodeReadResult() ([][2]byte, error) {
	if len(this.Data) < 2 || len(this.Data)%2 != 0 {
		return nil, errors.New("数据长度错误")
	}
	list := [][2]byte(nil)
	for i := 2; i < len(this.Data); i += 2 {
		list = append(list, [2]byte{this.Data[i], this.Data[i+1]})
	}
	return list, nil
}

// DecodeWriteResult 解析写入数据结果
func (this *Frame) DecodeWriteResult() error {
	return ErrMap[conv.Uint16(this.Data)]
}

var (
	ErrMap = map[uint16]error{
		0x0000: nil,
		0x0001: errors.New("服务被取消"),
		0x0101: errors.New("本地节点错误"),
		0x0102: errors.New("令牌超时"),
		0x0103: errors.New("重试失败"),
		0x0104: errors.New("发送帧太多"),
		0x0105: errors.New("节点地址错误"),
		0x0106: errors.New("节点地址重复"),
		0x0201: errors.New("目标节点不在网络中"),
		0x0202: errors.New("Unit missing"),
		0x0203: errors.New("Third node missing"),
		0x0204: errors.New("Destination node busy"),
		0x0205: errors.New("Response timeout"),
	}
)

/*



 */

/*
Read
例读取DM区地址100,连续10个地址的数据 82 006400 000A
*/
type Read struct {
	Area    Area   //区域
	Address uint32 //起始2字节地址+1字节位地址 0-7FFF0F
	Length  uint16 //长度,数量
}

func (this Read) Bytes() g.Bytes {
	data := []byte{uint8(this.Area)}
	data = append(data, conv.Bytes(this.Address)[1:]...)
	data = append(data, conv.Bytes(this.Length)...)
	return data
}

type Write struct {
	Area    Area   //区域
	Address uint32 //起始2字节地址+1字节位地址 0-7FFF0F
	Value   []byte //写入的值
}

func (this Write) Bytes() g.Bytes {
	data := []byte{uint8(this.Area)}
	data = append(data, conv.Bytes(this.Address)[1:]...)
	data = append(data, conv.Bytes(uint16(len(this.Value)))...)
	data = append(data, this.Value...)
	return data
}
