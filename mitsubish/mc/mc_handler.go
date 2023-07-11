package mc

// NewReadPkg 新建读取
func NewReadPkg(addr uint32, block Block, length uint16) *Pkg {
	return &Pkg{
		Addr:  addr,
		Block: block,
		Order: OrderRead,
		Point: length,
	}
}

func NewReadBytes(addr uint32, block Block, length uint16) []byte {
	return NewReadPkg(addr, block, length).Bytes()
}

// NewWritePkg 新建写入,协议是小端,自动转成大端,1点有2字节,不足会补0
func NewWritePkg(addr uint32, block Block, data []byte) *Pkg {
	if len(data)%2 != 0 {
		//一个点位2字节,不是偶数需要在前面补0
		data = append([]byte{0}, data...)
	}
	return &Pkg{
		Addr:  addr,
		Block: block,
		Order: OrderWrite,
		Point: uint16(len(data) / 2),
		Data:  data,
	}
}

func NewWriteBytes(addr uint32, block Block, data []byte) []byte {
	return NewWritePkg(addr, block, data).Bytes()
}
