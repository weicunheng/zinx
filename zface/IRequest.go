package zface

type IRequest interface {

	// 获取当前链接的方法
	GetConnection() IConnection

	// 获取数据
	GetData() []byte

	// 获取MsgId
	GetMsgId() uint32
}
