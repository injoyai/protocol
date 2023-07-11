package ip

import (
	"encoding/hex"
	"testing"
)

func TestDecodeEth(t *testing.T) {
	s := "8cf3195ad7702cea7ffebf8100000000000000000000000000007f0000010808080810e1005000000000000000000000000000000000000102030000"
	bs := []byte{44, 234, 127, 254, 191, 129, 140, 243, 25, 90, 215, 112, 8, 0, 69, 0, 0, 66, 87, 70, 0, 0, 30, 6, 175, 197, 192, 168, 10, 66, 192, 168, 10, 24, 0, 102, 192, 11, 0, 103, 56, 250, 24, 158, 164, 201, 80, 24, 32, 0, 15, 230, 0, 0, 3, 0, 0, 26, 2, 240, 128, 50, 3, 0, 0, 166, 156, 0, 2, 0, 5, 0, 0, 4, 1, 255, 3, 0, 1, 1, 0, 0, 0, 0}
	bs, _ = hex.DecodeString(s)
	dp, err := DecodeMac(bs)
	if err != nil {
		t.Error(err)
	}
	t.Log(dp)
	t.Log(bs)
	t.Log(dp.Bytes())
	t.Log(dp.Payload())
	dp2, err := DecodeIP(dp.Payload())
	if err != nil {
		t.Error(err)
	}
	t.Log(dp2)
	dp3, err := DecodeTCP(dp2.Payload())
	if err != nil {
		t.Error(err)
	}
	t.Log(dp3)
	t.Log(dp3.Payload())
	t.Log(string(dp3.Payload()))

}
