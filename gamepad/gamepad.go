package gamepad

import (
	"errors"
	"fmt"
	"github.com/injoyai/conv"
)

/*
Decode

	[0~4)字节是左上摇杆
		0.1字节是左右方向,默认是0x0080,左0x0000,右0xFFFF
		2.3字节是上下方向,默认是0xFF7F,上是0x0000,下是0xFFFF
	[4~8)字节是右下摇杆
		4.5字节是左右方向,默认是0x0080,左0x0000,右0xFFFF
		6.7字节是上下方向,默认是0xFF7F,上是0x0000,下是0xFFFF
	[8~10)字节是
		8
		9字节是默认0x80,RT变0x00,LT过度0xFF
	10字节默认是0x0,A是0x01,B是0x02,X是0x04,Y是0x08
	10字节,默认是0x0,LB是0x10,RB是0x20,截图键是0x40,菜单键是0x80
	11字节是方向键,默认0x0,上是0x04,下是0x14,左是0x1c,右是0x0c
*/
func Decode(bs []byte) (*Status, error) {
	if len(bs) < 12 {
		return nil, errors.New("无效数据,至少需要12字节")
	}

	fn := func(b1, b2 byte, half uint16) float64 {
		_half := float64(half)
		f := float64(conv.Uint16([]byte{b2, b1}))
		return (f - _half) / _half
	}

	s := &Status{
		Joystick1: Joystick{
			X: fn(bs[0], bs[1], 0x8000),
			Y: -fn(bs[2], bs[3], 0x7FFF),
		},
		Joystick2: Joystick{
			X: fn(bs[4], bs[5], 0x8000),
			Y: -fn(bs[6], bs[7], 0x7FFF),
		},
		Direction:  Direction(bs[11]),
		A:          bs[10]&0x01 == 0x01,
		B:          bs[10]&0x02 == 0x02,
		X:          bs[10]&0x04 == 0x04,
		Y:          bs[10]&0x08 == 0x08,
		LB:         bs[10]&0x10 == 0x10,
		RB:         bs[10]&0x20 == 0x20,
		Screenshot: bs[10]&0x40 == 0x40,
		Menu:       bs[10]&0x80 == 0x80,
	}
	if bs[9] > 0x80 {
		s.LT = float64(0xFF-bs[9]) / float64(0x80)
	} else if bs[9] < 0x80 {
		s.RT = 1 - float64(bs[9])/float64(0x80)
	}

	return s, nil
}

type Status struct {
	Joystick1  Joystick  //左上摇杆
	Joystick2  Joystick  //右下摇杆
	Direction  Direction //方向键
	A, B, X, Y bool      //按键
	LB, RB     bool      //
	LT, RT     float64   //
	Screenshot bool      //截图键
	Menu       bool      //菜单键
}

func (this *Status) String() string {
	s := ""
	if this.Joystick1.X != 0 {
		s += fmt.Sprintf("左上摇杆 X:%.2f\n", this.Joystick1.X)
	}
	if this.Joystick1.Y != 0 {
		s += fmt.Sprintf("左上摇杆 Y:%.2f\n", this.Joystick1.Y)
	}
	if this.Joystick2.X != 0 {
		s += fmt.Sprintf("右下摇杆 X:%.2f\n", this.Joystick2.X)
	}
	if this.Joystick2.Y != 0 {
		s += fmt.Sprintf("右下摇杆 Y:%.2f\n", this.Joystick2.Y)
	}
	if this.Direction != 0 {
		s += fmt.Sprintf("方向键:%s\n", this.Direction)
	}
	if this.A {
		s += "A\n"
	}
	if this.B {
		s += "B\n"
	}
	if this.X {
		s += "X\n"
	}
	if this.Y {
		s += "Y\n"
	}
	if this.Screenshot {
		s += "截图\n"
	}
	if this.Menu {
		s += "菜单\n"
	}
	if this.LB {
		s += "LB\n"
	}
	if this.RB {
		s += "RB\n"
	}
	if this.LT != 0 {
		s += fmt.Sprintf("LT:%.2f\n", this.LT)
	}
	if this.RT != 0 {
		s += fmt.Sprintf("RT:%.2f\n", this.RT)
	}
	return s
}

func (this *Status) Valid() bool {
	return this.Joystick1.X != 0 ||
		this.Joystick1.Y != 0 ||
		this.Joystick2.X != 0 ||
		this.Joystick2.Y != 0 ||
		this.A || this.B || this.X || this.Y || this.Screenshot || this.Menu ||
		this.LB || this.RB ||
		this.LT != 0 || this.RT != 0 ||
		this.Direction > 0
}

type Joystick struct {
	X float64 //左(-)右(+),[0,1]
	Y float64 //下(-)上(+),[0,1]
}

const (
	Up        Direction = 4
	UpRight   Direction = 8
	Right     Direction = 12
	DownRight Direction = 16
	Down      Direction = 20
	DownLeft  Direction = 24
	Left      Direction = 28
	UpLeft    Direction = 32
)

// Direction 方向键
type Direction uint8

func (this Direction) String() string {
	s := ""
	if this.Right() {
		s += "右"
	} else if this.Left() {
		s += "左"
	}
	if this.Up() {
		s += "上"
	} else if this.Down() {
		s += "下"
	}
	return s
}

func (this Direction) Up() bool {
	return this == 32 || this == 4 || this == 8
}

func (this Direction) Down() bool {
	return this == 16 || this == 20 || this == 24
}

func (this Direction) Left() bool {
	return this == 24 || this == 28 || this == 32
}

func (this Direction) Right() bool {
	return this == 8 || this == 12 || this == 16
}

// Clock 时钟方向,1~12 1.5,3,4.5,6,7.5,9,10.5,12
func (this Direction) Clock() float64 {
	n := float64(this/4-1) * 1.5
	n = conv.Select(n == 0, 12, n)
	return n
}
