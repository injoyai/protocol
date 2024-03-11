package ads

import (
	"fmt"
	"github.com/injoyai/base/g"
	"github.com/injoyai/conv"
)

// InfoResp 读取名称和版本号响应
type InfoResp struct {
	Result       uint32 //ADS error number , 翻译: ADS错误编号
	MajorVersion uint8  //Major version number , 翻译: 主要版本号
	MinorVersion uint8  //Minor version number , 翻译: 次要版本号
	VersionBuild uint16 //Build number , 翻译: 构建号
	DeviceName   string //Device name , 翻译: 设备名称
}

type ReadReq struct {
	IndexGroup  uint32 //Index Group of the data which should be read , 翻译: 应读取的数据的索引组
	IndexOffset uint32 //Index Offset of the data which should be read , 翻译: 应读取的数据的索引偏移
	Length      uint32 //Length of the data (in bytes) which should be read , 翻译: 应读取的数据的长度
}

func (this ReadReq) Bytes() g.Bytes {
	data := []byte(nil)
	data = append(data, g.Bytes(conv.Bytes(this.IndexGroup)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.IndexOffset)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.Length)).Reverse()...)
	return data
}

type ReadResp struct {
	Result uint32 //ADS error number , 翻译: ADS错误编号
	Length uint32 //Length of data which are supplied back , 翻译: 读取的数据的长度
	Data   []byte //Data which are supplied back , 翻译: 读取的数据
}

func DecodeRead(bs g.Bytes) (*ReadResp, error) {
	if len(bs) < 8 {
		return nil, fmt.Errorf("基础长度错误: 预期(8) 得到(%d)", len(bs))
	}
	if conv.Uint32(bs[4:8].Reverse())+4 != uint32(bs.Len()) {
		return nil, fmt.Errorf("数据长度错误: 预期(%d) 得到(%d)", conv.Uint32(bs[4:8].Reverse())+4, bs.Len())
	}
	result := Result(conv.Uint32(bs[:4].Reverse()))
	return &ReadResp{
		Result: conv.Uint32(bs[:4].Reverse()),
		Length: conv.Uint32(bs[4:8].Reverse()),
		Data:   bs[8:].Reverse(),
	}, result.Err()
}

type WriteReq struct {
	IndexGroup  uint32 //Index Group of the data which should be written , 翻译: 应写入的数据的索引组
	IndexOffset uint32 //Index Offset of the data which should be written , 翻译: 应写入的数据的索引偏移
	//Length      uint32 //Length of the data in bytes which  written , 翻译: 写入的数据长度（以字节为单位）
	Data []byte //Data which are written in ADS device , 翻译: 写入ADS设备的数据
}

func (this WriteReq) Bytes() g.Bytes {
	data := []byte(nil)
	data = append(data, g.Bytes(conv.Bytes(this.IndexGroup)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.IndexOffset)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(uint32(len(this.Data)))).Reverse()...)
	data = append(data, g.Bytes(this.Data).Reverse()...)
	return data
}

type WriteResp struct {
	Result Result //ADS error number , 翻译: ADS错误编号
}

func DecodeWrite(bs g.Bytes) (*WriteResp, error) {
	if bs.Len() != 4 {
		return nil, fmt.Errorf("基础长度错误: 预期(4) 得到(%d)", bs.Len())
	}
	result := Result(conv.Uint32(bs[:4].Reverse()))
	return &WriteResp{
		Result: result,
	}, result.Err()
}

type ReadStateReq struct{}

func (this ReadStateReq) Bytes() g.Bytes {
	return nil
}

type ReadStateResp struct {
	Result      uint32 //ADS error number , 翻译: ADS错误编号
	ADSState    uint16 //ADS state (see data typeADSSTATE of the ADS-DLL) , 翻译: ADS状态（参见ADS-DLL的数据类型ADSSTATE）
	DeviceState uint16 //Device state , 翻译: 设备状态
}

type WriteControlReq struct {
	ADSState    uint16 //New ADS state (see data typeADSSTATE of the ADS-DLL) , 翻译: ADS状态（参见ADS-DLL的数据类型ADSSTATE）
	DeviceState uint16 //New Device state , 翻译: 设备状态
	//Length      uint32 //Length of the data in bytes , 翻译: 数据长度（以字节为单位）
	Data []byte //Additional Data which are sent to the ADS device , 翻译: 发送给ADS设备的附加数据
}

func (this WriteControlReq) Bytes() g.Bytes {
	data := []byte(nil)
	data = append(data, g.Bytes(conv.Bytes(this.ADSState)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.DeviceState)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(uint32(len(this.Data)))).Reverse()...)
	data = append(data, g.Bytes(this.Data).Reverse()...)
	return data
}

type WriteControlResp struct {
	Result uint32 //ADS error number , 翻译: ADS错误编号
}

