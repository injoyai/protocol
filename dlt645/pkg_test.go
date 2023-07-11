package dlt645

import (
	"encoding/hex"
	"testing"
)

func TestDecode(t *testing.T) {
	bs, _ := hex.DecodeString("68AAAAAAAAAAAA68910833333333343333337E16")
	pkg, err := Decode(bs)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(pkg.Data.Sub0x33())
	f, err := pkg.Data.Sub0x33ReverseHexToFloat64(2)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(f)
}

func TestEncode(t *testing.T) {
	p := &EnPkg{
		No:      []byte{0xaa, 0xaa, 0xaa, 0xaa, 0xaa, 0xaa},
		Control: ControlGet,
		Mark:    MarkGetUse,
		Data:    []byte{0, 0, 0, 01},
	}
	t.Log(p.HEX())
	p2, err := Decode(p.Bytes())
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(p2.No)
	t.Log(p2.Control)
	t.Log(p2.Mark.HEX())
	t.Log(p2.Data.Sub0x33ReverseHexToInt())
}
