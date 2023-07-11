package s7

import (
	"errors"
	"github.com/injoyai/conv"
	"regexp"
)

// DecodeAddress 地址解析 todo 待完成
func DecodeAddress(addr string) (uint32, error) {
	if b, _ := regexp.MatchString(`^\d+$`, addr); b {
		return conv.Uint32(addr) * 8, nil
	}
	if b, _ := regexp.MatchString(`^\d+\.\d$`, addr); b {
		f := conv.Float64(addr)
		return uint32(f)*8 + uint32(f*10)%10, nil
	}
	if b, _ := regexp.MatchString(`^(VD|VW)\d+$`, addr); b {

	}
	if b, _ := regexp.MatchString(`^(VB|)\d+\.\d$`, addr); b {

	}
	return 0, errors.New("未知地址类型")
}
