package bacnet

import (
	"testing"
)

func TestPackage_Bytes(t *testing.T) {
	p := &Pkg{
		BVLC: BVLC{
			IP:       []byte{192, 168, 10, 15},
			Port:     DefaultPort,
			Function: FunctionOriginalBroadcast,
		},
		NPDU: NPDU{
			Control: Control{
				HasAPDU:     false,
				HasDNet:     true,
				HasSNet:     false,
				NetResponse: false,
				Priority:    0,
			},
			DNET:    0xFFFF,
			DLEN:    0,
			DADR:    nil,
			SNET:    0,
			SLEN:    0,
			SADR:    nil,
			NetType: 0,
			Vendor:  0,
		},
		APDU: APDU{
			Type: UnConfirmedReq,
			Flag: Flag{
				SEG: false,
				MOR: false,
				SA:  false,
			},
			MaxSeg:   0,
			MaxResp:  0,
			InvokeID: 0,
		},
	}
	t.Log(p.Bytes().HEX())
}
