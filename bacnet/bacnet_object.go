package bacnet

import "github.com/injoyai/base/g"

type Object struct {
	Tag   Tag
	Class bool
	Data  []byte
}

func (this Object) Bytes() g.Bytes {
	return nil
}

type Tag uint8

const (
	TagNull                   Tag = iota
	TagBoolean                    //bool
	TagUnsignedInteger            //uint
	TagSignedInteger              //int
	TagReal                       //float32
	TagDouble                     //float64
	TagOctetString                //八进制字符?
	TagCharacterString            //字节字符?
	TagBitString                  //二进制字符?
	TagEnumerated                 //枚举
	TagData                       //
	TagTime                       //时间
	TagBACnetObjectIdentifier     //
)

// ObjectType
// 对象是对现实设备中某一特征的抽象
type ObjectType uint8

const (

	// AnalogInputObjectType 模拟输入对象类型
	// 模拟输入对象类型定义为一个标准对象，其属性表示一个模拟输入的外部可见一致性代码。
	AnalogInputObjectType ObjectType = iota

	// AnalogOutputObjectType 模拟输出对象类型
	// 模拟输出对象类型定义为一个标准对象，其属性表示一个模拟输出的外部可见一致性代码。
	AnalogOutputObjectType

	// AnalogValueObjectType 模拟值对象类型
	// 模拟值对象类型定义为一个标准对象，其属性表示一个模拟值的外部可见一致性代码。
	// BACnet设备的一个“模拟值”是驻留在这个设备的内存中的一个控制系统参数。
	AnalogValueObjectType

	// BinaryInputObjectType 二进制输入对象类型
	// 二进制输入对象类型定义为一个标准对象，它的属性表示二进制输入的外部可见一致性代码。
	// “二进制输入”是物理设备或硬件的输入，该输入只存在两种状态，即“活动（ACTIVE）”状态和“非活动（INACTIVE）”状态。
	// 二进制输入的主要用途是指明机械设备状态，如：风机或水泵是否运行。活动表示设备开或运转，非活动表示设备关或未运行。
	BinaryInputObjectType

	// BinaryOutputObjectType 二进制输出对象类型
	// 二进制输出对象类型定义为一个标准对象，它的属性表示二进制输出的外部可见一致性代码。
	// “二进制输出”是物理设备或硬件的输出，该输出只存在两种状态，即“活动”状态和“非活动”状态。
	// 二进制输出的主要用途是切换机械设备状态，如：风机或水泵的开和关。活动表示设备开或运转，非活动表示设备关或未运行。
	BinaryOutputObjectType

	// BinaryValueObjectType 二进制值对象类型
	// 二进制值对象类型定义为一个标准对象，它的属性表示二进制值的外部可见一致性代码。
	// “二进制值”是驻留在BACnet设备内存中的控制系统参数。这个参数只存在两种状态即：“活动”状态和“非活动”状态。
	BinaryValueObjectType

	// CalendarObjectType 日期表对象类型
	// 日期表对象类型定义为一个标准对象，用于描述日期列表，例如，“节假日”、“特别日”或简单的日期列表。
	CalendarObjectType

	// CommandObjectType 命令对象类型
	// 命令对象类型定义为一个标准对象，其属性反映了多操作命令过程的外部可见一致性代码。
	// 命令对象的作用是，根据写入到命令对象自己的当前值属性中的“操作代码（action code）”，向一组对象属性写入一组值。
	// 无论何时，只要命令对象的当前值属性被写入，就会触发命令对象采取一组改变其它对象的属性值的操作。
	CommandObjectType

	// DeviceObjectType 设备对象类型
	// 设备对象类型定义为一个标准对象，其属性表示BACnet设备的外部可见一致性代码。
	// 每个BACnet设备有且只有一个设备对象。
	// 每个设备对象由它的对象标识符属性确定，该属性在BACnet设备中乃至整个BACnet互联网中都是唯一的。
	DeviceObjectType

	// EventEnrollmentObjectType 事件登记对象类型
	// 事件登记对象类型定义为一个标准对象，表示BACnet系统内管理事件的信息。
	// “事件”是指满足预先规定条件的所有对象的任何属性值的变化。
	// 事件登记对象主要用于定义一个事件和提供在事件发生与通告消息向一个或多个接收者进行传输这两者之间的联系。
	EventEnrollmentObjectType

	// FileObjectType 文件对象类型
	// 文件对象类型定义为一个标准对象，用于定义可以通过文件服务（见第14节）访问的数据文件的属性。
	FileObjectType

	// GroupObjectType 组对象类型
	// 组对象类型定义为一个标准对象，其属性表示一个其它对象的集合以及这些对象的一个或多个属性。
	// 组对象提供一种快速的方式，可以一次确定组的成员，从而简化BACnet设备间的信息交换。
	// 一个组对象可以是任何对象类型的组合。
	GroupObjectType

	// LoopObjectType 环对象类型
	// 环对象类型定义为一个标准对象，其属性表示任何形式的反馈控制环路的外部可见一致性代码。
	// 环对象通过提供三个独立的无单位增益常数，可以具有广泛的适用性。
	// 每个增益常数由控制算法具休确定，如何使用不同的算法确定增益常数的方法，由生产商自行确定。
	LoopObjectType

	// MultiStateInputObjectType 多态输入对象类型
	// 多态输入对象类型定义了一个标准对象，它的当前值属性表示对象驻留的BACnet设备内算法处理的结果。
	MultiStateInputObjectType

	// MultiStateOutputObjectType 多态输出对象类型
	// 多态输出对象类型定义了一个标准对象，它的属性表示这个对象驻留的BACnet设备内的处理程序或一个或多个物理输出的期望状态。
	MultiStateOutputObjectType

	// NotificationClassObjectType 通告类对象类型
	// 通告类对象类型定义了一个标准对象，表示在BACnet系统内事件通告发布所需的信息。
	NotificationClassObjectType

	// ProgramObjectType 程序对象类型
	// 程序对象类型定义了一个标准对象，它的属性表示应用程序的外部可视一致性代码。
	// 在本协议中，应用程序是指对一个在BACnet设备中的处理过程的抽象表示，这个处理过程执行一个指令集，对某个数据结构集合进行操作。
	ProgramObjectType

	// ScheduleObjectType 时间表对象类型
	// 时间表对象类型定义了一个标准对象，用于描述一个周期性的时间表。
	// 这个时间表中确定了某事件在一个日期范围内可能重复发生，同时表示有些日期是事件不发生的日期。
	ScheduleObjectType
)
