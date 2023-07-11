package s7

import (
	"encoding/hex"
	"errors"
	"github.com/injoyai/conv"

	"strconv"
)

// Area 区块
type Area struct {
	DataType DataType //数据类型 1 bit 2 word 后面的是啥? 3 dint 4 real 5 counter
	Block    Block    //DB块编号，0是M区 1是DB区 区块
	Addr     uint32   //内部地址 偏移量最大为0x07FFFF  M取需要*8   200.1>>>200*8+1
	Size     uint16   //访问数据的字节个数
}

func (this *Area) encode() []byte {
	data := make([]byte, 12)
	data[0] = 0x12                //固定
	data[1] = 0x0a                //剩余长度
	data[2] = 0x10                //寻址模式 any-type(x010)
	data[3] = byte(this.DataType) //数据类型
	if this.DataType == Bit {
		this.Size = 1
	}
	data[4] = byte(this.Size >> 8)
	data[5] = byte(this.Size)
	data[6] = byte(0) //this.Block >> 8)
	data[7] = this.Block.Bytes()[0]
	data[8] = this.Block.Bytes()[1]
	addr := this.Addr // int(this.Addr)*8 + int(this.Addr*10)%10
	data[9] = byte(addr >> 16)
	data[10] = byte(addr >> 8)
	data[11] = byte(addr)
	return data
}

// Param 参数
type Param struct {
	OrderType OrderType //指令类型,4是读 5是写
	Area      Area      //数据块
}

func (this *Param) encode() []byte {
	buf := make([]byte, 14)
	buf[0] = byte(this.OrderType) // 4读 5写
	buf[1] = 1                    // 读取块数,目前固定读1
	copy(buf[2:], this.Area.encode())
	return buf
}

func (this *Param) len() int {
	return 1*8 + 6
}

type Write struct {
	Value []byte //写入内容
}

func (this *Write) encode(Type DataType) []byte {
	if len(this.Value) == 0 {
		return []byte{}
	}
	data := []byte(nil)
	data = append(data, 0)
	length := uint16(len(this.Value))
	if Type == Bit {
		data = append(data, 3)
	} else {
		data = append(data, 4)
		length *= 8
	}
	data = append(data, conv.Bytes(length)...)
	data = append(data, this.Value...)
	return data
}

func (this Write) len() int {
	if len(this.Value) == 0 {
		return 0
	}
	return len(this.Value) + 4
}

type Pkg struct {
	MsgID uint16 //消息id
	Param Param  //参数
	Write Write  //写入数据
}

func (this *Pkg) SetMsgID(msgID uint16) *Pkg {
	this.MsgID = msgID
	return this
}

func (this *Pkg) len() int {
	return 17 + this.Param.len() + this.Write.len()
}

// Encode
// 03 00 00 16 11 E0 00 00 00 01 00 C1 02 10 00 C2 02 03 00 C0 01 0A
func (this *Pkg) Encode() []byte {
	data := []byte{0x03, 0x00}                                          //报文头,2字节
	data = append(data, conv.Bytes(uint16(this.len()))...)              //报文长度2字节
	data = append(data, []byte{0x02, 0xF0, 0x80}...)                    //固定协议标识,3字节
	data = append(data, byte(0x32))                                     //协议id,1字节0x32
	data = append(data, byte(JobRequest))                               //命令类型,1字节,1是发送,2是响应
	data = append(data, []byte{0x00, 0x00}...)                          //保留字段,2字节,固定0x0000
	data = append(data, conv.Bytes(this.MsgID)...)                      //消息id,(小端,不过没啥问题,传啥回啥)
	data = append(data, conv.Bytes(uint16(this.Param.len()))...)        //参数字段长度,2字节
	data = append(data, conv.Bytes(uint16(this.Write.len()))...)        //数据字段长度,2字节
	data = append(data, this.Param.encode()...)                         //参数
	data = append(data, this.Write.encode(this.Param.Area.DataType)...) //数据
	return data
}

func (this *Pkg) HEX() string {
	return hex.EncodeToString(this.Encode())
}

type DePkg struct {
	MsgID     uint16    //消息id
	OrderType OrderType //指令类型
	DataType  DataType  //数据类型
	Result    Result    //执行结果
	Value     []byte    //读取的内容
	Size      int       //读取的长度
}

func Decode(bs []byte) (*DePkg, error) {
	if len(bs) < 17 {
		return nil, errors.New("数据长度小于17" + hex.EncodeToString(bs))
	}
	length := conv.Int(bs[2:4])
	if len(bs) != length {
		return nil, errors.New("数据长度错误:" + hex.EncodeToString(bs))
	}
	p := new(DePkg)
	p.MsgID = conv.Uint16(bs[11:13])

	lenData := conv.Int(bs[13:15])
	lenInside := conv.Int(bs[15:17])
	if len(bs) < 17+2+lenInside+lenData {
		return nil, errors.New("数据长度错误:" + hex.EncodeToString(bs))
	}
	inside := bs[17+lenData : 17+2+lenInside+lenData]
	if len(inside) != lenInside+2 {
		return nil, errors.New("内部数据长度错误:" + hex.EncodeToString(bs))
	}
	p.OrderType = OrderType(inside[0]) //指令类型
	//lenArea := int(inside[1])          //读取区块数
	// 无数据模式
	if len(inside) == 2 {
		p.Result = Success
		return p, nil
	}
	p.Result = Result(inside[2]) //执行结果 0xff 标识成功
	if p.Result != Success {
		return nil, p.Result
	}
	if p.OrderType == ReadVar {
		p.DataType = DataType(inside[3] - 2)                                       //数据类型
		lenValue, err := strconv.ParseInt(hex.EncodeToString(inside[4:6]), 16, 64) //数据长度
		if err != nil {
			return nil, err
		}
		lenValue /= 8
		//if int(lenValue)/8*2 != len(inside[6:]) {
		//	return nil, errors.New("变量值长度错误:" + hex.EncodeToString(inside))
		//}
		p.Value = inside[6:]
	}

	return p, nil
}
