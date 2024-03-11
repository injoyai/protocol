package knx

type Service uint16

func (this Service) Bytes() []byte {
	return []byte{byte(this >> 8), byte(this)}
}

const (

	// SearchReqService 搜索请求服务
	SearchReqService Service = 0x0201

	// SearchResService 搜索响应服务
	SearchResService Service = 0x0202

	// DescrReqService 描述请求服务
	DescrReqService Service = 0x0203

	// DescrResService 描述响应服务
	DescrResService Service = 0x0204

	// ConnReqService 连接请求服务
	ConnReqService Service = 0x0205

	// ConnResService 连接响应服务
	ConnResService Service = 0x0206

	// ConnStateReqService 连接状态请求服务
	ConnStateReqService Service = 0x0207

	// ConnStateResService 连接状态响应服务
	ConnStateResService Service = 0x0208

	// DiscReqService 断开请求服务
	DiscReqService Service = 0x0209

	// DiscResService 断开响应服务
	DiscResService Service = 0x020a

	// TunnelReqService 数据请求服务
	TunnelReqService Service = 0x0420

	// TunnelResService 数据响应服务
	TunnelResService Service = 0x0421

	// RoutingIndService x
	RoutingIndService  Service = 0x0530
	RoutingLostService Service = 0x0531
	RoutingBusyService Service = 0x0532
)
