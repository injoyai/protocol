package dnp3

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/injoyai/base/bytes/crypt/crc"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"github.com/injoyai/io"
)

/*
参考
https://cn-sec.com/archives/933620.html
https://www.yii666.com/article/544647.html
https://blog.csdn.net/weixin_43047908/article/details/119206553

*/

const (
	Prefix uint16 = 0x0564
)

/*
Pkg


例 读
05 64 14 c4 00 04 01 00 e9 b6 c1 c1 01 3c 02 06 3c 03 06 3c 04 06 3c 01 06 62 01

例 05 64 12 c4 03 00 04 00 15 2d c1 c1 02 32 01 07 01 fa 7d 0b 46 0d 01 c8 63
例 05 64 1d 44 02 00 01 00 ef 56 c4 c1 91 90 00 46 04 5b 01 0d 00 00 00 00 00 00 21 e5 00 00 00 00 09 04 00 04 a0 77
*/
type Pkg struct {
	Header Header
	Body   Body
}

// Bytes 0564 12 C4 0300 0400 152D
func (this *Pkg) Bytes() g.Bytes {
	data := []byte(nil)
	bodyBytes := this.Body.Bytes()
	data = append(data, this.Header.Bytes(uint8(len(bodyBytes)+3))...)
	data = append(data, bodyBytes...)
	return data
}

//===============================Header===============================

// Header 链路层,10字节=报文头(2)+报文长度(1)+控制码(1)+源地址(2)+目的地(2)+CRC校验(2)
type Header struct {
	Control LinkControl //控制码
	From    uint16      //源地址
	To      uint16      //目的地
}

func (this Header) CRC16(length uint8) uint16 {
	data := conv.Bytes(Prefix)
	data = append(data, length)
	data = append(data, this.Control.Byte())
	data = append(data, byte(this.To), byte(this.To>>8))
	data = append(data, byte(this.From), byte(this.From>>8))
	return conv.Uint16(crc.Encrypt16(data, crc.CRC16_DNP).Reverse().Bytes())
}

func (this Header) Bytes(length uint8) g.Bytes {
	data := conv.Bytes(Prefix)
	data = append(data, length)                                          //长度
	data = append(data, this.Control.Byte())                             //控制码
	data = append(data, byte(this.To), byte(this.To>>8))                 //目的地
	data = append(data, byte(this.From), byte(this.From>>8))             //源地址
	data = append(data, crc.Encrypt16(data, crc.CRC16_DNP).Reverse()...) //crc校验
	return data
}

type LinkControl struct {
	ToSlave  bool           //发送方向 第一位
	IsMaster bool           //是否是主设备 第二位
	Correct  bool           //主站标识纠错 第三位
	Function HeaderFunction //功能码,后4位
}

func (this LinkControl) Byte() byte {
	b := byte(0)
	if this.ToSlave {
		b |= 0x80
	}
	if this.IsMaster {
		b |= 0x40
	}
	if this.IsMaster && this.Correct {
		b |= 0x20
	}
	b |= this.Function.Byte()
	return b
}

// HeaderFunction 后4位是功能码
type HeaderFunction uint8

func (this HeaderFunction) Byte() byte {
	return byte(this) & 0x0F
}

const (
	LinkReset     HeaderFunction = 0  //链路重置 ,从站地址:同意
	ProcessReset  HeaderFunction = 1  //进程重置 ,从站地址:拒绝
	RequestSend   HeaderFunction = 3  //请求发送数据
	Send          HeaderFunction = 4  //发送数据
	LinkStatus    HeaderFunction = 9  //查询当前链路的状态
	LinkStatusAck HeaderFunction = 11 //回应当前链路状态
)

//===============================Transfer===============================

// PkgNo 包序,从1开始,最大值63,不填写表示只有一个包
type PkgNo struct {
	IsLast  bool  //是否是最后一包
	IsFirst bool  //是否是第一包
	Current uint8 //当前,最大值64
}

func (this PkgNo) Byte() byte {
	b := byte(0)
	if this.IsLast {
		//表示最后一个包
		b |= 0x80
	}
	if this.IsFirst {
		//表示是第一个包
		b |= 0x40
	}
	b |= this.Current & 0x3F
	if b == 0 {
		b = 0xC1
	}
	return b
}

