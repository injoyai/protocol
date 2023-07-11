package s7

import (
	"github.com/injoyai/conv"
	"testing"
)

func TestValue(t *testing.T) {
	t.Log(float64(float32(1.1)))
	t.Log(conv.Float64(float32(1.1)))
	t.Log(GetFloat([]byte{0x3f, 0x8c, 0xcc, 0xcd}))
	t.Log(GetFloat([]byte{0x3f, 0x8c, 0xcc, 0xcd}))
	t.Log(GetBool([]byte{1, 0, 0, 0}, 0))
	t.Log(GetBool([]byte{}))
	t.Log(GetString([]byte{4, 116, 101, 115, 116, 0, 0, 0, 0, 0}))
	t.Log(GetInt([]byte{0, 0, 0, 4}))
}

func TestGetBool(t *testing.T) {
	t.Log(conv.BINStr(2))
	t.Log(conv.BINStr([]byte{2, 0}))
	t.Log(GetBool([]byte{2, 0, 0, 0}, 1))
}
