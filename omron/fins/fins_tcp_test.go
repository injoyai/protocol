package fins

import (
	"github.com/injoyai/conv"
	"testing"
)

func TestTCPRead(t *testing.T) {
	t.Log(conv.Bytes(Command(0x02)))
	b := TCPMemoryRead(AreaDM, 100<<8, 10, 1, 1)
	t.Log(b.HEX())
}

func TestShakeHand(t *testing.T) {
	t.Log(TCPShakeHand(250).HEX())
}

func TestTCPMemoryWrite(t *testing.T) {
	b := TCPMemoryWrite(AreaDM, 100<<8, [][2]byte{{1, 2}, {3, 4}}, 1, 1)
	t.Log(b.HEX())
}
