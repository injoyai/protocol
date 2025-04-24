package j60

import (
	"github.com/injoyai/conv"
	"testing"
)

func TestID(t *testing.T) {
	t.Log(ID(1, 2))
	t.Log(ID(2, 2))
}

func TestUintToFloat(t *testing.T) {
	t.Log(UintToFloat(24, -20, 200, 7))
}

func TestDecodeControlResp(t *testing.T) {
	ls := [][8]byte{
		{172, 128, 169, 185, 54, 117, 127, 46},
		{234, 94, 96, 22, 54, 15, 127, 46},
		{234, 172, 177, 40, 54, 1, 127, 47},
		{221, 136, 207, 91, 201, 155, 128, 45},
		{99, 134, 127, 253, 200, 136, 128, 44},

		{157, 255, 199, 250, 127, 21, 128, 49},
		{189, 0, 176, 249, 127, 133, 127, 50},

		{85, 103, 72, 252, 127, 93, 127, 52},
		{87, 103, 40, 5, 128, 94, 127, 51},
	}

	for _, v := range ls {
		resp := DecodeControlResp(v)
		t.Log(resp)
	}

}

func TestDecodeControlResp2(t *testing.T) {
	t.Log(FloatToUint(10, -40, 40, 16))
	t.Log(UintToFloat(32606, -40, 40, 16))
	t.Log(conv.BIN([]byte{189, 0, 176}))
}
