package zface

type IMessage interface {
	GetMsgData() []byte
	GetMsgId() uint32
	GetMsgLen() uint32

	SetMsgData([]byte)
	SetMsgLen(uint32)
	SetMsgId(uint32)
}
