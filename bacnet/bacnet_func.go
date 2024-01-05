package bacnet

/*

在BACnet中，有十种网络层协议报文，它们的作用是进行路由器自动配置，路由表的维护，和网络层拥塞控制。下面介绍这十种报文。



*/

/*
WhoIsRouterToNetwork
格式为：网络报文类型域是X‘00’，后面有2个字节的网络号。作用是：①节点用来确定通达某目标网络的下一个路由器；②帮助路由器更新路由表。当省略2字节的网络号时，接收此报文的路由器要返回其所有可通达的目标网络的列表。
*/
func WhoIsRouterToNetwork() {

}

/*
IAmRouterToNetwork
格式为：网络报文类型域是X‘01’，后面有2个字节的网络号。作用是列出通过发送此报文的路由器可以到达的网络号。
*/
func IAmRouterToNetwork() {

}

/*
ICouldBeRouterToNetwork
格式为：网络报文类型域是X‘02’，后面有2个字节的网络号和1个字节的性能指标。作用是响应包含有特定网络号的Who-Is-Router-To-Network报文，由能够建立到达特定目标网络的点到点连接的半路由器使用，其网络号就是所响应的报文中包含的特定网络的网络号。性能指标表明这种连接的质量。
*/
func ICouldBeRouterToNetwork() {

}

/*
RejectMessageToNetwork
格式为：网络报文类型域是X‘03’，后面有1个字节的原因说明和2个字节的网络号。作用是一个拒绝报文，直接发给生成被拒绝的报文的节点，网络号就是被拒绝报文中的网络号。在拒绝原因字节中是一个无符号的整数，其值所表示的意义如下：

(0)：其它差错。

(1)：本路由器不能直接连接到所指示的目标网络以及不能发现任何一个能够连接到所指示的目标网络。

(2)：本路由器忙，目前不能接收关于所指示目标网络的报文。

(3)：这是一个不可识别的网络层报文类型。

(4)：报文太长，不能路由到所指示的目标网络。
*/
func RejectMessageToNetwork() {

}

/*
RouterBusyToNetwork
格式为：网络报文类型域是X‘04’，后面是可选择的2个字节的网络号。作用是被路由器用来通知停止接收通过本路由器向某特定目标网络或者所有网络发送的报文。此报文通常用广播MAC地址发向相应的网络。如果没有可选择的2个字节的网络号，则表示到所有网络的报文都不接收。
*/
func RouterBusyToNetwork() {

}

/*
RouterAvailableToNetwork
格式为：网络报文类型域是X‘05’，后面是可选择的2个字节的网络号。作用是被路由器用来通知开始或者重新开始接收通过本路由器向某特定目标网络或者所有网络发送的报文。此报文通常用广播MAC地址发向相应的网络。如果没有可选择的2个字节的网络号，则表示到所有网络的报文都可接收。
*/
func RouterAvailableToNetwork() {}

/*
InitializeRouterTable
格式为：网络报文类型域是X‘06’。作用是初始化一个路由器的路由表或者查询当前路由表的内容。此报文有一个数据段，包含有初始化路由表的信息。
*/
func InitializeRouterTable() {}

/*
InitializeRouterTableAck
格式为：网络报文类型域是X‘07’。作用是对初始化路由表报文的应答，表示路由器的路由表已经改变，或者已被查询。此报文的数据段具有与它应答的初始化路由表报文相同的格式。
*/
func InitializeRouterTableAck() {}

/*
EstablishConnectionToNetwork
格式为：网络报文类型域是X‘08’， 后面有2个字节的网络号和1个字节的“中止时间值”。作用是命令一个半路由器创建一个通达指定网络的点到点连接。2个字节的网络号指出要半路由器连接的目标网络。1个字节的“中止时间值”规定了在没有NPDU到达的情况下，连接保留的时间。当此值为0时，表示连接永久保留。
*/
func EstablishConnectionToNetwork() {}

// DisconnectConnectionToNetwork 格式为：网络报文类型域是X‘09’， 后面有2个字节的网络号。作用是命令一个路由器释放所建立的点到点连接。
func DisconnectConnectionToNetwork() {}
