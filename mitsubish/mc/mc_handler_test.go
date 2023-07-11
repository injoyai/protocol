package mc

import (
	"encoding/hex"
	"github.com/injoyai/io"
	"github.com/injoyai/io/dial"
	"strings"
	"testing"
)

func TestNewWriteBytes(t *testing.T) {
	bs := NewWriteBytes(7000, BlockD, []byte{0x00, 0x0C})
	t.Log(hex.EncodeToString(bs))
	if strings.ToUpper(hex.EncodeToString(bs)) != "500000FFFF03000E00100001140000581B00A801000C00" {
		t.Error("编码有误")
	}
}

func TestNewReadPkgBytes(t *testing.T) {
	bs := NewReadBytes(7000, BlockD, 5)
	t.Log(hex.EncodeToString(bs))
	if strings.ToUpper(hex.EncodeToString(bs)) != "500000FFFF03000C00100001040000581B00A80500" {
		t.Error("编码有误")
	}
}

func TestServer(t *testing.T) {
	t.Error(dial.RunTCPServer(10089, func(s *io.Server) {
		s.Debug()
		s.SetDealFunc(func(msg *io.IMessage) {
			msg.Client.Write([]byte{0xD0, 0x00, 0x00, 0xFF, 0xFF, 0x03, 0x00, 0x0C, 0x00, 0x00, 0x00, 0x0C, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
		})
	}))
}
