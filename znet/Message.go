package znet

type Message struct {
	MsgId  uint32
	MsgLen uint32
	Msg    []byte
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
