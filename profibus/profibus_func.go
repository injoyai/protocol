package profibus

import (
	"bufio"
	"io"
)

func Read() {

}

func ReadFunc(r *bufio.Reader) ([]byte, error) {
	for {

		b, err := r.ReadByte()
		if err != nil {
			return nil, err
		}

		switch b {
		case SC:

			return []byte{b}, nil

		case SD1:

			length := 5
			buf := make([]byte, length)
			n, err := io.ReadAtLeast(r, buf, length)
			if err != nil {
				return nil, err
			}
			return append([]byte{b}, buf[:n]...), nil

		case SD3:

			length := 13
			buf := make([]byte, length)
			n, err := io.ReadAtLeast(r, buf, length)
			if err != nil {
				return nil, err
			}
			return append([]byte{b}, buf[:n]...), nil

		case SD4:

			length := 2
			buf := make([]byte, length)
			n, err := io.ReadAtLeast(r, buf, length)
			if err != nil {
				return nil, err
			}
			return append([]byte{b}, buf[:n]...), nil

		case SD2:

			b, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			length := int(b) + 6
			buf := make([]byte, length)
			n, err := io.ReadAtLeast(r, buf, length)
			if err != nil {
				return nil, err
			}
			return append([]byte{b}, buf[:n]...), nil

		}

	}
}
