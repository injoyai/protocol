package dlt645

import (
	"github.com/injoyai/base/bytes"
)

// Mark
// 第三字节是序号(例如0是总,1对应1通道),
// 第四字节是天数(例如:1是上一日)
type Mark bytes.Entity

func (this Mark) String() string {
	return this.HEX()
}

func (this Mark) HEX() string {
	return bytes.Entity(this).HEX()
}

func (this Mark) ReverseAdd0x33() Mark {
	return Mark(bytes.Entity(this).Reverse().Add0x33())
}

var (
	MarkSetReset     Mark = []byte{0x12, 0x34, 0x56, 0x01}                                                 //电表清零
	MarkSetSwitchOn  Mark = []byte{0x99, 0x12, 0x31, 0x23, 0x59, 0x59, 0x00, 0x1B, 0x12, 0x34, 0x56, 0x01} //合闸
	MarkSetSwitchOff Mark = []byte{0x99, 0x12, 0x31, 0x23, 0x59, 0x59, 0x00, 0x1A, 0x12, 0x34, 0x56, 0x01} //断闸

	MarkGetUse              Mark = []byte{0x00, 0x00, 0x00, 0x00} //抄读电量
	MarkGetElectricA        Mark = []byte{0x02, 0x02, 0x01, 0x00} //电流A
	MarkGetElectricB        Mark = []byte{0x02, 0x02, 0x01, 0x00} //电流B
	MarkGetElectricC        Mark = []byte{0x02, 0x02, 0x01, 0x00} //电流C
	MarkGetElectricAll      Mark = []byte{0x02, 0x02, 0x01, 0x00} //电流All
	MarkGetVoltageA         Mark = []byte{0x02, 0x01, 0x01, 0x00} //电压A
	MarkGetVoltageB         Mark = []byte{0x02, 0x01, 0x02, 0x00} //电压B
	MarkGetVoltageC         Mark = []byte{0x02, 0x01, 0x03, 0x00} //电压C
	MarkGetVoltageAll       Mark = []byte{0x02, 0x01, 0xFF, 0x00} //电压All
	MarkGetPowerActiveSum   Mark = []byte{0x02, 0x03, 0x00, 0x00} //瞬时有功功率Sum
	MarkGetPowerActiveA     Mark = []byte{0x02, 0x03, 0x00, 0x00} //瞬时有功功率A
	MarkGetPowerActiveB     Mark = []byte{0x02, 0x03, 0x01, 0x00} //瞬时有功功率B
	MarkGetPowerActiveC     Mark = []byte{0x02, 0x03, 0x02, 0x00} //瞬时有功功率C
	MarkGetPowerActiveAll   Mark = []byte{0x02, 0x03, 0xFF, 0x00} //瞬时有功功率All
	MarkGetPowerReactiveSum Mark = []byte{0x02, 0x04, 0x00, 0x00} //瞬时无功功率Sum
	MarkGetPowerReactiveA   Mark = []byte{0x02, 0x04, 0x00, 0x00} //瞬时无功功率A
	MarkGetPowerReactiveB   Mark = []byte{0x02, 0x04, 0x01, 0x00} //瞬时无功功率B
	MarkGetPowerReactiveC   Mark = []byte{0x02, 0x04, 0x02, 0x00} //瞬时无功功率C
	MarkGetPowerReactiveAll Mark = []byte{0x02, 0x04, 0xFF, 0x00} //瞬时无功功率All
	MarkGetPowerApparentSum Mark = []byte{0x02, 0x05, 0x00, 0x00} //瞬时视在功率Sum
	MarkGetPowerApparentA   Mark = []byte{0x02, 0x05, 0x00, 0x00} //瞬时视在功率A
	MarkGetPowerApparentB   Mark = []byte{0x02, 0x05, 0x01, 0x00} //瞬时视在功率B
	MarkGetPowerApparentC   Mark = []byte{0x02, 0x05, 0x02, 0x00} //瞬时视在功率C
	MarkGetPowerApparentAll Mark = []byte{0x02, 0x05, 0xFF, 0x00} //瞬时视在功率All
	MarkGetFactorSum        Mark = []byte{0x02, 0x06, 0x00, 0x00} //功率因数Sum
	MarkGetFactorA          Mark = []byte{0x02, 0x06, 0x01, 0x00} //功率因数A
	MarkGetFactorB          Mark = []byte{0x02, 0x06, 0x02, 0x00} //功率因数B
	MarkGetFactorC          Mark = []byte{0x02, 0x06, 0x03, 0x00} //功率因数C
	MarkGetFactorAll        Mark = []byte{0x02, 0x06, 0xFF, 0x00} //功率因数All
	MarkGetSaveUse          Mark = []byte{0x05, 0x06, 0x01, 0x00} //冻结数据
)

