package s7

import "github.com/injoyai/conv"

func GetInt(bytes []byte) int {
	return conv.Int(bytes)
}

func GetFloat(bytes []byte) float64 {
	return conv.Float64(bytes)
}

func GetString(bytes []byte) string {
	if len(bytes) > 0 {
		length := int(bytes[0])
		bytes = bytes[1:]
		if len(bytes) >= length {
			return string(bytes[:length])
		} else {
			return string(bytes)
		}
	}
	return ""
}

func GetBool(bytes []byte, idxs ...uint8) bool {
	idx := conv.GetDefaultUint8(0, idxs...)
	idx = conv.SelectUint8(idx <= 7, idx, 7)
	if len(bytes) > 0 {
		return conv.BIN(bytes[0])[7-idx]
	}
	return false
}
