package znet

import "zinx/zface"

type Request struct {

	// 当前链接
	conn zface.IConnection
	// 客户端请求的数据
	//data [] byte
	data zface.IMessage
}

// 获取当前链接的方法
func (r *Request) GetConnection() zface.IConnection{
	return r.conn
}

// 获取数据
func (r *Request)  GetData() []byte{
	return r.data.GetMsgData()
}

func (r *Request)  GetMsgId() uint32{
	return r.data.GetMsgId()
}