package iec104

import (
	"encoding/hex"
	"github.com/injoyai/ios"
	"github.com/injoyai/ios/client"
	"github.com/injoyai/ios/client/redial"
	"github.com/injoyai/logs"
	io2 "io"
	"strings"
	"testing"
)

func TestNewAck(t *testing.T) {
	b := NewAck(299)
	t.Log(b.HEX())
}

func TestHandshake(t *testing.T) {

	handler := Handler(func(writer io2.Writer, a Response) error {

		logs.Debug(a)

		//if a.Type == TypeZHTotal && a.Reason == ReasonBreak {
		//	//结束总召唤
		//	totalDone <- struct{}{}
		//}

		return nil

	}).Do

	<-redial.TCP(":2404", func(c *client.Client) {
		//握手
		c.CloseWithErr(Handshake(c))
		c.Logger.Debug()
		c.Logger.WithHEX()
		c.OnReadFrom = ReadFrom
		//c.SetReadFunc(ReadFunc)

		//totalDone := make(chan struct{}, 1)

		//go func() {
		//	for {
		//		c.Write(NewZHTotal(1))
		//		<-totalDone
		//	}
		//}()

		c.OnDealMessage = func(c *client.Client, msg ios.Acker) {
			c.CloseWithErr(handler(c, msg.Payload()))
		}

		//c.SetDealFunc(func(c *io.Client, msg io.Message) {
		//	c.CloseWithErr(handler(c, msg.Bytes()))
		//})

	}).Done()
}

func TestDecode(t *testing.T) {
	s := "681306000200098214000100010700a11000891500"
	s = "681904000200240100000100020000b6f39d3f0073692e0fbc0717"
	s = strings.ReplaceAll(s, " ", "")
	t.Log(s)
	bs, err := hex.DecodeString(s)
	if err != nil {
		t.Log(err)
		return
	}
	a, err := Decode(bs)
	if err != nil {
		t.Log(err)
		return
	}
	t.Log(a)
}

func TestNewZHTotal(t *testing.T) {
	t.Log(NewZHTotal(1).HEX())
}

func TestNewRead(t *testing.T) {
	t.Log(NewRead(1, 1).HEX())
}
