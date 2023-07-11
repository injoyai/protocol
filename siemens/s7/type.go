package s7

// OrderType 命令类型
type OrderType byte

const (
	CPUServer          OrderType = 0x00
	SetupCommunication OrderType = 0xF0
	ReadVar            OrderType = 0x04 //读
	WriteVar           OrderType = 0x05 //写
	DownloadRequest    OrderType = 0x1A
	DownloadBlock      OrderType = 0x1B
	DownloadEnd        OrderType = 0x1C
	UploadStart        OrderType = 0x1D
	Upload             OrderType = 0x1E
	UploadEnd          OrderType = 0x1F
	PLCStart           OrderType = 0x28
	PLCStop            OrderType = 0x29
)

// DataType 数据类型
// 示例: 变量的示例地址是DB123X 2.1，
// 它访问DB块123的第三个Byte的第二个Bit
type DataType byte

const (
	Bit     DataType = 0x01 //一个无符号的bit
	Byte    DataType = 0x02 //一个8位的数字
	Char    DataType = 0x03 //一个字符
	Uint16  DataType = 0x04 //两个字节宽的无符号整数
	Int16   DataType = 0x05 //两个字节宽的有符号整数
	Uint32  DataType = 0x06 //四字节宽的无符号整数
	Int32   DataType = 0x07 //四字节宽的有符号整数
	Float32 DataType = 0x08 //四个字节宽的IEEE浮点数
	Count   DataType = 0x1C //PLC程序计数器使用的计数器类型
)

type Block byte

func (this Block) Bytes() [2]byte {
	switch this {
	case Input, Output, Marker, CPU:
		return [2]byte{0, byte(this)}
	}
	// DataBlock
	return [2]byte{1, byte(this)}
}

var (
	MapBlock = map[string]Block{
		"I":   Input,
		"Q":   Output,
		"M":   Marker,
		"DB":  DataBlock,
		"DI":  DI,
		"L":   L,
		"V":   V,
		"C":   Counter,
		"T":   Timer,
		"CPU": CPU,
	}
)

const (

	//内存地址

	// CPU 状态
	CPU Block = 0x03
	// Input 数字和模拟输入模块值，映射到存储器
	Input Block = 0x81
	// Output 类似的存储器映射输出
	Output Block = 0x82
	// Marker 任意标记变量或标志寄存器驻留在这里
	Marker Block = 0x83
	// DataBlock DB区域是存储设备不同功能所需数据的最常见位置，这些数据块编号为地址的一部分。
	DataBlock Block = 0x84
	// DI 背景数据块
	DI Block = 0x85
	// L 局部变量
	L Block = 0x86
	// V 全局变量
	V Block = 0x87
	// Counter PLC程序使用的不同计数器的值
	Counter Block = 0x1E
	// Timer PLC程序使用的不同定时器的值
	Timer Block = 0x1F
)

// MsgType 消息的一般类型（有时称为ROSCTR类型），消息的其余部分在很大程度上取决于Message Type和功能代码
type MsgType byte

const (
	// JobRequest 主站发送的请求（例如读/写存储器，读/写块，启动/停止设备，通信设置）
	JobRequest MsgType = 0x01
	// Ack 从站发送的简单确认没有数据字段（从未见过它由S300 / S400设备发送）
	Ack MsgType = 0x02
	// AckData 带有可选数据字段的应答，包含对Job Request的回复
	AckData MsgType = 0x03
	// Userdata 原始协议的扩展，参数字段包含请求/响应Id（用于编程/调试，SZL读取，安全功能，时间设置，循环读取...）
	Userdata MsgType = 0x04
)

type Result byte

func (this Result) Err() error {
	if this != Success {
		return this
	}
	return nil
}

func (this Result) Error() string {
	switch this {
	case Reserved:
		return "预留"
	case HardwareError:
		return "硬件错误"
	case AccessNotAllowed:
		return "对象不允许访问"
	case InvalidAddress:
		return "无效地址"
	case DataTypeNotSupported:
		return "数据类型不支持"
	case DateTypeInconsistent:
		return "日期类型不一致"
	case ObjectNotExist:
		return "对象不存在"
	case Success:
		return "成功"
	}
	return "未知"
}

const (
	Reserved             Result = 0x00 //未定义，预留
	HardwareError        Result = 0x01 //硬件错误
	AccessNotAllowed     Result = 0x03 //对象不允许访问
	InvalidAddress       Result = 0x05 //无效地址，所需的地址超出此PLC的极限
	DataTypeNotSupported Result = 0x06 //数据类型不支持
	DateTypeInconsistent Result = 0x07 //日期类型不一致
	ObjectNotExist       Result = 0x0A //对象不存在
	Success              Result = 0xFF //成功
)