//===============================Body===============================

type BodyFunction uint8

func (this BodyFunction) Byte() byte {
	return byte(this)
}

const (
	Confirm                   BodyFunction = 0x00 //确认
	Read                      BodyFunction = 0x01 //读数据
	Write                     BodyFunction = 0x02 //写数据
	Select                    BodyFunction = 0x03 //选择
	Operate                   BodyFunction = 0x04 //执行
	DirectOperate             BodyFunction = 0x05 //直接执行
	DirectOperateNoAck        BodyFunction = 0x06 //直接执行，不需要响应
	ImmediateFreeze           BodyFunction = 0x07 //立即冻结
	ImmediateFreezeNoAck      BodyFunction = 0x08 //立即冻结，不需要响应
	FreezeClear               BodyFunction = 0x09 //冻结清除
	FreezeClearNoAck          BodyFunction = 0x0A //冻结清除，不需要响应
	FreezeWithTime            BodyFunction = 0x0B //带时间冻结
	FreezeWithTimeNoAck       BodyFunction = 0x0C //带时间冻结，不需要响应
	ColdRestart               BodyFunction = 0x0D //冷启动
	WarmRestart               BodyFunction = 0x0E //热启动
	InitDataToDefault         BodyFunction = 0x0F //初始化数据到默认
	InitAPP                   BodyFunction = 0x10 //初始化应用
	StartAPP                  BodyFunction = 0x11 //启动应用
	StopAPP                   BodyFunction = 0x12 //停止应用
	SaveConfig                BodyFunction = 0x13 //保存配置
	EnableUnsolicitedMessage  BodyFunction = 0x14 //自发报文使能
	DisableUnsolicitedMessage BodyFunction = 0x15 //自发报文禁止
	AssignClass               BodyFunction = 0x16 //分配类
	DelayMeasurement          BodyFunction = 0x17 //延时测量
	Response                  BodyFunction = 0x81 //响应
	UnsolicitedMessage        BodyFunction = 0x82 //自发报文
)

type Body struct {
	PkgNo    PkgNo        //包序,这独立,是总体的包序
	Control  BodyControl  //应用控制,分包情况
	Function BodyFunction //功能码 例如 0x01(读) 0x02(写)
	Datas    []Data       //数据
}

func (this Body) Bytes() g.Bytes {
	data := []byte(nil)
	data = append(data, this.PkgNo.Byte())
	data = append(data, this.Control.Byte())
	data = append(data, this.Function.Byte())
	for _, v := range this.Datas {
		data = append(data, v.Bytes()...)
	}
	data = append(data, crc.Encrypt16(data, crc.CRC16_DNP).Reverse()...)
	return data
}

type Data struct {
	DataType  DataType      //数据类型,例如时间类型,模拟输入等
	Qualifier BodyQualifier //限定词
	Range     []byte        //范围(变程)
	Data      []byte        //数据
}

func (this Data) Bytes() []byte {
	data := []byte(nil)
	data = append(data, this.DataType.Bytes()...)
	data = append(data, this.Qualifier.Byte())
	switch this.Qualifier.Code {
	case 0, 1, 2, 3, 4, 5:
		//当限定词码取值 0～5时，变程段包含1个开始范围 (Start Range) 和1个结束范围 (Stop range)
		data = append(data, this.Range...)
	case 6:

	case 7, 8, 9:
		//当限定词码取值为 7～ 9 时，则变程段由一个计数值所组成，它指明所讨论的数据对象的数目
		data = append(data, this.Range...)
	default:
	}
	data = append(data, this.Data...)
	return data
}

type BodyControl struct {
	IsFirst     bool  //是否是报文第一分片
	IsFinlay    bool  //是否是报文最后一分片
	NeedAck     bool  //是否需要应答
	Unsolicited bool  //是否主动发起
	No          uint8 //分片序号，从1开始，最大值15
}

func (this BodyControl) Byte() byte {
	b := uint8(0)
	if this.IsFirst {
		b |= 0x80
	}
	if this.IsFinlay {
		b |= 0x40
	}
	if this.NeedAck {
		b |= 0x20
	}
	if this.Unsolicited {
		b |= 0x10
	}
	b |= this.No & 0x1F
	if b == 0 {
		b = 0xC1
	}
	return b
}

