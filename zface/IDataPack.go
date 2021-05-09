package zface

type IDataPack interface {
	// 获取Message Head的长度
	GetHeadLen() int
	// 封包，接受消息，封装成bytes数组
	Pack(message IMessage)([]byte, error)
	// 拆包, 接受bytes数组，返回拆包后的Message
	UnPack([]byte) (IMessage, error)
}