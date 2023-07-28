package iec104

import (
	"errors"
	"fmt"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"io"
	"math"
	"time"
)

/*
参考
https://blog.csdn.net/wgd0707/article/details/122344581
https://blog.csdn.net/m0_55987469/article/details/130267870
https://blog.csdn.net/chenyitao736866376/article/details/99120024


遥脉: 电度量,是指对现场某装置所发出的脉冲信号进行周期累计的一种远程计数操作
遥信: 数字输入DI 开关量
遥测: 模拟输入AI
遥控: 数字输出DO 开关量
遥调: 模拟输出AO

*/

// NewZHTotal 总召唤
// @slave RTU地址 从站地址
func NewZHTotal(slave uint16) g.Bytes {
	return (&APDU{
		APCI: APCI{},
		ASDU: ASDU{
			Type:   TypeZHTotal,
			Slave:  slave,
			Reason: ReasonStart,
			Info: []Info{
				{
					Addr: [3]byte{},
					QOI:  0x14,
				},
			},
		},
	}).Bytes()
}

// NewSTARTDT_C 启动 U帧
func NewSTARTDT_C() g.Bytes {
	return []byte{0x68, 0x04, STARTDT_C, 0x00, 0x00, 0x00}
}

// NewSTARTDT_A 启动确认 U帧
func NewSTARTDT_A() g.Bytes {
	return []byte{0x68, 0x04, STARTDT_A, 0x00, 0x00, 0x00}
}

// NewSTOP_C 停止 U帧
func NewSTOP_C() g.Bytes {
	return []byte{0x68, 0x04, STOP_C, 0x00, 0x00, 0x00}
}

// NewSTOP_A 停止确认 U帧
func NewSTOP_A() g.Bytes {
	return []byte{0x68, 0x04, STOP_A, 0x00, 0x00, 0x00}
}

// NewTESTFR_C 测试 U帧 心跳的作用
func NewTESTFR_C() g.Bytes {
	return []byte{0x68, 0x04, TESTFR_C, 0x00, 0x00, 0x00}
}

// NewTESTFR_A 测试确认 U帧
func NewTESTFR_A() g.Bytes {
	return []byte{0x68, 0x04, TESTFR_A, 0x00, 0x00, 0x00}
}

// NewAck S格式的确认报文
func NewAck(msgID uint16) g.Bytes {
	return (&APDU{
		APCI: APCI{
			Control1: Order_A,
			Control2: 0,
			Control3: byte(msgID),
			Control4: byte(msgID >> 8),
		},
	}).Bytes()
}

/*




 */

const (
	Prefix = 0x68
)

type APCI struct {
	Control1 byte
	Control2 byte
	Control3 byte
	Control4 byte
}

// WriteNo 发送序号
func (this APCI) WriteNo() uint16 {
	return uint16(this.Control2) + uint16(this.Control1)
}

// ReadNo 接收序号
func (this APCI) ReadNo() uint16 {
	return uint16(this.Control4) + uint16(this.Control3)
}

func (this APCI) Bytes(length int) g.Bytes {
	return []byte{
		Prefix,           //起始字节
		byte(length + 4), //数据总长度
		this.Control1,    //发送序号
		this.Control2,    //发送序号
		this.Control3,    //接收序号
		this.Control4,    //接收序号
	}
}

/*
ASDU
VSQ 可变结构限定词
	SQ
	0: 由信息对象地址寻址的单个信息元素或元素集合,每个信息元素分别带有信息体地址
	1: 单个元素或信息元素同类集合的序列,所有信息共同体用信息体地址,后续信息体地址依次加1
|	D7	|	D6	|	D5	|	D4	|	D3	|	D2	|	D1	|	D0	|
|	SQ	|		Num(信息对象的数量)
*/
type ASDU struct {
	Type   Type   //类型标识
	Reason Reason //传送原因
	Slave  uint16 //公共地址,即RTU站址
	Info   []Info //信息对象
}

func (this ASDU) Bytes() g.Bytes {
	data := []byte(nil)
	data = append(data, byte(this.Type))                         //类型标识,遥信...
	data = append(data, byte(len(this.Info)))                    //可变结构限定词,信息对象数量
	data = append(data, byte(this.Reason>>8), byte(this.Reason)) //传送原因
	data = append(data, byte(this.Slave>>8), byte(this.Slave))   //公共地址,从站地址
	for _, v := range this.Info {
		data = append(data, v.Bytes()...)
	}
	return data
}

// Info 信息对象
type Info struct {
	Addr [3]byte   //信息对象地址, 操作地址
	QOI  byte      //信息元素集, 操作类型
	Time time.Time //信息对象时标(可选)
}

func (this Info) Bytes() g.Bytes {
	data := append(this.Addr[:], this.QOI)
	if !this.Time.IsZero() {
		//60*1000=60000 < 65535
		mill := this.Time.Second()*1000 + int(this.Time.UnixNano()/1e6)%1000
		data = append(data,
			byte(mill),
			byte(mill/256),
			byte(this.Time.Minute()),
			byte(this.Time.Hour()),
			byte(this.Time.Day()),
			byte(this.Time.Month()),
			byte(this.Time.Year()-2000),
		)
	}
	return data
}

type Request = APDU

type APDU struct {
	APCI
	ASDU
}

