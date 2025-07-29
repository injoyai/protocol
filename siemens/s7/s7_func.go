package s7

import (
	"github.com/injoyai/conv"
	"io"
)

//// ReadFunc 西门子协议数据读取
//func ReadFunc(r *bufio.Reader) (bytes []byte, err error) {
//	readOne := func() (byte, error) {
//		buf := make([]byte, 1)
//		_, err := r.Read(buf)
//		if err != nil {
//			panic(err)
//		}
//		return buf[0], err
//	}
//	for {
//		b, err := readOne()
//		if err != nil {
//			return
//		}
//		if b == 0x03 {
//			b2, err := readOne()
//			if err != nil {
//				return
//			}
//			if b2 == 0x00 && len(bytes) == 0 {
//				bytes = []byte{b, b2}
//			} else {
//				bytes = append(bytes, b, b2)
//			}
//		} else {
//			bytes = append(bytes, b)
//		}
//		if len(bytes) >= 17 {
//			length := conv.Int(bytes[2:4])
//			if len(bytes) == length {
//				break
//			}
//		}
//	}
//	return
//}

// ReadFrom 西门子协议数据读取
func ReadFrom(r io.Reader) (bytes []byte, err error) {
	readOne := func() (byte, error) {
		buf := make([]byte, 1)
		_, err := r.Read(buf)
		if err != nil {
			panic(err)
		}
		return buf[0], err
	}
	for {
		b, err := readOne()
		if err != nil {
			return
		}
		if b == 0x03 {
			b2, err := readOne()
			if err != nil {
				return
			}
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
