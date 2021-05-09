package zface

type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 开启服务
	Server()

	// 添加一个路由
	AddRouter(msgId uint32, router IRouter)
}