package s7

import (
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/g"
	"github.com/injoyai/io"
	"github.com/injoyai/io/dial"
	"github.com/injoyai/logs"
	"testing"
	"time"
)

func TestReadDB(t *testing.T) {
	dial.RedialTCP("192.168.10.66:102", func(c *io.Client) {
		c.Debug()
		c.SetPrintWithHEX()
		c.SetWriteWithNil()
		c.SetReadFunc(ReadFunc)
		c.WriteBytes(ShakeHandS7200Smart1.Bytes())
		c.WriteBytes(ShakeHandS7200Smart2.Bytes())
		c.SetDealFunc(func(msg *io.IMessage) {
			if ShakeHandS7200Smart1Regexp.Regexp(msg.HEX()) ||
				ShakeHandS7200Smart2Regexp.Regexp(msg.HEX()) {
				return
			}
			p, err := Decode(msg.Bytes())
			if !logs.PrintErr(err) {
				g.Done(conv.String(p.MsgID), p)
			}
		})
		go func() {

			{
				time.Sleep(time.Second)
				p := ReadMBit(3 * 8)
				c.Write(p.Encode())
				v, err := g.Wait(conv.String(p.MsgID))
				if err != nil {
					t.Error(err)
				} else {
					logs.Debug(GetInt(v.(*Response).Value))
				}
			}

			{
				time.Sleep(time.Second)
				p := ReadDBBit(300 * 8).SetMsgID(100)
				c.Write(p.Encode())
				v, err := g.Wait(conv.String(p.MsgID))
				if err != nil {
					t.Error(err)
				} else {
					logs.Debug(GetBool(v.(*Response).Value, 0))
				}
			}

			{
				time.Sleep(time.Second)
				p := ReadDBByte(310*8, 4)
				c.Write(p.Encode())
				v, err := g.Wait(conv.String(p.MsgID))
				if err != nil {
					t.Error(err)
				} else {
					logs.Debug(GetInt(v.(*Response).Value))
				}
			}

			{
				time.Sleep(time.Second)
				p := ReadDBByte(325*8, 4)
				c.Write(p.Encode())
				v, err := g.Wait(conv.String(p.MsgID))
				if err != nil {
					t.Error(err)
				} else {
					logs.Debug(GetFloat(v.(*Response).Value))
				}
			}

		}()

	})

	for {

	}
}

func TestWriteDB(t *testing.T) {
	dial.RedialTCP("192.168.10.66:102", func(c *io.Client) {
		c.Debug()
		c.SetPrintWithHEX()
		c.SetReadFunc(ReadFunc)
		c.WriteBytes(ShakeHandS7200Smart1.Bytes())
		c.WriteBytes(ShakeHandS7200Smart2.Bytes())
		c.SetDealFunc(func(msg *io.IMessage) {
			if ShakeHandS7200Smart1Regexp.Regexp(msg.HEX()) ||
				ShakeHandS7200Smart2Regexp.Regexp(msg.HEX()) {
				return
			}
			p, err := Decode(msg.Bytes())
			if !logs.PrintErr(err) {
				g.Done(conv.String(p.MsgID), p)
			}
		})
		go func() {

			{
				time.Sleep(time.Second)
				p := WriteDBBit(300*8, false)
				c.Write(p.Encode())
				v, err := g.Wait(conv.String(p.MsgID))
				if err != nil {
					t.Error(err)
				} else {
					logs.Debug(GetFloat(v.(*Response).Value))
				}
			}

			{
				time.Sleep(time.Second)
				p := WriteDBByte(325*8, conv.Bytes(float32(1.011)))
				c.Write(p.Encode())
				v, err := g.Wait(conv.String(p.MsgID))
				if err != nil {
					t.Error(err)
				} else {
					logs.Debug(GetFloat(v.(*Response).Value))
				}
			}

			{
				time.Sleep(time.Second)
				p := WriteDBByte(330*8, append([]byte{5}, []byte("close")...))
				c.Write(p.Encode())
				v, err := g.Wait(conv.String(p.MsgID))
				if err != nil {
					t.Error(err)
				} else {
					logs.Debug(GetFloat(v.(*Response).Value))
				}
			}

			{
				time.Sleep(time.Second)
				p := WriteDBByte(320*8, []byte{1, 1, 0, 0})
				c.Write(p.Encode())
				v, err := g.Wait(conv.String(p.MsgID))
				if err != nil {
					t.Error(err)
				} else {
					logs.Debug(GetFloat(v.(*Response).Value))
				}
			}

			{
				time.Sleep(time.Second)
				p := WriteDBByte(310*8, conv.Bytes(uint32(0)))
				c.Write(p.Encode())
				v, err := g.Wait(conv.String(p.MsgID))
				if err != nil {
					t.Error(err)
				} else {
					logs.Debug(GetFloat(v.(*Response).Value))
				}
			}

		}()

	})

	for {

	}
}
