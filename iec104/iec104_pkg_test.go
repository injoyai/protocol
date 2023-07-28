package iec104

import (
	"encoding/hex"
	"github.com/injoyai/io"
	"github.com/injoyai/io/dial"
	"github.com/injoyai/logs"
	io2 "io"
	"strings"
	"testing"
	"time"
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

	<-dial.RedialTCP(":2404", func(c *io.Client) {
		//握手
		c.CloseWithErr(Handshake(c))
		c.Debug()
		c.SetPrintWithHEX()

		totalDone := make(chan struct{}, 1)

		go c.Timer(time.Second, func() error {
			res, err := c.WriteReadWithTimeout(NewRead(1, 2), time.Second*5)
			logs.PrintErr(err)
			handler(c, res)
			return nil
		})

		go c.Run()
		return

		go func() {
			for {
				c.Write(NewZHTotal(1))
				<-totalDone
			}
		}()

		c.SetDealFunc(func(msg *io.IMessage) {
			c.CloseWithErr(handler(msg.Client, msg.Bytes()))
		})
	}).DoneAll()
}

func TestPkg(t *testing.T) {
	t.Log(NewZHTotal(1).HEX())
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
