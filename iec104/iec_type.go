package iec104

// Type 类型标识
type Type uint8

func (this Type) String() string {
	switch this {
	case TypeYCMemo:
		return "带品质描述的遥测 归一化遥测（整型) 值占3字节"
	case TypeYCMemoSB:
		return "带3个字节时标的且具有品质描述的遥测值，每个遥测值占6个字节"
	case TypeYC3:
		return "不带时标的标度化值，每个遥测值占3个字节"
	case TypeYCSB:
		return "带3个字节时标的标度化值，每个遥测值占6个字节"
	case TypeYCMemoFloat:
		return "带品质描述的浮点值，每个遥测值占5个字节"
	case TypeYCMemoSBFloat:
		return "带3个字节时标且具有品质描述的浮点值，每个遥测值占8个字节"
	case TypeYC2:
		return "不带品质描述的遥测值，每个遥测值占2个字节"

	case TypeYXOne:
		return "单点遥信,对应数字输入DI,每个遥信占1个字节"
	case TypeYXOneSB3:
		return "带3个字节短时标的单点遥信"
	case TypeYXTwo:
		return "双点遥信,每个遥信占1个字节"
	case TypeYXTwoSB:
		return "带3个字节短时标的双点遥信"
	case TypeYXStatus:
		return "具有状态变位检测的成组单点遥信，每个字节包括8个遥信"
	case TypeYXOneSB7:
		return "带7个字节时标的单点遥信"
	case TypeYXTwoSB7:
		return "带7个字节时标的双点遥信"

	case TypeYM:
		return "不带时标的电度量，每个电度量占5个字节"
	case TypeYMSB8:
		return "带3个字节短时标的电度量，每个电度量占8个字节"
	case TypeYMSB12_:
		return "带7个字节长时标的电度量，每个电度量占12个字节"
	case TypeYMSB12:
		return "带7个字节长时标的电度量，每个电度量占12个字节"
	case TypeYMFloat:
		return "短浮点数下设定值"

	case TypeYKOne:
		return "单点遥控"
	case TypeYKTwo:
		return "双点遥控"
	case TypeInit:
		return "初始化结束"
	case TypeZHTotal:
		return "召唤全部数据"
	case TypeZHYM:
		return "召唤全部电度(遥脉)"
	case TypeRead:
		return "读取单个对象信息对象值"
	case TypeClock:
		return "时钟同步"
	case TypeTestSB:
		return "带时标的测试命令"
	}

	return "未知类型"
}

const (

	/*
		遥测 YC ,模拟输入AI
	*/

	TypeYCMemo        Type = 0x09 //带品质描述的遥测 归一化遥测（整型) 值占3字节
	TypeYCMemoSB      Type = 0x0A //带3个字节时标的且具有品质描述的遥测值，每个遥测值占6个字节
	TypeYC3           Type = 0x0B //不带时标的标度化值，每个遥测值占3个字节
	TypeYCSB          Type = 0x0C //带3个字节时标的标度化值，每个遥测值占6个字节
	TypeYCMemoFloat   Type = 0x0D //带品质描述的浮点值，每个遥测值占5个字节
	TypeYCMemoSBFloat Type = 0x0E //带3个字节时标且具有品质描述的浮点值，每个遥测值占8个字节
	TypeYC2           Type = 0x15 //不带品质描述的遥测值，每个遥测值占2个字节

	/*
		遥信 YX ,数字输入DI 开关量
	*/

	TypeYXOne    Type = 0x01 //单点遥信,对应数字输入DI,每个遥信占1个字节
	TypeYXOneSB3 Type = 0x02 //带3个字节短时标的单点遥信
	TypeYXTwo    Type = 0x03 //双点遥信,每个遥信占1个字节
	TypeYXTwoSB  Type = 0x04 //带3个字节短时标的双点遥信
	TypeYXStatus Type = 0x14 //具有状态变位检测的成组单点遥信，每个字节包括8个遥信
	TypeYXOneSB7 Type = 0x1E //带7个字节时标的单点遥信
	TypeYXTwoSB7 Type = 0x1F //带7个字节时标的双点遥信

	/*
		遥脉 YM ,电度量,是指对现场某装置所发出的脉冲信号进行周期累计的一种远程计数操作
	*/

	TypeYM      Type = 0x0F //不带时标的电度量，每个电度量占5个字节
	TypeYMSB8   Type = 0x10 //带3个字节短时标的电度量，每个电度量占8个字节
	TypeYMFloat Type = 0x21 //短浮点数下设定值
	TypeYMSB12_ Type = 0x24 //带7个字节长时标的电度量，每个电度量占12个字节
	TypeYMSB12  Type = 0x25 //带7个字节长时标的电度量，每个电度量占12个字节

	/*
		其他 遥控(YK) 召唤(全部数据/全部遥脉)
	*/

	TypeYKOne   Type = 0x2D //单点遥控
	TypeYKTwo   Type = 0x2E //双点遥控
	TypeInit    Type = 0x46 //初始化结束
	TypeZHTotal Type = 0x64 //召唤全部数据
	TypeZHYM    Type = 0x65 //召唤全部电度(遥脉)
	TypeRead    Type = 0x66 //读单个信息对象值
	TypeClock   Type = 0x67 //时钟同步
	TypeTestSB  Type = 0x6B //带时标的测试命令

)

