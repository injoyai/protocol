package j60

import (
	"errors"
	"fmt"
)

type ControlReq struct {
	Position float64 //角度,按角度传入
	Velocity float64 //角速度
	Torque   float64 //扭矩
	Kp       float64 //刚度
	Kd       float64 //阻尼
}

func (this ControlReq) Bytes() (*[8]byte, error) {
	if this.Position > 2400 || this.Position < -2400 {
		return nil, errors.New("角度超出范围[-2400,2400]")
	}
	//if this.Position > 40 || this.Position < -40 {
	//	return nil, errors.New("角度超出范围[-40,40]")
	//}
	if this.Velocity > 40 || this.Velocity < -40 {
		return nil, errors.New("角速度超出范围[-40,40]")
	}
	if this.Torque > 40 || this.Torque < -40 {
		return nil, errors.New("扭矩超出范围[-40,40]")
	}
	if this.Kp > 1023 || this.Kp < 0 {
		return nil, errors.New("刚度超出范围[0,1023]")
	}
	if this.Kd > 51 || this.Kd < 0 {
		return nil, errors.New("阻尼超出范围[0,51]")
	}
	position := FloatToUint(this.Position/60, -40, 40, 16)
	velocity := FloatToUint(this.Velocity, -40, 40, 14)
	torque := FloatToUint(this.Torque, -40, 40, 16)
	kp := FloatToUint(this.Kp, 0, 1023, 10)
	kd := FloatToUint(this.Kd, 0, 51, 8)
	data := [8]byte{}
	data[0] = byte(position)
	data[1] = byte(position >> 8)
	data[2] = byte(velocity)
	data[3] = byte(((velocity >> 8) & 0x3f) | ((kp & 0x03) << 6))
	data[4] = byte(kp >> 2)
	data[5] = byte(kd)
	data[6] = byte(torque)
	data[7] = byte(torque >> 8)
	return &data, nil
}

func FloatToUint(f float64, min, max float64, bits uint) uint32 {
	span := max - min
	offset := min
	return (uint32)((f - offset) * float64(uint32(1)<<bits-1) / span)
	return (uint32)((f - offset) * ((float64)((int(1) << bits) - 1)) / span)
}

func UintToFloat(n uint32, min, max float64, bits uint) float64 {
	span := max - min
	offset := min
	return (float64)(n)*span/((float64)((int(1)<<bits)-1)) + offset
}

// ID cmd左移5位，占6位，id占5位
func ID(id uint8, cmd uint8) uint32 {
	return uint32(id&0x1F) | (uint32(cmd) << 5)
}

func DID(id uint32) (uint8, uint8, bool) {
	return uint8(id & 0xF), uint8(id>>5) & 0x3F, id&0x10 == 0x10
}

type ControlResp struct {
	Position    float64 //角度
	Velocity    float64 //角速度
	Torque      float64 //扭矩
	Temperature float64 //温度
}

func (this *ControlResp) String() string {
	return fmt.Sprintf("角度:%f,角速度:%f,扭矩:%f,温度:%f", this.Position, this.Velocity, this.Torque, this.Temperature)
}

func DecodeControlResp(bs [8]byte) *ControlResp {
	data := new(ControlResp)
	position := uint32(bs[0]) | (uint32(bs[1]) << 8) | (uint32(bs[2]&0x0F) << 16) //20位
	velocity := uint32(bs[2]&0xF0) | (uint32(bs[3]) << 4) | (uint32(bs[4]) << 12) //20位
	torque := uint32(bs[5]) + uint32(bs[6])<<8                                    //16位
	flag := bs[7] & 0x01
	temp := uint32(bs[7] >> 1) //7位

	data.Position = UintToFloat(position, -40, 40, 20) * 60
	data.Velocity = UintToFloat(velocity, -40, 40, 20)
	data.Torque = UintToFloat(torque, -40, 40, 16)
	_ = flag
	data.Temperature = UintToFloat(temp, -20, 200, 7)
	return data
}
