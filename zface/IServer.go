package zface

type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 开启服务
	Server()
}