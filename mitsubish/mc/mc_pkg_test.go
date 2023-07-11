package mc

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"strings"
	"testing"
)

func TestPkg_Bytes(t *testing.T) {
	{
		p := &Pkg{
			Addr:  7000,
			Block: BlockD,
			Order: OrderWrite,
			Point: 1,
			Data:  []byte{0x00, 0x0C},
		}
		t.Log(p.Bytes().HEX())
		if p.Bytes().HEX() != "500000FFFF03000E00100001140000581B00A801000C00" {
			t.Error("编码有误")
		}
	}
	{
		p := &Pkg{
			Addr:  7000,
			Block: BlockD,
			Order: OrderRead,
			Point: 5,
		}
		t.Log(p.Bytes().HEX())
		if p.Bytes().HEX() != "500000FFFF03000C00100001040000581B00A80500" {
			t.Error("编码有误")
		}
	}
}

func TestDecode(t *testing.T) {
	s := "D0 00 00 FF FF 03 00 0C 00 00 00 0C 00 00 00 00 00 00 00 00 00"
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
	t.Log(p.Data.HEX())
}

func TestReadFunc(t *testing.T) {
	s := "D0 00 00 FF FF 03 00 0C 00 00 00 0C 00 00 00 00 00 00 00 00 00"
	s = strings.ReplaceAll(s, " ", "")
	bs, err := hex.DecodeString(s)
	if err != nil {
		t.Error(err)
		return
	}
	bs, err = ReadFunc(bufio.NewReader(bytes.NewReader(bs)))
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(hex.EncodeToString(bs))
	if strings.ToUpper(hex.EncodeToString(bs)) != s {
		t.Error("读取错误")
	}
}
