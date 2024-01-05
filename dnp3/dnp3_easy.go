package dnp3

// ConfirmPkg 确认包,客户端连接到服务端的时候,服务端会下发一个主动上报,需要确认
func ConfirmPkg(form, to uint16) *Pkg {
	return &Pkg{
		Header: Header{
			Control: LinkControl{
				ToSlave:  true,
				IsMaster: true,
				Correct:  false,
				Function: Send,
			},
			From: form,
			To:   to,
		},
		Body: Body{
			PkgNo: PkgNo{
				IsLast:  true,
				IsFirst: true,
				Current: 0,
			},
			Control: BodyControl{
				IsFirst:     true,
				IsFinlay:    true,
				NeedAck:     false,
				Unsolicited: true,
				No:          0,
			},
			Function: Confirm,
		},
	}
}

type ReadDataReq struct {
	DataType uint16
}

func ReadPkg(from, to uint16, datas []Data) *Pkg {
	return &Pkg{
		Header: Header{
			Control: LinkControl{
				ToSlave:  true,
				IsMaster: true,
				Correct:  false,
				Function: Send,
			},
			From: from,
			To:   to,
		},
		Body: Body{
			PkgNo:    PkgNo{},
			Control:  BodyControl{},
			Function: Read,
			Datas:    datas,
		},
	}
}

func WritePkg(from, to uint16, datas []Data) *Pkg {
	return &Pkg{
		Header: Header{
			Control: LinkControl{
				ToSlave:  true,
				IsMaster: true,
				Correct:  false,
				Function: Send,
			},
			From: from,
			To:   to,
		},
		Body: Body{
			PkgNo:    PkgNo{},
			Control:  BodyControl{},
			Function: Write,
			Datas:    datas,
		},
	}
}