var MarkMap = map[string]MarkInfo{

	//下发设置
	MarkSetSwitchOn.HEX():  newMarkInfo(MarkGetUse, "合闸", "", ControlSetSwitch, 0),
	MarkSetSwitchOff.HEX(): newMarkInfo(MarkGetUse, "断闸", "", ControlSetSwitch, 0),
	MarkSetReset.HEX():     newMarkInfo(MarkSetReset, "清零", "", ControlSetResetUse, 0),

	//下发抄读
	MarkGetUse.HEX():              newMarkInfo(MarkGetUse, "电量", "kwh", ControlGet, 2).SetLength(4, 1),
	MarkGetElectricA.HEX():        newMarkInfo(MarkGetElectricA, "电流A", "A", ControlGet, 3).SetLength(3, 1),
	MarkGetElectricB.HEX():        newMarkInfo(MarkGetElectricB, "电流B", "A", ControlGet, 3).SetLength(3, 1),
	MarkGetElectricC.HEX():        newMarkInfo(MarkGetElectricC, "电流C", "A", ControlGet, 3).SetLength(3, 1),
	MarkGetElectricAll.HEX():      newMarkInfo(MarkGetElectricAll, "电流All", "A", ControlGet, 3).SetLength(9, 3),
	MarkGetVoltageA.HEX():         newMarkInfo(MarkGetVoltageA, "电压A", "V", ControlGet, 3).SetLength(2, 1),
	MarkGetVoltageB.HEX():         newMarkInfo(MarkGetVoltageB, "电压B", "V", ControlGet, 3).SetLength(2, 1),
	MarkGetVoltageC.HEX():         newMarkInfo(MarkGetVoltageC, "电压C", "V", ControlGet, 3).SetLength(2, 1),
	MarkGetVoltageAll.HEX():       newMarkInfo(MarkGetVoltageAll, "电压All", "V", ControlGet, 3).SetLength(6, 3),
	MarkGetPowerActiveSum.HEX():   newMarkInfo(MarkGetPowerActiveSum, "瞬时有功功率Sum", "KW", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerActiveA.HEX():     newMarkInfo(MarkGetPowerActiveA, "瞬时有功功率A", "KW", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerActiveB.HEX():     newMarkInfo(MarkGetPowerActiveB, "瞬时有功功率B", "KW", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerActiveC.HEX():     newMarkInfo(MarkGetPowerActiveC, "瞬时有功功率C", "KW", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerActiveAll.HEX():   newMarkInfo(MarkGetPowerActiveAll, "瞬时有功功率All", "KW", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerReactiveSum.HEX(): newMarkInfo(MarkGetPowerReactiveSum, "瞬时无功功率Sum", "KVar", ControlGet, 3).SetLength(12, 4),
	MarkGetPowerReactiveA.HEX():   newMarkInfo(MarkGetPowerReactiveA, "瞬时无功功率A", "KVar", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerReactiveB.HEX():   newMarkInfo(MarkGetPowerReactiveB, "瞬时无功功率B", "KVar", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerReactiveC.HEX():   newMarkInfo(MarkGetPowerReactiveC, "瞬时无功功率C", "KVar", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerReactiveAll.HEX(): newMarkInfo(MarkGetPowerReactiveAll, "瞬时无功功率All", "KVar", ControlGet, 3).SetLength(12, 4),
	MarkGetPowerApparentSum.HEX(): newMarkInfo(MarkGetPowerApparentSum, "瞬时视在功率Sum", "KVA", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerApparentA.HEX():   newMarkInfo(MarkGetPowerApparentA, "瞬时视在功率A", "KVA", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerApparentB.HEX():   newMarkInfo(MarkGetPowerApparentB, "瞬时视在功率B", "KVA", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerApparentC.HEX():   newMarkInfo(MarkGetPowerApparentC, "瞬时视在功率C", "KVA", ControlGet, 3).SetLength(3, 1),
	MarkGetPowerApparentAll.HEX(): newMarkInfo(MarkGetPowerApparentAll, "瞬时视在功率All", "KVA", ControlGet, 3).SetLength(12, 4),
	MarkGetFactorSum.HEX():        newMarkInfo(MarkGetFactorSum, "功率因数Sum", "", ControlGet, 3).SetLength(2, 1),
	MarkGetFactorA.HEX():          newMarkInfo(MarkGetFactorA, "功率因数A", "", ControlGet, 3).SetLength(2, 1),
	MarkGetFactorB.HEX():          newMarkInfo(MarkGetFactorB, "功率因数B", "", ControlGet, 3).SetLength(2, 1),
	MarkGetFactorC.HEX():          newMarkInfo(MarkGetFactorC, "功率因数C", "", ControlGet, 3).SetLength(2, 1),
	MarkGetFactorAll.HEX():        newMarkInfo(MarkGetFactorAll, "功率因数All", "", ControlGet, 3).SetLength(6, 4),
	MarkGetSaveUse.HEX():          newMarkInfo(MarkGetSaveUse, "日冻结电量", "", ControlGet, 1).SetLength(20, 4),
}

type MarkInfo struct {
	Mark     Mark    //标示
	Name     string  //名称
	Unit     string  //单位
	Length   int     //数据长度(字节)
	Number   int     //数据数量(数据长度/数据数量)
	Decimals uint8   //小数位数
	Control  Control //控制码
}

func (this MarkInfo) SetLength(length, number int) MarkInfo {
	this.Length = length
	this.Number = number
	return this
}

func (this MarkInfo) String() string {
	return this.Name
}

func newMarkInfo(mark Mark, name, unit string, control Control, decimals uint8) MarkInfo {
	return MarkInfo{
		Mark:     mark,
		Name:     name,
		Unit:     unit,
		Decimals: decimals,
		Control:  control,
	}
}
