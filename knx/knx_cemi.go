package knx

type MessageCode uint8

const (
	// LBusmonIndCode is the message code for L_Busmon.ind.
	LBusmonIndCode MessageCode = 0x2B

	// LDataReqCode is the message code for L_Data.req.
	LDataReqCode MessageCode = 0x11

	// LDataIndCode is the message code for L_Data.ind.
	LDataIndCode MessageCode = 0x29

	// LDataConCode is the message code for L_Data.con.
	LDataConCode MessageCode = 0x2E

	// LRawReqCode is the message code for L_Raw.req.
	LRawReqCode MessageCode = 0x10

	// LRawIndCode is the message code for L_Raw.ind.
	LRawIndCode MessageCode = 0x2D

	// LRawConCode is the message code for L_Raw.con.
	LRawConCode MessageCode = 0x2F
)

type Message struct {
	Code     MessageCode //消息类型
	Info     []byte      //附加字段
	Control1 uint8       //未知作用,0xBC
	Control2 uint8       //未知作用,0xE0
	Source   uint16      //源,可以填0
	Target   uint16      //目标
	APDU     []byte
}

func NewMessage(code MessageCode, apdu []byte) {

}
