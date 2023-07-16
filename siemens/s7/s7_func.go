package s7

import (
	"bufio"
	"github.com/injoyai/conv"
	"github.com/injoyai/goutil/g"
)

// ReadFunc 西门子协议数据读取
func ReadFunc(r *bufio.Reader) (bytes []byte, err error) {
	readOne := func() byte {
		buf := make([]byte, 1)
		_, err := r.Read(buf)
		g.PanicErr(err)
		return buf[0]
	}
	defer g.Recover(&err)
	for {
		b := readOne()
		if b == 0x03 {
			b2 := readOne()
			if b2 == 0x00 && len(bytes) == 0 {
				bytes = []byte{b, b2}
			} else {
				bytes = append(bytes, b, b2)
			}
		} else {
			bytes = append(bytes, b)
		}
		if len(bytes) >= 17 {
			length := conv.Int(bytes[2:4])
			if len(bytes) == length {
				break
			}
		}
	}
	return
}
