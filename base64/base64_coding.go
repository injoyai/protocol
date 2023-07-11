package base64

const (
	str = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
)

var m = func() map[uint8]uint8 {
	m := map[uint8]uint8{}
	for i, v := range str {
		m[uint8(v)] = uint8(i)
	}
	return m
}()

func Encode(bs []byte) string {

	//判断长度是否是3的倍数
	padding := 3 - len(bs)%3
	for i := 0; padding < 3 && i < padding; i++ {
		//填充0x00
		bs = append(bs, 0x00)
	}

	bytes := []byte(nil)
	for i := 0; i <= len(bs)-3; i += 3 {
		//3字节转4字节,每字节6位(0-63),对应64个字符
		bytes = append(bytes, str[bs[i]>>2])
		bytes = append(bytes, str[bs[i]<<6>>2+bs[i+1]>>4])
		bytes = append(bytes, str[bs[i+1]<<4>>2+bs[i+2]>>6])
		bytes = append(bytes, str[bs[i+2]<<2>>2])
	}

	//修改补码数据为'='
	for n := len(bytes); padding < 3 && n > padding && padding > 0; padding-- {
		bytes[n-padding] = '='
	}

	return string(bytes)
}

func Decode(s string) []byte {

	//确保长度是4的倍数,不返回错误
	for len(s)%4 != 0 {
		s = s + "="
	}

	bytes := []byte(s)
	length := len(bytes)
	padding := 0

	//替换填充的'='
	for i := 0; len(bytes) > 0 && i < 3; i++ {
		if bytes[length-1-i] == '=' {
			bytes[length-1-i] = 'A'
			padding++
			continue
		}
		break
	}

	//4字节取后6位(根据字典),组成3字节,
	result := []byte(nil)
	for i := 0; i <= length-4; i += 4 {
		bs0 := m[bytes[i]]
		bs1 := m[bytes[i+1]]
		bs2 := m[bytes[i+2]]
		bs3 := m[bytes[i+3]]
		result = append(result, bs0<<2+bs1>>4)
		result = append(result, bs1<<4+bs2>>2)
		result = append(result, bs2<<6+bs3)
	}

	return result[:len(result)-padding]
}
