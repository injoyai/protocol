package dlt645

import (
	"github.com/injoyai/base/bytes"
	"github.com/injoyai/conv"
	"testing"
)

func TestControl_Byte(t *testing.T) {
	t.Log(ControlGet.Byte())
	t.Log(ControlSet.Byte())
	t.Log(bytes.Entity([]byte{0x34, 0x89, 0x67, 0x45, 0x4e, 0x33, 0x8c, 0x8c, 0x56, 0x64, 0x45, 0xcc}).Sub0x33().Reverse().HEX())
}

func TestMark(t *testing.T) {
	t.Log(MarkSetSwitchOff.ReverseAdd0x33().HEX())
	t.Log(Control(conv.BINStr(byte(0x1c))[3:8]).String())

}
