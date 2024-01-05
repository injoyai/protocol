package dnp3

import (
	"encoding/hex"
	"github.com/injoyai/io"
	"github.com/injoyai/io/dial"
	"github.com/injoyai/logs"
	"strings"
	"testing"
	"time"
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

func TestDecode2(t *testing.T) {
	s := "0564\n0b\nc4\n0400\n0300\ne42b\nc3\nc2\n01\n3c0106\na385"
	s = strings.ReplaceAll(s, "\n", "")
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

func TestEncode2(t *testing.T) {
	p := ReadPkg(0x63, 0x64, []Data{
		{
			DataType: 0x3C01,
			Qualifier: BodyQualifier{
				Code: 6,
			},
		},
	})
	t.Log(p.Bytes().HEX())
}

func TestConfirmPkg(t *testing.T) {
	p := ConfirmPkg(0x63, 0x64)
	t.Log(p.Bytes().HEX())
}

func TestConnect(t *testing.T) {
	<-dial.RedialTCP("127.0.0.1:20000", func(c *io.Client) {
		c.Debug()
		c.SetPrintWithHEX()
		//c.SetReadFunc(ReadFunc)
		c.SetDealFunc(func(msg *io.IMessage) {
			p, err := Decode(msg.Bytes())
			if err != nil {
				logs.Err(err)
				return
			}
			switch p.Body.Function {
			case UnsolicitedMessage:
				c.Tag().Set("from", p.Header.To)
				c.Tag().Set("to", p.Header.From)
				msg.Write(ConfirmPkg(p.Header.To, p.Header.From).Bytes())
			}

		})
		c.GoTimerWriter(time.Second*5, func(w *io.IWriter) error {
			from := c.Tag().GetUint16("from", 3)
			to := c.Tag().GetUint16("to", 4)
			if from > 0 && to > 0 {
				_, err := w.Write(ReadPkg(from, to, []Data{
					{
						DataType: 0x3C01,
						Qualifier: BodyQualifier{
							Code: 6,
						},
					},
				}).Bytes())
				return err
			}
			return nil
		})
	}).DoneAll()
}
