package zface

type IMsgHandler interface {
	// 1. 从map中获取router, 并执行
	DoMsgHandler(request IRequest)
	// 2. 添加router
	AddRouter(msgId uint32, router IRouter)
}
