package zface

type IMsgHandler interface {
	// 1. 从map中获取router, 并执行
	DoMsgHandler(request IRequest)
	// 2. 添加router
	AddRouter(msgId uint32, router IRouter)
	// 3. 启动工作池
	StartWorkerPool()
	// 4. 将消息交给TaskQueue，由Worker处理
	SendReqToTaskQueue(request IRequest)
}