/*
BodyQualifier
0: 预留
1~3: 索引规模,规定前置于每个数据对象的索引规模或对象的规模
	 在请求报文中，当限定词码 (Qualifier Code) 等于 11时,1、2、3分别代表数据对象前的索引是1、2、4个字节。0无效。 4、5、6、7保留。
4-7: 限定词码,用以规定变程 (Range) 意义。当限定词码取值 0～5时，变程段包含1个开始范围 (Start Range) 和1个结束范围 (Stop range)。当限定词码取值 6时，则 Range 段的长为零 (即无变程段 )，因为所指定的是所要求的数据类型的全部数据对象。当限定词码取值为 7～ 9 时，则变程段由一个计数值所组成，它指明所讨论的数据对象的数目。
*/
type BodyQualifier struct {
	Index uint8 //3位,1~3
	Code  uint8 //4位,4~7
}

func (this BodyQualifier) Byte() byte {
	b := uint8(0)
	b |= this.Index << 5 >> 1
	b |= this.Code << 4 >> 4
	return b
}

// DataType 数据类型
type DataType uint32

func (this DataType) Bytes() []byte {
	return []byte{byte(this >> 8), byte(this)}
}

func (this DataType) Name() string {
	switch this {
	case InputBit:
		return "1bit数字输入"
	case InputStatus:
		return "状态数字输入"
	case InputDisplacement:
		return "不带时间的数字输入变位"
	case InputDisplacementTime:
		return "带时间的数字输入变位"
	case InputDisplacementRelativeTime:
		return "带相对时间的数字输入变位"
	case OutputBit:
		return "1bit数字输出(1位)"
	case OutputStatus:
		return "状态数字输出(8位)"
	case OutputRelay:
		return "继电器输出块"
	case OutputPattern:
		return "方式(pattern)控制块"
	case InputCount32:
		return "32位计数器值"
	case InputCount16:
		return "16位计数器值"
	case InputCountSub32:
		return "32位计数器差值"
	case InputCountSub16:
		return "16位计数器差值"
	case InputCountNoStatus32:
		return "32位无状态计数器值"
	case InputCountNoStatus16:
		return "16位无状态计数器值"
	case InputCountNoStatusSub32:
		return "32位无状态计数器差值"
	case InputCountNoStatusSub16:
		return "16位无状态计数器差值"
	case InputCountFreeze32:
		return "32位冻结计数器值"
	case InputCountFreeze16:
		return "16位冻结计数器值"
	case InputCountFreezeSub32:
		return "32位冻结计数器差值"
	case InputCountFreezeSub16:
		return "16位冻结计数器差值"
	case InputCountFreezeTime32:
		return "32位带时间冻结计数器值"
	case InputCountFreezeTime16:
		return "16位带时间冻结计数器值"
	case InputCountFreezeNoStatus32:
		return "32位无状态冻结计数器值"
	case InputCountFreezeNoStatus16:
		return "16位无状态冻结计数器值"
	case InputCountFreezeNoStatusSub32:
		return "32位无状态冻结计数器差值"
	case InputCountFreezeNoStatusSub16:
		return "16位无状态冻结计数器差值"
	case InputCountLimitNoTime32:
		return "32位无时间计数器值越限事件"
	case InputCountLimitNoTime16:
		return "16位无时间计数器值越限事件"
	case InputCountLimitNoTimeSub32:
		return "32位无时间计数器差值越限事件"
	case InputCountLimitNoTimeSub16:
		return "16位无时间计数器差值越限事件"
	case InputCountLimitTime32:
		return "32位带时间计数器值越限事件"
	case InputCountLimitTime16:
		return "16位带时间计数器值越限事件"
	case InputCountLimitTimeSub32:
		return "32位带时间计数器差值越限事件"
	case InputCountLimitTimeSub16:
		return "16位带时间计数器差值越限事件"
	case InputCountLimitFreezeNoTime32:
		return "32位无时间冻结计数器值越限事件"
	case InputCountLimitFreezeNoTime16:
		return "16位无时间冻结计数器值越限事件"
	case InputCountLimitFreezeNoTimeSub32:
		return "32位无时间冻结计数器差值越限事件"
	case InputCountLimitFreezeNoTimeSub16:
		return "16位无时间冻结计数器差值越限事件"
	case InputCountLimitFreezeTime32:
		return "32位带时间冻结计数器值越限事件"
	case InputCountLimitFreezeTime16:
		return "16位带时间冻结计数器值越限事件"
	case InputCountLimitFreezeTimeSub32:
		return "32位带时间冻结计数器差值越限事件"
	case InputCountLimitFreezeTimeSub16:
		return "16位带时间冻结计数器差值越限事件"
	case InputAnalog32:
		return "32位模拟量输入值"
	case InputAnalog16:
		return "16位模拟量输入值"
	case InputAnalogNoStatus32:
		return "32位无状态模拟量输入值"
	case InputAnalogNoStatus16:
		return "16位无状态模拟量输入值"
	case InputAnalogFreeze32:
		return "32位冻结模拟量输入值"
	case InputAnalogFreeze16:
		return "16位冻结模拟量输入值"
	case InputAnalogFreezeTime32:
		return "32位带时间冻结模拟量输入值"
	case InputAnalogFreezeTime16:
		return "16位带时间冻结模拟量输入值"
	case InputAnalogFreezeNoStatus32:
		return "32位无状态冻结模拟量输入值"
	case InputAnalogFreezeNoStatus16:
		return "16位无状态冻结模拟量输入值"
	case InputAnalogLimitNoTime32:
		return "32位无时间模拟量输入值越限事件"
	case InputAnalogLimitNoTime16:
		return "16位无时间模拟量输入值越限事件"
	case InputAnalogLimitTime32:
		return "32位带时间模拟量输入越限事件"
	case InputAnalogLimitTime16:
		return "16位带时间模拟量输入越限事件"
	case InputAnalogLimitFreezeNoTime32:
		return "32位无时间冻结模拟量输入值越限事件"
	case InputAnalogLimitFreezeNoTime16:
		return "16位无时间冻结模拟量输入值越限事件"
	case InputAnalogLimitFreezeTime32:
		return "32位带时间冻结模拟量输入值越限事件"
	case InputAnalogLimitFreezeTime16:
		return "16位带时间冻结模拟量输入值越限事件"
	case OutputAnalog32:
		return "32位模拟量输出"
	case OutputAnalog16:
		return "16位模拟量输出"
	case OutputAnalogBlock32:
		return "32位模拟量输出块(block)"
	case OutputAnalogBlock16:
		return "16位模拟量输出块(block)"
	case DateTime:
		return "日期和时间"
	case DateTimeSustain:
		return "持续日期和时间"
	case TimeDelayAlmost:
		return "近似延时(16位,秒)"
	case TimeDelayAccurate:
		return "精确延时(16位,秒)"
	case Class0:
		return "Class0类数据,所有非1,2,3类数据"
	case Class1:
		return "Class1类数据(通常为某组信息对象的变化)"
	case Class2:
		return "Class2类数据(通常为某组信息对象的变化)"
	case Class3:
		return "Class3类数据(通常为某组信息对象的变化)"
	}
	return ""
}