func (this *APDU) Bytes() g.Bytes {
	asdu := this.ASDU.Bytes()
	data := this.APCI.Bytes(len(asdu)) //APCI,控制信息
	data = append(data, asdu...)       //数据域ASDU
	return data
}

/*



 */

// Handshake 握手
func Handshake(r io.ReadWriter) error {
	active := NewSTARTDT_C()
	if _, err := r.Write(active); err != nil {
		return err
	}
	buf := make([]byte, 6)
	if _, err := io.ReadAtLeast(r, buf, 6); err != nil {
		return err
	}
	want := NewSTARTDT_A().HEX()
	obtain := g.Bytes(buf).HEX()
	if obtain != want {
		return fmt.Errorf("握手错误,预期(%s),得到(%s)", want, obtain)
	}
	return nil
}

func Decode(bs []byte) (a Response, err error) {

	if len(bs) < 6 {
		return a, errors.New("数据长度小于6字节")
	}
	length := int(bs[1]) + 2
	if length < 6 {
		return a, fmt.Errorf("数据长度字节错误%X", bs)
	}
	a.APCI.Control1 = bs[2]
	a.APCI.Control2 = bs[3]
	a.APCI.Control3 = bs[4]
	a.APCI.Control4 = bs[5]

	if length == 6 {
		return
	}

	if len(bs) != length {
		return a, fmt.Errorf("数据长度错误,预期(%d),得到(%d)", length, len(bs))
	}

	/*

		ASDU 解析

	*/

	a.Type = Type(bs[6])
	//可变结构限定词
	//上报数据原因
	a.Reason = Reason(uint16(bs[9])<<8 + uint16(bs[8]))
	//公共地址,RTU地址
	a.Slave = uint16(bs[11])<<8 + uint16(bs[10])
	bs = bs[12:]

	var addr uint32
	for i := 0; len(bs) > 3; i++ { //i < a.ValueLen()
		if !a.Orderly() || i == 0 {
			addr = conv.Uint32([]byte{bs[2], bs[1], bs[0]})
			bs = bs[3:]
		} else {
			addr++
		}

		func() {
			defer g.Recover(nil)
			val := &Value{Addr: addr}
			switch a.Type {
			case TypeYCMemo:
				//带品质描述的遥测 归一化遥测（整型) 值占3字节
				val = &Value{
					Addr:  addr,
					Value: conv.Uint16([]byte{bs[1], bs[0]}),
					Memo:  Memo(bs[2]),
				}
				bs = bs[3:]
			case TypeYXOne:
				//单点遥信,对应数字输入DI,每个遥信占1个字节
				val.Value = bs[1] == 1 //1是合,0是分
				bs = bs[1:]
			case TypeYXTwo:
				//双点遥信, 每个遥信占1个字节
				val.Value = bs[1] == 2 //2是合,1是分
				bs = bs[1:]
			case TypeYMSB12, TypeYMSB12_:
				//带7个字节长时标的电度量，每个电度量占12个字节
				val.Value = math.Float32frombits(conv.Uint32(g.Bytes(bs[:5]).Reverse().Bytes()))
				val.SetTime(bs[5:12])
				bs = bs[12:]
			}
			a.Values = append(a.Values, val)
		}()
	}

	return
}

type Response struct {
	APCI

	Type Type //类型标识

	// Reason 上报实时值(消息体)原因
	Reason Reason

	/*
		VSQ 可变结构限定词
		第0位(左1)表示是否连续地址
		后7位表示上报实时值(信息体)的数量
	*/
	VSQ byte

	// Slave 公共地址,即RTU站址
	Slave uint16

	// Values 实时值    (信息体)
	Values []*Value
}

// Orderly 实时值(信息体)是否连续地址
func (this Response) Orderly() bool {
	return this.VSQ >= (1 << 7)
}

// ValueLen 实时值(信息体)长度
func (this Response) ValueLen() int {
	return int(this.VSQ << 1 >> 1)
}

func (this Response) String() string {
	return fmt.Sprintf(`
发送序号: %d	接收序号: %d
公共地址: %d	上报数量: %d
上报类型: %d(%s)
上报原因: %d(%s)
%s`,
		this.APCI.WriteNo(), this.APCI.ReadNo(),
		this.Slave, this.ValueLen(),
		this.Type, this.Type,
		this.Reason, this.Reason,
		func() string {
			s := ""
			for _, v := range this.Values {
				s += fmt.Sprintf("地址: %d	值: %v", v.Addr, v.Value)
				if !v.Time.IsZero() {
					s += fmt.Sprintf("	时间: %s", v.Time.String())
				}
				s += "\n"
			}
			return s
		}())
}

type Value struct {
	Addr  uint32      //数据地址
	Value interface{} //数据值
	Memo  Memo        //品质描述
	Time  time.Time   //时间
}

func (this *Value) SetTime(bs g.Bytes) {
	switch bs.Len() {
	case 7:
		mill := int(bs[1])<<8 + int(bs[0])
		this.Time = time.Date(
			int(bs[6])+2000,   //年
			time.Month(bs[5]), //月
			int(bs[4]<<3>>3),  //日
			int(bs[3]),        //时
			int(bs[2]),        //分
			mill/1000,         //秒
			(mill%1000)*1e6,   //纳秒
			time.Local,
		)
	}
}