const (
	Order_A   = 0x01 //命令确认
	STARTDT_C = 0x07 //启动请求
	STARTDT_A = 0x0B //启动确认
	STOP_C    = 0x13 //停止请求
	STOP_A    = 0x23 //停止确认
	TESTFR_C  = 0x43 //测试请求
	TESTFR_A  = 0x83 //测试确认
)

type Reason uint16

func (this Reason) String() string {
	switch this {
	case ReasonCycle:
		return "周期循环"
	case ReasonScan:
		return "背景扫描"
	case ReasonUpload:
		return "主动上报"
	case ReasonInit:
		return "初始化"
	case ReasonAsk:
		return "请求或被请求"
	case ReasonStart:
		return "激活"
	case ReasonStartAck:
		return "激活确认"
	case ReasonStop:
		return "停止激活"
	case ReasonStopAck:
		return "停止激活确认"
	case ReasonBreak:
		return "激活终止"
	case ReasonZHResponse:
		return "响应总召唤"
	case ReasonResponse1:
		return "响应第1组召唤"
	}
	return "未知原因"
}

const (
	ReasonCycle      Reason = 0x0001 //周期循环
	ReasonScan       Reason = 0x0002 //背景扫描
	ReasonUpload     Reason = 0x0003 //主动上报
	ReasonInit       Reason = 0x0004 //初始化
	ReasonAsk        Reason = 0x0005 //请求或被请求
	ReasonStart      Reason = 0x0006 //激活
	ReasonStartAck   Reason = 0x0007 //激活确认
	ReasonStop       Reason = 0x0008 //停止激活
	ReasonStopAck    Reason = 0x0009 //停止激活确认
	ReasonBreak      Reason = 0x000A //激活终止
	ReasonZHResponse Reason = 0x0014 //响应总召唤
	ReasonResponse1  Reason = 0x0015 //响应第一组召唤
)

// Memo 品质描述
type Memo uint8

// Valid 数据是否有效
func (this Memo) Valid() bool {
	return this.IV() && !this.OV()
}

// Over 数据是否溢出
func (this Memo) Over() bool {
	return this.OV()
}

func (this Memo) IV() bool {
	return (this>>7)%2 == 1
}

func (this Memo) NT() bool {
	return (this>>6)%2 == 1
}

func (this Memo) SB() bool {
	return (this>>5)%2 == 1
}

func (this Memo) BL() bool {
	return (this>>4)%2 == 1
}

func (this Memo) OV() bool {
	return (this>>3)%2 == 1
}