var (
	AllDataType = []DataType{
		InputBit,
		InputStatus,
		InputDisplacement,
		InputDisplacementTime,
		InputDisplacementRelativeTime,
		OutputBit,
		OutputStatus,
		OutputRelay,
		OutputPattern,
		InputCount32,
		InputCount16,
		InputCountSub32,
		InputCountSub16,
		InputCountNoStatus32,
		InputCountNoStatus16,
		InputCountNoStatusSub32,
		InputCountNoStatusSub16,
		InputCountFreeze32,
		InputCountFreeze16,
		InputCountFreezeSub32,
		InputCountFreezeSub16,
		InputCountFreezeTime32,
		InputCountFreezeTime16,
		InputCountFreezeNoStatus32,
		InputCountFreezeNoStatus16,
		InputCountFreezeNoStatusSub32,
		InputCountFreezeNoStatusSub16,
		InputCountLimitNoTime32,
		InputCountLimitNoTime16,
		InputCountLimitNoTimeSub32,
		InputCountLimitNoTimeSub16,
		InputCountLimitTime32,
		InputCountLimitTime16,
		InputCountLimitTimeSub32,
		InputCountLimitTimeSub16,
		InputCountLimitFreezeNoTime32,
		InputCountLimitFreezeNoTime16,
		InputCountLimitFreezeNoTimeSub32,
		InputCountLimitFreezeNoTimeSub16,
		InputCountLimitFreezeTime32,
		InputCountLimitFreezeTime16,
		InputCountLimitFreezeTimeSub32,
		InputCountLimitFreezeTimeSub16,
		InputAnalog32,
		InputAnalog16,
		InputAnalogNoStatus32,
		InputAnalogNoStatus16,
		InputAnalogFreeze32,
		InputAnalogFreeze16,
		InputAnalogFreezeTime32,
		InputAnalogFreezeTime16,
		InputAnalogFreezeNoStatus32,
		InputAnalogFreezeNoStatus16,
		InputAnalogLimitNoTime32,
		InputAnalogLimitNoTime16,
		InputAnalogLimitTime32,
		InputAnalogLimitTime16,
		InputAnalogLimitFreezeNoTime32,
		InputAnalogLimitFreezeNoTime16,
		InputAnalogLimitFreezeTime32,
		InputAnalogLimitFreezeTime16,
		OutputAnalog32,
		OutputAnalog16,
		OutputAnalogBlock32,
		OutputAnalogBlock16,
		DateTime,
		DateTimeSustain,
		TimeDelayAlmost,
		TimeDelayAccurate,
		Class0,
		Class1,
		Class2,
		Class3,
	}
)

