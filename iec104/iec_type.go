package iec104

// Type 类型标识
type Type uint8

func (this Type) String() string {
	switch this {
	case M_ME_NA_1:
		return "带品质描述的遥测 归一化遥测（整型) 值占3字节"
	case TypeYCMemoSB:
		return "带3个字节时标的且具有品质描述的遥测值，每个遥测值占6个字节"
	case M_ME_NB_1:
		return "不带时标的标度化值，每个遥测值占3个字节"
	case TypeYCSB:
		return "带3个字节时标的标度化值，每个遥测值占6个字节"
	case M_ME_NC_1:
		return "带品质描述的浮点值，每个遥测值占5个字节"
	case TypeYCMemoSBFloat:
		return "带3个字节时标且具有品质描述的浮点值，每个遥测值占8个字节"
	case M_ME_ND_1:
		return "不带品质描述的遥测值，每个遥测值占2个字节"

	case M_SP_NA_1:
		return "单点遥信,1字节"
	case TypeYXOneSB3:
		return "带3个字节短时标的单点遥信"
	case M_DP_NA_1:
		return "双点遥信,1字节"
	case TypeYXTwoSB:
		return "带3个字节短时标的双点遥信"
	case M_ST_NA_1:
		return "步位置信息,带描述"
	case TypeYXStatus:
		return "具有状态变位检测的成组单点遥信，每个字节包括8个遥信"
	case M_SP_TB_1:
		return "带7个字节时标的单点遥信"
	case M_MP_TB_1:
		return "双点遥信,带时标,8字节"

	case M_IT_NA_1:
		return "不带时标的电度量，每个电度量占5个字节"
	case TypeYMSB8:
		return "带3个字节短时标的电度量，每个电度量占8个字节"
	case M_IT_TB_1, M_IT_TB_1_:
		return "带7个字节长时标的电度量，每个电度量占12个字节"

	case C_SC_NA_1:
		return "单点遥控"
	case C_DC_NA_1:
		return "双点遥控"
	case C_SE_NA_1:
		return "规一化设定值"
	case C_SE_NB_1:
		return "标度化设定值"
	case C_SE_NC_1:
		return "短浮点设定值"
	case M_EI_NA_1:
		return "初始化结束"
	case C_IC_NA_1:
		return "总召唤"
	case C_CI_NA_1:
		return "累积量召唤"
	case C_RD_NA_1:
		return "读取单个对象信息对象值"
	case TypeClock:
		return "时钟同步"
	case C_TS_TA_1:
		return "带时标的测试命令"
	}

	return "未知类型"
}

const (

	/*
		遥测 YC ,模拟输入AI
	*/

	M_ME_NA_1         Type = 0x09 //带品质描述的遥测 归一化遥测（整型) 值占3字节
	TypeYCMemoSB      Type = 0x0A //带3个字节时标的且具有品质描述的遥测值，每个遥测值占6个字节
	M_ME_NB_1         Type = 0x0B //不带时标的标度化值，每个遥测值占3个字节
	TypeYCSB          Type = 0x0C //带3个字节时标的标度化值，每个遥测值占6个字节
	M_ME_NC_1         Type = 0x0D //带品质描述的浮点值，每个遥测值占5个字节
	TypeYCMemoSBFloat Type = 0x0E //带3个字节时标且具有品质描述的浮点值，每个遥测值占8个字节
	M_ME_ND_1         Type = 0x15 //不带品质描述的遥测值，每个遥测值占2个字节

	/*
		遥信 YX ,数字输入DI 开关量
	*/

	M_SP_NA_1    Type = 0x01 //单点遥信,1字节
	TypeYXOneSB3 Type = 0x02 //带3个字节短时标的单点遥信
	M_DP_NA_1    Type = 0x03 //双点遥信,1字节
	TypeYXTwoSB  Type = 0x04 //带3个字节短时标的双点遥信
	M_ST_NA_1    Type = 0x05 //步位置信息,带描述,无时标
	TypeYXStatus Type = 0x14 //具有状态变位检测的成组单点遥信，每个字节包括8个遥信
	M_SP_TB_1    Type = 0x1E //单点遥信,带描述,带时标,8字节
	M_MP_TB_1    Type = 0x1F //双点遥信,带时标,8字节

	/*
		遥脉 YM ,电度量,是指对现场某装置所发出的脉冲信号进行周期累计的一种远程计数操作
	*/

	M_IT_NA_1  Type = 0x0F //累计量,5字节
	TypeYMSB8  Type = 0x10 //带3个字节短时标的电度量，每个电度量占8个字节
	M_IT_TB_1_ Type = 0x24 //带7个字节长时标的电度量，每个电度量占12个字节
	M_IT_TB_1  Type = 0x25 //带7个字节长时标的电度量，每个电度量占12个字节

	/*
		其他 遥控(YK) 召唤(全部数据/全部遥脉)
	*/

	C_SC_NA_1 Type = 0x2D //单点遥控
	C_DC_NA_1 Type = 0x2E //双点遥控
	C_SE_NA_1 Type = 0x30 //规一化设定值
	C_SE_NB_1 Type = 0x31 //标度化设定值
	C_SE_NC_1 Type = 0x32 //短浮点设定值
	M_EI_NA_1 Type = 0x46 //初始化结束

	C_IC_NA_1 Type = 0x64 //总召唤
	C_CI_NA_1 Type = 0x65 //累积量召唤
	C_RD_NA_1 Type = 0x66 //读命令
	TypeClock Type = 0x67 //时钟同步
	C_TS_TA_1 Type = 0x6B //带时标的测试命令

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
