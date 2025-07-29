package s7

import (
	"encoding/hex"
	"github.com/injoyai/base/maps/wait"
	"github.com/injoyai/conv"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/ios/client/redial"
	"github.com/injoyai/logs"
	"testing"
	"time"
)

func TestReadDB(t *testing.T) {
	redial.RunTCP("192.168.10.66:102", func(c *client.Client) {
		c.Logger.Debug()
		c.Logger.WithHEX()
		//c.SetWriteWithNil()
		c.OnReadFrom = ReadFrom
		c.Write(ShakeHandS7200Smart1.Bytes())
		c.Write(ShakeHandS7200Smart2.Bytes())
		c.OnDealMessage = func(c *client.Client, msg ios.Acker) {
			if ShakeHandS7200Smart1Regexp.Regexp(hex.EncodeToString(msg.Payload())) ||
				ShakeHandS7200Smart2Regexp.Regexp(hex.EncodeToString(msg.Payload())) {
				return
			}
			p, err := Decode(msg.Payload())
			if !logs.PrintErr(err) {
				wait.Done(conv.String(p.MsgID), p)
			}
		}
		go func() {

			{
				time.Sleep(time.Second)
				p := ReadMBit(3 * 8)
				c.Write(p.Encode())
				v, err := wait.Wait(conv.String(p.MsgID))
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
				v, err := wait.Wait(conv.String(p.MsgID))
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
				v, err := wait.Wait(conv.String(p.MsgID))
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
				v, err := wait.Wait(conv.String(p.MsgID))
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
	redial.RunTCP("192.168.10.66:102", func(c *client.Client) {
		c.Logger.Debug()
		c.Logger.WithHEX()
		c.OnReadFrom = ReadFrom
		//c.SetReadFunc(ReadFrom)
		c.Write(ShakeHandS7200Smart1.Bytes())
		c.Write(ShakeHandS7200Smart2.Bytes())
		c.OnDealMessage = func(c *client.Client, msg ios.Acker) {
			if ShakeHandS7200Smart1Regexp.Regexp(hex.EncodeToString(msg.Payload())) ||
				ShakeHandS7200Smart2Regexp.Regexp(hex.EncodeToString(msg.Payload())) {
				return
			}
			p, err := Decode(msg.Payload())
			if !logs.PrintErr(err) {
				wait.Done(conv.String(p.MsgID), p)
			}
		}
		c.OnDealMessage = func(c *client.Client, msg ios.Acker) {
			if ShakeHandS7200Smart1Regexp.Regexp(hex.EncodeToString(msg.Payload())) ||
				ShakeHandS7200Smart2Regexp.Regexp(hex.EncodeToString(msg.Payload())) {
				return
			}
			p, err := Decode(msg.Payload())
			if !logs.PrintErr(err) {
				wait.Done(conv.String(p.MsgID), p)
			}
		}
		go func() {

			{
				time.Sleep(time.Second)
				p := WriteDBBit(300*8, false)
				c.Write(p.Encode())
				v, err := wait.Wait(conv.String(p.MsgID))
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
				v, err := wait.Wait(conv.String(p.MsgID))
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
				v, err := wait.Wait(conv.String(p.MsgID))
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
				v, err := wait.Wait(conv.String(p.MsgID))
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
				v, err := wait.Wait(conv.String(p.MsgID))
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
