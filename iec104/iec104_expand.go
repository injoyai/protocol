package iec104

import (
	"github.com/injoyai/base/g"
	"io"
)

type Handler func(writer io.Writer, a Response) error

// Do 处理服务端过来的数据
func (this Handler) Do(writer io.Writer, bytes g.Bytes) error {
	a, err := Decode(bytes)
	if err != nil {
		return err
	}
	switch a.Type {
	case 0:
		//没有APDU
		switch a.APCI.Control1 {
		case STARTDT_C:
			writer.Write(NewSTARTDT_A())
		case STOP_C:
			writer.Write(NewSTOP_A())
		case TESTFR_C:
			writer.Write(NewTESTFR_A())
		case Order_A:
			//确认信息
		}

	default:

		switch true {
		case a.Type == TypeZHTotal && a.Reason == ReasonZHResponse:
		//总召唤确认,响应数据
		case a.Type == TypeZHYM && a.Reason == ReasonZHResponse:
		//电度总召唤,响应数据
		default:

			if err := this(writer, a); err != nil {
				return err
			}
			//发送确认数据 S帧
			writer.Write(NewAck(a.APCI.WriteNo() + 2))

		}

	}
	return nil
}