const (
	InputBit                      DataType = 0x0101 //1bit数字输入
	InputStatus                   DataType = 0x0102 //状态数字输入
	InputDisplacement             DataType = 0x0201 //不带时间的数字输入变位
	InputDisplacementTime         DataType = 0x0202 //带时间的数字输入变位
	InputDisplacementRelativeTime DataType = 0x0203 //带相对时间的数字输入变位

	OutputBit     DataType = 0x0A01 //1bit数字输出(1位)
	OutputStatus  DataType = 0x0A02 //状态数字输出(8位)
	OutputRelay   DataType = 0x0C01 //继电器输出块
	OutputPattern DataType = 0x0C02 //方式(pattern)控制块

	InputCount32                     DataType = 0x1401 //32位计数器值
	InputCount16                     DataType = 0x1402 //16位计数器值
	InputCountSub32                  DataType = 0x1403 //32位计数器差值
	InputCountSub16                  DataType = 0x1404 //16位计数器差值
	InputCountNoStatus32             DataType = 0x1405 //32位无状态计数器值
	InputCountNoStatus16             DataType = 0x1406 //16位无状态计数器值
	InputCountNoStatusSub32          DataType = 0x1407 //32位无状态计数器差值
	InputCountNoStatusSub16          DataType = 0x1408 //16位无状态计数器差值
	InputCountFreeze32               DataType = 0x1501 //32位冻结计数器值
	InputCountFreeze16               DataType = 0x1502 //16位冻结计数器值
	InputCountFreezeSub32            DataType = 0x1503 //32位冻结计数器差值
	InputCountFreezeSub16            DataType = 0x1504 //16位冻结计数器差值
	InputCountFreezeTime32           DataType = 0x1505 //32位带时间冻结计数器值
	InputCountFreezeTime16           DataType = 0x1506 //16位带时间冻结计数器值
	InputCountFreezeNoStatus32       DataType = 0x1509 //32位无状态冻结计数器值
	InputCountFreezeNoStatus16       DataType = 0x1510 //16位无状态冻结计数器值
	InputCountFreezeNoStatusSub32    DataType = 0x1511 //32位无状态冻结计数器差值
	InputCountFreezeNoStatusSub16    DataType = 0x1512 //16位无状态冻结计数器差值
	InputCountLimitNoTime32          DataType = 0x1601 //32位无时间计数器值越限事件
	InputCountLimitNoTime16          DataType = 0x1602 //16位无时间计数器值越限事件
	InputCountLimitNoTimeSub32       DataType = 0x1603 //32位无时间计数器差值越限事件
	InputCountLimitNoTimeSub16       DataType = 0x1604 //16位无时间计数器差值越限事件
	InputCountLimitTime32            DataType = 0x1605 //32位带时间计数器值越限事件
	InputCountLimitTime16            DataType = 0x1606 //16位带时间计数器值越限事件
	InputCountLimitTimeSub32         DataType = 0x1607 //32位带时间计数器差值越限事件
	InputCountLimitTimeSub16         DataType = 0x1608 //16位带时间计数器差值越限事件
	InputCountLimitFreezeNoTime32    DataType = 0x1701 //32位无时间冻结计数器值越限事件
	InputCountLimitFreezeNoTime16    DataType = 0x1702 //16位无时间冻结计数器值越限事件
	InputCountLimitFreezeNoTimeSub32 DataType = 0x1703 //32位无时间冻结计数器差值越限事件
	InputCountLimitFreezeNoTimeSub16 DataType = 0x1704 //16位无时间冻结计数器差值越限事件
	InputCountLimitFreezeTime32      DataType = 0x1705 //32位带时间冻结计数器值越限事件
	InputCountLimitFreezeTime16      DataType = 0x1706 //16位带时间冻结计数器值越限事件
	InputCountLimitFreezeTimeSub32   DataType = 0x1707 //32位带时间冻结计数器差值越限事件
	InputCountLimitFreezeTimeSub16   DataType = 0x1708 //16位带时间冻结计数器差值越限事件

	InputAnalog32                  DataType = 0x1E01 //32位模拟量输入值
	InputAnalog16                  DataType = 0x1E02 //16位模拟量输入值
	InputAnalogNoStatus32          DataType = 0x1E03 //32位无状态模拟量输入值
	InputAnalogNoStatus16          DataType = 0x1E04 //16位无状态模拟量输入值
	InputAnalogFreeze32            DataType = 0x1F01 //32位冻结模拟量输入值
	InputAnalogFreeze16            DataType = 0x1F02 //16位冻结模拟量输入值
	InputAnalogFreezeTime32        DataType = 0x1F03 //32位带时间冻结模拟量输入值
	InputAnalogFreezeTime16        DataType = 0x1F04 //16位带时间冻结模拟量输入值
	InputAnalogFreezeNoStatus32    DataType = 0x1F05 //32位无状态冻结模拟量输入值
	InputAnalogFreezeNoStatus16    DataType = 0x1F06 //16位无状态冻结模拟量输入值
	InputAnalogLimitNoTime32       DataType = 0x2001 //32位无时间模拟量输入值越限事件
	InputAnalogLimitNoTime16       DataType = 0x2002 //16位无时间模拟量输入值越限事件
	InputAnalogLimitTime32         DataType = 0x2003 //32位带时间模拟量输入越限事件
	InputAnalogLimitTime16         DataType = 0x2004 //16位带时间模拟量输入越限事件
	InputAnalogLimitFreezeNoTime32 DataType = 0x2101 //32位无时间冻结模拟量输入值越限事件
	InputAnalogLimitFreezeNoTime16 DataType = 0x2102 //16位无时间冻结模拟量输入值越限事件
	InputAnalogLimitFreezeTime32   DataType = 0x2103 //32位带时间冻结模拟量输入值越限事件
	InputAnalogLimitFreezeTime16   DataType = 0x2104 //16位带时间冻结模拟量输入值越限事件

	OutputAnalog32      DataType = 0x2801 //32位模拟量输出
	OutputAnalog16      DataType = 0x2802 //16位模拟量输出
	OutputAnalogBlock32 DataType = 0x2903 //32位模拟量输出块(block)
	OutputAnalogBlock16 DataType = 0x2904 //16位模拟量输出块(block)

	DateTime          DataType = 0x3201 //日期和时间
	DateTimeSustain   DataType = 0x3202 //持续日期和时间
	TimeDelayAlmost   DataType = 0x3401 //近似延时(16位,秒)
	TimeDelayAccurate DataType = 0x3402 //精确延时(16位,秒)

	Class0 DataType = 0x3C01 //class0类数据,所有非1,2,3类数据
	Class1 DataType = 0x3C02 //class1类数据(通常为某组信息对象的变化)
	Class2 DataType = 0x3C03 //class2类数据(通常为某组信息对象的变化)
	Class3 DataType = 0x3C04 //class3类数据(通常为某组信息对象的变化)
)

