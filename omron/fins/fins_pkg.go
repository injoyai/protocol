package fins

type Control struct {
	Gateway  bool //是否使用网关
	Order    bool //是否命令,否则是响应
	Response bool //是否需要响应
}

func (this Control) Byte() (b byte) {
	if this.Gateway {
		b += byte(1 << 7)
	}
	if this.Order {
		b += byte(1 << 6)
	}
	if this.Response {
		b += byte(1 << 0)
	}
	return
}

type Code struct {
}

/*
Request
数据帧类型如下:
	0x00：connect requst 连接请求数据帧
	0x01：connect Response，连接请求确认数据；
	0x02：data，数据传输
*/
type Request struct {

	/*
		0x00：connect requst 连接请求数据帧
		0x01：connect Response，连接请求确认数据；
		0x02：data，数据传输；
	*/
	Type uint8

	/*
		1… …. = Gateway bit，是否使用网关，0x01表示使用；
		.1… …. = Data Type bit，数据类型比特位，0x01表示为响应，0x00表示命令；
		…0. …. = Reserved bit，第一个保留比特位，默认置0；
		…0 …. = Reserved bit，第二个保留比特位，默认置0；
		…. 0… = Reserved bit，第三个保留比特位，默认置0；
		…. .0… = Reserved bit，第四个保留比特位，默认置0；
		…. …0. = Reserved bit，第五个保留比特位，默认置0；
		…. …1 = Response setting bit，第一个保留比特位响应标志为，0x01表示非必需回应，0x00表示必须进行回应。
	*/
	Control //（Information Control Field）信息控制码：

	/*
		00：本地网络
		01 to 7F：远程网络
	*/
	TargetNetwork uint8 //（Destination network address）目标网络地址。 00：本地网络 01 to 7F：远程网络

	/*
		01 to 7E：SYSMAC NET 网络节点号
		01 to 3E：SYSMAC LINK 网络节点号
		FF：广播节点号
	*/
	TargetNode uint8 //Destination node number）目标节点号。

	/*
		00：PC（CPU）
		FE：SYSMAC NET连接单元或者SYSMAC LINK单元连接网络
		10 to 1F：CPU 总线单元
	*/
	SourceUnitNumber uint8 //（Source unit number）源单元号。

	/*
		00：本地网络
		01 to 7F：远程网络
	*/
	SourceNetwork uint8 //（Source network address）源网络地址。

	/*
		01 to 7E：SYSMAC NET 网络节点号
		01 to 3E：SYSMAC LINK 网络节点号
		FF：广播节点号
	*/
	SourceNode uint8 //（Source node number）源节点号

	/*
		00：PC（CPU）
		FE：SYSMAC NET连接单元或者SYSMAC LINK单元连接网络
		10 to 1F：CPU 总线单元
	*/
	SourceUnitAddress uint8 //（Source Unit address）源单元地址

	ServiceID uint8 //（Service ID） 序列号 范围00-FF

	Code
	Data []byte
}

/*
Bytes
46 49 4e 53 00 00 00 15 00 00 00 02 00 00 00 00 80 00 02 00 7a 00 00 00 ef 05 05 01 00
FINS(4)		数据长度(4)	保留(3)  类型 错误码(4)	Header
*/
func (this *Request) Bytes() []byte {
	data := []byte("FINS")          //协议ID,4字节
	data = append(data, 0, 0, 0, 0) //后续数据长度,占位,4字节
	data = append(data, 0, 0, 0)    //保留,3字节
	data = append(data, this.Type)  //数据帧类型,1字节
	data = append(data, 0, 0, 0, 0) //错误码,

	data = append(data, this.Control.Byte()) //（Information Control Field）信息控制码
	data = append(data, 0)                   //(Reserved）预留 一般为0x00
	data = append(data, 0x02)                //（Gateway count）网关数量，一般为0x02。
	data = append(data, this.TargetNetwork)  //（Destination network address）目标网络地址。
	data = append(data, this.TargetNode)     //（Destination node number）目标节点号
	data = append(data, this.SourceUnitNumber)
	data = append(data, this.SourceNetwork)
	data = append(data, this.SourceNode)
	data = append(data, this.SourceUnitAddress)
	data = append(data, this.ServiceID)

	data = append(data, '*', '\n')
	return data
}
