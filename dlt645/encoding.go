package dlt645

import (
	"encoding/hex"
	"errors"
)

func Encode(addr, mark, pwd string, ctl Control, value string) ([]byte, error) {
	addrBs, err := hex.DecodeString(addr)
	if err != nil || len(addrBs) != 6 {
		return nil, errors.New("通讯地址有误:" + addr)
	}
	markBs, err := hex.DecodeString(mark)
	if err != nil || len(markBs) != 4 {
		return nil, errors.New("数据标识有误:" + mark)
	}
	pwdBs, err := hex.DecodeString(pwd)
	if err != nil || (len(pwdBs) != 3 && len(pwdBs) != 0) {
		return nil, errors.New("密码有误:" + pwd)
	}
	valueBs, err := hex.DecodeString(value)
	if err != nil {
		return nil, errors.New("数据域有误:" + value)
	}
	p := &EnPkg{
		No:       addrBs,
		Control:  ctl,
		Password: pwdBs,
		Mark:     markBs,
		Data:     valueBs,
	}
	return p.Bytes(), nil
}

func Read(addr, mark string, pwd string) ([]byte, error) {
	return Encode(addr, mark, pwd, ControlGet, "")
}

func Write(addr, mark, pwd string, value string) ([]byte, error) {
	return Encode(addr, mark, pwd, ControlSet, value)
}
