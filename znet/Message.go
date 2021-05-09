package znet
/*
消息模块，Message 是 Request对象数据载体
*/
type Message struct {
	MsgId  uint32
	MsgLen uint32
	Msg    []byte
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		MsgId: id,
		MsgLen: uint32(len(data)),
		Msg: data,
	}
}

func (m *Message) GetMsgData() []byte {
	return m.Msg
}
func (m *Message) GetMsgId() uint32 {
	return m.MsgId
}
func (m *Message) GetMsgLen() uint32 {
	return m.MsgLen
}

func (m *Message) SetMsgData(data []byte) {
	m.Msg = data
}
func (m *Message) SetMsgLen(len uint32) {
	m.MsgLen = len
}
func (m *Message) SetMsgId(id uint32) {
	m.MsgId = id
}
