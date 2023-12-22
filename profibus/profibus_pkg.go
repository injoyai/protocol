package profibus

import (
	"errors"
	"fmt"
	"github.com/injoyai/base/g"
)

const (
	SD1 = 0x10 //测试连接.类似ping
	SD2 = 0x68 //可变数据帧
	SD3 = 0xA2 //定长数据帧
	SD4 = 0xDC //令牌传递
	SC  = 0xE5 //短应答
	ED  = 0x16 //帧尾
)

/*
Pkg
参考文档
https://zhuanlan.zhihu.com/p/643231451
https://zhuanlan.zhihu.com/p/596431958
https://blog.csdn.net/fishwang89/article/details/14010439

*/
type Pkg struct {
	SD       uint8    //帧头
	DA       uint8    //从站地址
	Master   uint8    //主站地址
	Function Function //功能码
	SlaveAP  uint8
	MasterAP uint8
	Data     []byte //数据
}

func (this *Pkg) sum(data []byte) uint8 {
	var sum uint8
	for i, v := range data {
		if i > 0 {
			sum += v
		}
	}
	return sum
}

func (this *Pkg) Suffix(data []byte) []byte {
	data = append(data, this.sum(data))
	data = append(data, ED)
	return data
}

func (this *Pkg) Bytes() g.Bytes {

	switch this.SD {
	case SC:

		return []byte{SC}

	case SD1:

		data := []byte{
			SD1,
			this.DA,
			this.Master,
			this.Function.Byte(),
		}
		return this.Suffix(data)

	case SD3:

		data := []byte{
			SD3,
			this.DA,
			this.Master,
			this.Function.Byte(),
		}
		data = append(data, this.Data...)
		return this.Suffix(data)

	case SD4:

		return []byte{
			SD4,
			this.DA,
			this.Master,
		}

	case SD2:

		data := []byte{
			SD2,
			byte(len(this.Data)) + 5,
			byte(len(this.Data)) + 5,
			SD2,
			this.DA,
			this.Master,
			this.Function.Byte(),
			this.SlaveAP,
			this.MasterAP,
		}
		data = append(data, this.Data...)
		return this.Suffix(data)

	}

	return []byte{}
}

type Function struct {
	Frame    bool  //帧类型，“=1”，请求帧；“=0”，应答/回答帧。
	StnType  uint8 //站类型和 FDL 状态（帧类型 BIT6=0）。具体功能见表1-1
	Function uint8 //功能码，具体解释见 表1-2
}

func (this Function) Byte() (b byte) {
	if this.Frame {
		b += 1 << 6
	}
	b += (this.StnType % 4) << 4
	b += this.Function % 16
	return 0
}

func Decode(bs []byte) (*Pkg, error) {

	if len(bs) == 0 {
		return nil, errors.New("数据长度为0")
	}

	switch bs[0] {
	case SC:

		if len(bs) != 1 {
			return nil, fmt.Errorf("数据长度错误,预期(%d),得到(%d)", 1, len(bs))
		}

		return &Pkg{SD: SC}, nil

	case SD1:

		if len(bs) != 6 {
			return nil, fmt.Errorf("数据长度错误,预期(%d),得到(%d)", 6, len(bs))
		}

		return &Pkg{
			SD:     bs[0],
			DA:     bs[2],
			Master: bs[1],
			Function: Function{
				Frame:    bs[3]<<1>>7 == 1,
				StnType:  bs[3] << 2 >> 6,
				Function: bs[3] << 4 >> 4,
			},
		}, nil

	case SD3:

		if len(bs) != 14 {
			return nil, fmt.Errorf("数据长度错误,预期(%d),得到(%d)", 14, len(bs))
		}
		return &Pkg{
			SD:     SD3,
			DA:     bs[2],
			Master: bs[1],
			Function: Function{
				Frame:    bs[3]<<1>>7 == 1,
				StnType:  bs[3] << 2 >> 6,
				Function: bs[3] << 4 >> 4,
			},
			Data: bs[4:12],
		}, nil

	case SD4:

		return &Pkg{
			SD:     SD4,
			DA:     bs[2],
			Master: bs[1],
		}, nil

	case SD2:

		length := int(bs[1] + 7)
		if len(bs) != length {
			return nil, fmt.Errorf("数据长度错误,预期(%d),得到(%d)", length, len(bs))
		}
		bs = bs[3:]

		return &Pkg{
			SD:     SD2,
			DA:     bs[2],
			Master: bs[1],
			Function: Function{
				Frame:    bs[3]<<1>>7 == 1,
				StnType:  bs[3] << 2 >> 6,
				Function: bs[3] << 4 >> 4,
			},
			SlaveAP:  bs[4],
			MasterAP: bs[5],
			Data:     bs[6 : len(bs)-2],
		}, nil

	}

	return nil, errors.New("未知类型数据")
}
