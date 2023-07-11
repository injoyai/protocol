package dlt645

import (
	"encoding/hex"
	"errors"
	"github.com/injoyai/base/bytes"
	"github.com/injoyai/conv"
)

type EnPkg struct {
	No       bytes.Entity //表地址,传入正序
	Control  Control      //控制码,枚举
	Password bytes.Entity //密码,可选
	Mark     bytes.Entity //标识符,枚举
	Data     bytes.Entity //数据域
}

func (this EnPkg) HEX() string {
	return hex.EncodeToString(this.Bytes())
}

func (this *EnPkg) Bytes() []byte {
	bs := []byte{0x68}                                                //起始帧
	bs = append(bs, this.No.Reverse()...)                             //表地址
	bs = append(bs, byte(0x68))                                       //中间帧
	bs = append(bs, this.Control.Byte())                              //控制帧
	length := this.Password.Len() + this.Mark.Len() + this.Data.Len() //数据域长度
	bs = append(bs, byte(length))                                     //数据域长度
	bs = append(bs, this.Mark.Reverse().Add0x33()...)                 //标识符
	bs = append(bs, this.Data.Reverse().Add0x33()...)                 //数据域
	bs = append(bs, bytes.Entity(bs).SumByte())                       //校验和
	bs = append(bs, byte(0x16))                                       //结束帧
	return bs
}

type DePkg struct {
	No      string       `json:"no"`
	Control Control      `json:"control"`
	Mark    Mark         `json:"mark"`
	Data    bytes.Entity `json:"data"`
}

// Result 默认结果,根据国标(部分,常用)
func (this *DePkg) Result() (float64, error) {
	decimals := int(MarkMap[this.Mark.HEX()].Decimals)
	return this.Data.Sub0x33ReverseHEXToFloat64(decimals)
}

func Decode(bs bytes.Entity) (*DePkg, error) {

	//去除0xfe , 4个fe是唤醒设备作用
	for len(bs) > 0 && bs[0] == 0xfe {
		bs = bs[1:]
	}

	//校验长度
	if bs.Len() < 12 {
		return nil, errors.New("基本数据长度错误(小于12字节):" + bs.HEX())
	}

	//校验起始帧
	if bs[0] != 0x68 || bs[len(bs)-1] != 0x16 {
		return nil, errors.New("基本数据结构错误(首尾非0x68,0x16):" + bs.HEX())
	}

	//校验控制码
	controlBin := conv.BINStr(bs[8])
	if controlBin[1] != '0' {
		return nil, errors.New("从站应答标示异常:" + controlBin)
	}

	//数据域长度
	length := int(bs[9])

	//校验数据域长度
	if length+12 != bs.Len() {
		return nil, errors.New("数据域长度错误:" + bs.HEX())
	}

	//解析包
	p := &DePkg{
		No:      bs[1:7].Reverse().HEX(), //解析表地址
		Control: Control(controlBin[3:8]),
	}

	// 是否有标识和数据域
	if length > 4 {
		p.Mark = Mark(bs[10:14])
		p.Data = bs[14 : 10+length]
	}

	return p, nil
}
