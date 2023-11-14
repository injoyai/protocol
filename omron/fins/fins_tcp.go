package fins

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
	"io"
)

type Command uint32

func (this Command) Bytes() []byte {
	return conv.Bytes(uint32(this))
}

const (
	CommandConnectRequest  Command = 0x00 //连接请求
	CommandConnectResponse Command = 0x01 //连接响应
	CommandData            Command = 0x02 //数据
)

type Area uint8

const (
	CIO    Area = 0x30
	WR     Area = 0x31
	HR     Area = 0x32
	AR     Area = 0x33
	AreaDM Area = 0x82
)

const (
	Prefix = "FINS"
)

/*
TCPFrame
例: 46494e530000001a000000020000000080000200010000010000010182006400000a
46494e53	FINS
0000001a	Length
00000002	Command
00000000	Error Code
80			ICF
00			RSV
02			GCT
00			DNA
01			DA1
00			DA2
00			SNA
01			SA1
00			SA2
00			SID
0101		读取内存数据
82			DM区
006400		地址100
000a 		长度10
*/
type TCPFrame struct {
	Command   Command
	ErrorCode uint32
	Frame     *Frame
}

func (this TCPFrame) Bytes() g.Bytes {
	data := []byte(Prefix)
	var frameBytes []byte
	if this.Frame != nil {
		frameBytes = this.Frame.Bytes()
	}
	data = append(data, conv.Bytes(uint32(len(frameBytes)+8))...) //数据长度
	data = append(data, this.Command.Bytes()...)
	data = append(data, conv.Bytes(this.ErrorCode)...)
	data = append(data, frameBytes...)
	return data
}

type TCPNode struct {
	ClientNode uint32
	ServerNode uint32
}

type TCPFrameResp struct {
	Length    uint32 //数据长度
	Command   Command
	ErrorCode uint32
	Data      []byte
}

func (this *TCPFrameResp) DecodeConnectResp() (*TCPNode, error) {
	if this.Command != CommandConnectResponse {
		return nil, errors.New("数据非连接响应")
	}
	if len(this.Data) != 8 {
		return nil, errors.New("数据长度错误")
	}
	return &TCPNode{
		ClientNode: conv.Uint32(this.Data[:4]),
		ServerNode: conv.Uint32(this.Data[4:]),
	}, nil
}

func (this *TCPFrameResp) DecodeData() (*Frame, error) {
	return &Frame{
		ICF:  DecodeICF(this.Data[0]),
		RSV:  this.Data[1],
		GCT:  this.Data[2],
		DNA:  this.Data[3],
		DA1:  this.Data[4],
		DA2:  this.Data[5],
		SNA:  this.Data[6],
		SA1:  this.Data[7],
		SA2:  this.Data[8],
		SID:  this.Data[9],
		RC:   Code(conv.Uint16(this.Data[10:12])),
		Data: this.Data[12:],
	}, nil
}

/*



 */

func ShakeHand(c io.ReadWriter) (*TCPNode, error) {
	bs := TCPShakeHandRequest()
	if _, err := c.Write(bs); err != nil {
		return nil, err
	}
	bs = make([]byte, 24)
	n, err := io.ReadAtLeast(c, bs, 24)
	if err != nil {
		return nil, err
	}
	resp, err := Decode(bs[:n])
	if err != nil {
		return nil, err
	}
	return resp.DecodeConnectResp()
}

// TCPShakeHandRequest 建立连接之后,需要进行握手操作
func TCPShakeHandRequest() g.Bytes {
	data := []byte(Prefix)
	data = append(data, 0, 0, 0, 0x0c)
	data = append(data, 0, 0, 0, 0)
	data = append(data, 0, 0, 0, 0)
	data = append(data, 0, 0, 0, 0)
	return data
}

/*
TCPMemoryRead 地址需要注意,最后一字节是位
发送: 0101 82 006400 000A
返回: 010100000102030405060708090A
*/
func TCPMemoryRead(area Area, addr uint32, length uint16, node *TCPNode) g.Bytes {
	return TCPFrame{
		Command: CommandData,
		Frame: &Frame{
			GCT: 2,
			DA1: uint8(node.ServerNode),
			SA1: uint8(node.ClientNode),
			RC:  MemoryRead,
			Data: Read{
				Area:    area,
				Address: addr,
				Length:  length,
			}.Bytes(),
		},
	}.Bytes()
}

// TCPMemoryWrite 内存写入
func TCPMemoryWrite(area Area, addr uint32, value []byte, node *TCPNode) g.Bytes {
	return TCPFrame{
		Command: CommandData,
		Frame: &Frame{
			GCT: 2,
			DA1: uint8(node.ServerNode),
			SA1: uint8(node.ClientNode),
			RC:  MemoryWrite,
			Data: Write{
				Area:    area,
				Address: addr,
				Value:   value,
			}.Bytes(),
		},
	}.Bytes()
}

func Decode(bs g.Bytes) (*TCPFrameResp, error) {

	//判断基础长度
	if bs.Len() < 24 {
		return nil, fmt.Errorf("基础长度错误,预期(%d),得到(%d)", 20, bs.Len())
	}

	//判断报文头部
	prefix := string(bs[:4])
	if prefix != Prefix {
		return nil, fmt.Errorf("报文头部错误,预期(%s),得到(%s)", Prefix, prefix)
	}

	length := conv.Uint32(bs[4:8])

	//判断报文长度是否完整
	if length+8 != uint32(bs.Len()) {
		return nil, fmt.Errorf("报文长度不完整,预期(%d),得到(%d)", length+8, bs.Len())
	}

	return &TCPFrameResp{
		Length:    length,
		Command:   Command(conv.Uint32(bs[8:12])),
		ErrorCode: conv.Uint32(bs[12:16]),
		Data:      bs[16:],
	}, nil
}

func TCPReadFunc(r *bufio.Reader) ([]byte, error) {
loop:
	for {
		for i := 0; i < len(Prefix); i++ {
			b, err := r.ReadByte()
			if err != nil {
				return nil, err
			}
			if b != Prefix[i] {
				continue loop
			}
		}
		//读取长度
		bs := make([]byte, 4)
		if _, err := io.ReadAtLeast(r, bs, 4); err != nil {
			return nil, err
		}
		length := conv.Uint32(bs)
		bs = make([]byte, length)
		_, err := io.ReadAtLeast(r, bs, int(length))
		return bs, err
	}
}
