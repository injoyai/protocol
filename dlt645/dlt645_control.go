package dlt645

type Control string

func (this Control) Byte() (b byte) {
	for _, v := range this {
		b *= 2
		if v != '0' {
			b += 1
		}
	}
	return
}

func (this Control) String() string {
	switch this {
	case ControlNull:
		return "保留"
	case ControlSetTime:
		return "广播校时"
	case ControlGet:
		return "读数据"
	case ControlGetNext:
		return "读后续数据"
	case ControlGetNo:
		return "读通讯地址"
	case ControlSet:
		return "写数据"
	case ControlSetSwitch:
		return "断闸,合闸,报警,报警解除,保电,保电解除"
	case ControlSetNo:
		return "写通信地址"
	case ControlSetSave:
		return "冻结命令"
	case ControlSetRate:
		return "修改通信速率"
	case ControlSetPassword:
		return "修改密码"
	case ControlSetResetMax:
		return "最大需量清零"
	case ControlSetResetUse:
		return "电表清零"
	case ControlSetResetEvent:
		return "事件清零"
	}
	return "未知控制码"
}

var (
	/*
		二进制
		D7(传送方向:0是主站发出的命令帧,1是从站发出的应答帧)
		D6(从站应答标志:0是从站正确应答,1是从站异常应答)
		D5(后续帧标示:0无后续数据帧,有后续数据帧)
		D4~D0(功能码:
				00000:保留
				01000:广播校时
				10001:读数据
				10010:读后续数据
				10011:读通讯地址
				10100:写数据
				10101:写通信地址
				10110:冻结命令
				10111:修改通信速率
				11000:修改密码
				11001:最大需量清零
				11010:电表清零
				11011:事件清零
				)
	*/

	ControlNull          Control = "00000" //保留
	ControlSetTime       Control = "01000" //广播校时
	ControlGet           Control = "10001" //读数据
	ControlGetNext       Control = "10010" //读后续数据
	ControlGetNo         Control = "10011" //读通讯地址
	ControlSet           Control = "10100" //写数据
	ControlSetSwitch     Control = "11100" //断闸,合闸,报警,报警解除,保电,保电解除
	ControlSetNo         Control = "10101" //写通信地址
	ControlSetSave       Control = "10110" //冻结命令
	ControlSetRate       Control = "10111" //修改通信速率
	ControlSetPassword   Control = "11000" //修改密码
	ControlSetResetMax   Control = "11001" //最大需量清零
	ControlSetResetUse   Control = "11010" //电表清零
	ControlSetResetEvent Control = "11011" //事件清零
)
