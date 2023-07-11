package s7

import (
	"encoding/hex"
	"regexp"
	"strings"
)

type ShakeHand string

func (this ShakeHand) Bytes() []byte {
	s := strings.ReplaceAll(string(this), " ", "")
	bs, _ := hex.DecodeString(s)
	return bs
}

func (this ShakeHand) HEX() string {
	return strings.ToUpper(hex.EncodeToString(this.Bytes()))
}

func (this ShakeHand) Regexp(s string) bool {
	b, _ := regexp.MatchString(strings.ToUpper(string(this)), strings.ToUpper(s))
	return b
}

const (
	//握手

	ShakeHandS7200Smart1       ShakeHand = "03 00 00 16 11 E0 00 00 00 01 00 C1 02 10 00 C2 02 03 00 C0 01 0A"
	ShakeHandS7200Smart2       ShakeHand = "03 00 00 19 02 F0 80 32 01 00 00 CC C1 00 08 00 00 F0 00 00 01 00 01 03 C0"
	ShakeHandS7200Smart1Regexp ShakeHand = "^(0300001611d00001)[0-9]{4}(00c0010ac1021000c2020300)$"
	ShakeHandS7200Smart2Regexp ShakeHand = "^0300001B02F08032030000CCC1000800000000F0000001000100F0$"

	ShakeHandS72001 ShakeHand = "03 00 00 16 11 E0 00 00 00 01 00 C1 02 4D 57 C2 02 4D 57 C0 01 09"
	ShakeHandS72002 ShakeHand = "03 00 00 19 02 F0 80 32 01 00 00 00 00 00 08 00 00 F0 00 00 01 00 01 03 C0"

	ShakeHandS73001 ShakeHand = "03 00 00 16 11 E0 00 00 00 01 00 C0 01 0A C1 02 01 02 C2 02 01 02"
	ShakeHandS73002 ShakeHand = "03 00 00 19 02 F0 80 32 01 00 00 04 00 00 08 00 00 F0 00 00 01 00 01 01 E0"

	ShakeHandS74001 ShakeHand = "03 00 00 16 11 E0 00 00 00 01 00 C0 01 0A C1 02 01 00 C2 02 01 03"
	ShakeHandS74002 ShakeHand = "03 00 00 19 02 F0 80 32 01 00 00 04 00 00 08 00 00 F0 00 00 01 00 01 01 E0"

	ShakeHandS712001 ShakeHand = "03 00 00 16 11 E0 00 00 00 01 00 C0 01 0A C1 02 01 02 C2 02 01 00"
	ShakeHandS712002 ShakeHand = "03 00 00 19 02 F0 80 32 01 00 00 04 00 00 08 00 00 F0 00 00 01 00 01 01 E0"

	ShakeHandS715001 ShakeHand = "03 00 00 16 11 E0 00 00 00 01 00 C0 01 0A C1 02 01 02 C2 02 01 00"
	ShakeHandS715002 ShakeHand = "03 00 00 19 02 F0 80 32 01 00 00 04 00 00 08 00 00 F0 00 00 01 00 01 01 E0"
)
