package ip

import (
	"testing"
)

func TestDecodeTCP(t *testing.T) {
	bs := []byte{252, 111, 1, 187, 167, 206, 24, 192, 67, 218, 185, 221, 80, 24, 2, 0, 45, 48, 0, 0}
	t.Log(DecodeTCP(bs))
}
