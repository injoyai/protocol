package icmp

import (
	"testing"
	"time"
)

func TestPing(t *testing.T) {
	t.Log(Ping("192.168.10.104", time.Second))
}