func Decode(bs []byte) (*Pkg, error) {
	if len(bs) < 13 {
		return nil, fmt.Errorf("基础数据长度错误,预期(%d),得到(%d)", 13, len(bs))
	}

	if conv.Uint16(bs[:2]) != Prefix {
		return nil, fmt.Errorf("前缀错误,预期(%x),得到(%x)", Prefix, bs[:2])
	}

	length := int(bs[2])
	if len(bs) != length+7 {
		return nil, fmt.Errorf("数据长度错误,预期(%d),得到(%d)", length+7, len(bs))
	}

	p := &Pkg{
		Header: Header{
			Control: LinkControl{
				ToSlave:  bs[3] >= 0x80,
				IsMaster: bs[3]<<1 >= 0x80,
				Correct:  bs[3]<<2 >= 0x80,
				Function: HeaderFunction(bs[3] << 3 >> 3),
			},
			From: conv.Uint16([]byte{bs[7], bs[6]}),
			To:   conv.Uint16([]byte{bs[5], bs[4]}),
		},
		Body: Body{},
	}

	if conv.Uint16(bs[8:10]) != p.Header.CRC16(uint8(length)) {
		return nil, fmt.Errorf("CRC校验错误,预期(%x),得到(%x)", p.Header.CRC16(uint8(length)), conv.Uint16(bs[8:10]))
	}

	p.Body = Body{
		PkgNo: PkgNo{
			IsLast:  bs[10] >= 0x80,
			IsFirst: bs[10]<<1 >= 0x80,
			Current: bs[10] & 0x3F,
		},
		Control: BodyControl{
			IsFirst:     bs[11] >= 0x80,
			IsFinlay:    bs[11]<<1 >= 0x80,
			NeedAck:     bs[11]<<2 >= 0x80,
			Unsolicited: bs[11]<<3 >= 0x80,
			No:          bs[11] & 0x0F,
		},
		Function: BodyFunction(bs[12]),
	}

	bs = bs[13 : len(bs)-2]
	for len(bs) >= 2 {
		data := Data{
			DataType: DataType(conv.Uint16(bs[:2])),
		}
		if len(bs) > 2 {
			data.Qualifier = BodyQualifier{
				Index: bs[2] << 1 >> 5,
				Code:  bs[2] << 4 >> 4,
			}
			bs = bs[3:]
			switch data.Qualifier.Code {
			case 6:
			case 0, 1, 2, 3, 4, 5:
				if len(bs) < 2 {
					return nil, errors.New("数据域范围(变程)长度错误")
				}
				data.Range = bs[:2]
				bs = bs[2:]
			case 7, 8, 9:
				if len(bs) < 1 {
					return nil, errors.New("数据域范围(变程)长度错误")
				}
				data.Range = bs[:1]
				bs = bs[1:]
			}
			p.Body.Datas = append(p.Body.Datas, data)
		} else {
			bs = bs[2:]
		}

	}

	return p, nil
}

func ReadFunc(r *bufio.Reader) ([]byte, error) {

	for {
		b1, err := r.ReadByte()
		if err != nil {
			return nil, err
		}
		if b1 == byte(Prefix>>8) {

			b2, err := r.ReadByte()
			if err != nil {
				return nil, err
			}

			if uint16(b1)<<8+uint16(b2) == Prefix {

				length, err := r.ReadByte()
				if err != nil {
					return nil, err
				}

				buf := make([]byte, length+4)
				if _, err := io.ReadAtLeast(r, buf, int(length+4)); err != nil {
					return nil, err
				}

				return append([]byte{b1, b2, length}, buf...), nil

			}
		}
	}

}

func DealUnsolicited(p *Pkg, c *io.Client) error {
	switch p.Body.Function {
	case UnsolicitedMessage:
		c.Tag().Set("from", p.Header.To)
		c.Tag().Set("to", p.Header.From)
		_, err := c.Write(ConfirmPkg(p.Header.To, p.Header.From).Bytes())
		return err
	}
	return nil
}
