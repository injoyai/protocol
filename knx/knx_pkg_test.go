package knx

import (
	"testing"
)

func TestFrame_Bytes(t *testing.T) {
	p := &Frame{
		Service: ConnReqService,
		HAPIs: []HAPI{
			NewHAPIAddress(UDP, []byte{192, 168, 0, 100}, 20),
			NewHAPIAddress(UDP, []byte{192, 168, 0, 1}, 20),
			NewCRI(),
		},
	}
	t.Logf("%#v", p)
	t.Log(p.Bytes().HEX())
}
