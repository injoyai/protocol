package ads

// NewRead 创建一个读命令
func NewRead(target [6]byte, targetPort uint16, source [6]byte, sourcePort uint16, group, offset, length uint32) *Pkg {
	return &Pkg{
		Header: Header{
			Reserved: 0,
		},
		AMS: AMS{
			Target:     target,
			TargetPort: targetPort,
			Source:     source,
			SourcePort: sourcePort,
			CommandID:  Read,
			Flag: Flag{
				Broadcast:      false,
				InitCommand:    false,
				UDPCommand:     false,
				TimestampAdded: false,
				HighPriority:   false,
				SysCommand:     false,
				ADSCommand:     true,
				NoReturn:       false,
				Response:       false,
			},
			ErrorCode: 0,
			InvokeID:  1,
		},
		Data: ReadReq{
			IndexGroup:  group,
			IndexOffset: offset,
			Length:      length,
		}.Bytes(),
	}
}

func NewWrite(target [6]byte, targetPort uint16, source [6]byte, sourcePort uint16, group, offset uint32, data []byte) *Pkg {
	return &Pkg{
		Header: Header{
			Reserved: 0,
		},
		AMS: AMS{
			Target:     target,
			TargetPort: targetPort,
			Source:     source,
			SourcePort: sourcePort,
			CommandID:  Write,
			Flag: Flag{
				Broadcast:      false,
				InitCommand:    false,
				UDPCommand:     false,
				TimestampAdded: false,
				HighPriority:   false,
				SysCommand:     false,
				ADSCommand:     true,
				NoReturn:       false,
				Response:       false,
			},
			ErrorCode: 0,
			InvokeID:  1,
		},
		Data: WriteReq{
			IndexGroup:  group,
			IndexOffset: offset,
			Data:        data,
		}.Bytes(),
	}
}
