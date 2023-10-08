package dns

// Request dns报文基础结构部分
type Request struct {
	ID            uint16 //事务ID
	Flag          uint16 //标志
	QuestionCount uint16 //请求问题数
	AnswerRRS     uint16 //回答资源记录数
	AuthorityRRS  uint16 //权威名称服务器数
	AdditionalRRS uint16 //附加资源记录数
}
