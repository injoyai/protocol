package dnp3

import (
	"encoding/hex"
	"strings"
	"testing"
)

func TestPkg_Bytes(t *testing.T) {
	p := &Pkg{
		Header: Header{
			Control: LinkControl{
				ToSlave:  true,
				IsMaster: true,
				Correct:  false,
				Function: Send,
			},
			From: 0x0100,
			To:   0x0004,
		},
		Body: Body{
			PkgNo:    PkgNo{},
			Control:  BodyControl{},
			Function: Read,
			Datas: []Data{
				{
					DataType: Class1,
					Qualifier: BodyQualifier{
						Code: 6,
					},
				},
				{
					DataType: Class2,
					Qualifier: BodyQualifier{
						Code: 6,
					},
				},
				{
					DataType: Class3,
					Qualifier: BodyQualifier{
						Code: 6,
					},
				},
				{
					DataType: Class0,
					Qualifier: BodyQualifier{
						Code: 6,
					},
				},
			},
		},
	}
	t.Log(p.Bytes().HEX())
}

func TestPkgNo_Byte(t *testing.T) {
	p := PkgNo{}
	t.Logf("%#v", p.Byte())
}

func TestDecode(t *testing.T) {
	s := "05 64 14 c4 00 04 01 00 e9 b6 c1 c1 01 3c 02 06 3c 03 06 3c 04 06 3c 01 06 62 01"
	s = strings.ReplaceAll(s, " ", "")
	bs, err := hex.DecodeString(s)
	if err != nil {
		t.Error(err)
		return
	}
	p, err := Decode(bs)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", *p)
	t.Log(p.Bytes().HEX() == s)
}