type AddNoticeReq struct {
	IndexGroup       uint32   //Index Group of the data which should be sent per notification , 翻译: 应发送通知的数据的索引组
	IndexOffset      uint32   //Index Offset of the data which should be sent per notification , 翻译: 应发送通知的数据的索引偏移
	Length           uint32   //Length of the data (in bytes) which should be sent per notification , 翻译: 应发送通知的数据的长度
	TransmissionMode uint32   //See description of the structure ADSTRANSMODE at the ADS-DLL , 翻译: 参见ADS-DLL的数据结构ADSTRANSMODE的说明
	MaxDelay         uint32   //At the latest after this time, the ADS Device Notification is called. The unit is 1ms 翻译: ADS设备通知最晚调用的时间。单位为1ms
	CycleTime        uint32   //The ADS server checks if the value changes in the time slice. The unit is 1ms 翻译: ADS服务器在时间片内检查值是否发生变化。单位为1ms
	Reserved         [16]byte //Must be set to 0 翻译: 必须为0
}

func (this AddNoticeReq) Bytes() g.Bytes {
	data := []byte(nil)
	data = append(data, g.Bytes(conv.Bytes(this.IndexGroup)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.IndexOffset)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.Length)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.TransmissionMode)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.MaxDelay)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.CycleTime)).Reverse()...)
	data = append(data, g.Bytes(this.Reserved[:]).Reverse()...)
	return data
}

type AddNoticeResp struct {
	Result             uint32 //ADS error number , 翻译: ADS错误编号
	NotificationHandle uint32 //Handle of notification , 翻译: 通知句柄
}

type DelNoticeReq struct {
	NotificationHandle uint32 //Handle of notification , 翻译: 通知句柄
}

func (this DelNoticeReq) Bytes() g.Bytes {
	return g.Bytes(conv.Bytes(this.NotificationHandle)).Reverse()
}

type DelNoticeResp struct {
	Result uint32 //ADS error number , 翻译: ADS错误编号
}

type GetInfoReq struct {
	Length         uint32         //Size of data in byte 翻译: 数据大小（以字节为单位）
	Stamps         uint32         //Number of elements of type AdsStampHeader 翻译: 类型为AdsStampHandler的元素个数
	AdsStampHeader AdsStampHeader //Array with elements of type AdsStampHeader 翻译: 类型为AdsStampHeader的数组
}

func (this GetInfoReq) Bytes() g.Bytes {
	data := []byte(nil)
	data = append(data, g.Bytes(conv.Bytes(this.Length)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.Stamps)).Reverse()...)
	{
		data = append(data, g.Bytes(conv.Bytes(this.AdsStampHeader.TimeStamp)).Reverse()...)
		data = append(data, g.Bytes(conv.Bytes(this.AdsStampHeader.Samples)).Reverse()...)
		{
			data = append(data, g.Bytes(conv.Bytes(this.AdsStampHeader.AdsNotificationSample.NotificationHandle)).Reverse()...)
			data = append(data, g.Bytes(conv.Bytes(this.AdsStampHeader.AdsNotificationSample.SampleSize)).Reverse()...)
			data = append(data, g.Bytes(this.AdsStampHeader.AdsNotificationSample.Data).Reverse()...)
		}
	}
	return data
}

type AdsStampHeader struct {
	//The timeStamp is coded after the Windows FILETIME format.l.e the value contains the number of 100-nanosecond intervals. which passed since 1.1.1601. In addition, the local time change is not considered. Thus the time stamp is present as universal Coordinatend time(UTC)
	//翻译: 时间戳以Windows FILETIME格式编码。l.e该值包含100纳秒的间隔数。自601年1月1日起通过。此外，不考虑当地时间的变化。因此，时间戳显示为通用协调时间（UTC）
	TimeStamp             uint64
	Samples               uint32                //Number of elements of type AdsNotificationSample 翻译: 类型为AdsNotificationHSample的元素个数
	AdsNotificationSample AdsNotificationSample //Array with elements of type AdsNotificationSample 翻译: 类型为AdsNotificationHSample的数组
}

type AdsNotificationSample struct {
	NotificationHandle uint32 //Handle of notification , 翻译: 通知句柄
	SampleSize         uint32 //Size of data in byte 翻译: 数据大小（以字节为单位）
	Data               []byte //Data 翻译: 数据
}

type ReadWriteReq struct {
	IndexGroup  uint32 //Index Group, in which the data should be  written 翻译: 数据所在索引组
	IndexOffset uint32 //Index Offset, in which the data should be written 翻译: 数据所在索引偏移
	ReadLength  uint32 //Length of data to be read 翻译: 需要读取的数据长度
	//WriteLength uint32 //Length of data to be written 翻译: 需要写入的数据长度
	Data []byte //Data to be written in the ADS device 翻译: 需要写入的数据
}

func (this ReadWriteReq) Bytes() g.Bytes {
	data := []byte(nil)
	data = append(data, g.Bytes(conv.Bytes(this.IndexGroup)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.IndexOffset)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(this.ReadLength)).Reverse()...)
	data = append(data, g.Bytes(conv.Bytes(uint32(len(this.Data)))).Reverse()...)
	data = append(data, g.Bytes(this.Data).Reverse()...)
	return g.Bytes(data).Reverse()
}

type ReadWriteResp struct {
	Result uint32 //ADS error number , 翻译: ADS错误编号
	Length uint32 //Length of data which are supplied back
	Data   []byte //Data which are supplied back
}
