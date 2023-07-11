package base64

import (
	"encoding/base64"
	"testing"
)

func TestEncode(t *testing.T) {
	bs := []byte("6666")
	t.Log(base64.StdEncoding.EncodeToString(bs)) //NjY2Ng==
	t.Log(Encode(bs))                            //NjY2Ng==
	bs = []byte{0x00, 0x01}
	t.Log(base64.StdEncoding.EncodeToString(bs)) //AAE=
	t.Log(Encode(bs))                            //AAE=
	bs = []byte{}
	t.Log(base64.StdEncoding.EncodeToString(bs))
	t.Log(Encode(bs))
}

func TestDecode(t *testing.T) {
	t.Log(string(Decode("NjY2Ng==")))
	t.Log(Decode("AAE="))
	t.Log(Decode(""))
}
